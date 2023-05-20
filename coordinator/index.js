const express = require('express')
const bodyParser = require('body-parser')
const crypto = require('crypto')
const axios = require('axios')
const config = require('config')
const args = require('yargs').argv;


const no_parallel = args.no_parallel ?? false
const app = express()
app.use(bodyParser.json({ limit: '50mb', extended: true }))

const port = config.get('server.port')
const url = config.get('server.host')

let workers = []
let pool = {
  standard: {
    length: 0,
    A: [],
    B: [],
    C: [],
    D: []
  },
  trial: []
}

const Profiles = {
  A: 'A',
  B: 'B',
  C: 'C',
  D: 'D'
}

function matchProfile(profile) {
  switch (profile) {
    case Profiles.A:
      return Profiles.A
    case Profiles.B:
      return Profiles.C
    case Profiles.C:
      return Profiles.B
    case Profiles.D:
      return Profiles.A
  }
}

function approx(a, b) {
  return a - 10 < b && b < a + 10
}

function profile(cpu, memory) {
  if (cpu.cpu < 20 && memory.memory < 20) {
    return Profiles.A
  }
  if ((cpu.cpu > 50 && memory.memory > 50) || cpu.cpu > 80 || memory.memory > 80) {
    return Profiles.D
  }
  if (approx(cpu.cpu, memory.memory)) {
    if (cpu.tick < memory.tick) {
      return Profiles.B
    } else {
      return Profiles.C
    }
  } else {
    return cpu.cpu > memory.memory ? Profiles.B : Profiles.C
  }
}

function profileOrder(profile) {
  switch (profile) {
    case Profiles.A:
      return [Profiles.A, Profiles.B, Profiles.C, Profiles.D]
    case Profiles.B:
      return [Profiles.C, Profiles.A, Profiles.B, Profiles.D]
    case Profiles.C:
      return [Profiles.B, Profiles.A, Profiles.C, Profiles.D]
    case Profiles.D:
      return [Profiles.A, Profiles.B, Profiles.C, Profiles.D]
  }
  return [Profiles.A, Profiles.B, Profiles.C, Profiles.D]
}
let registeredFunctions = []


app.post('/register', (req, res) => {
  var uuid = crypto.randomUUID()
  console.log('New container registered with id ' + uuid)
  workers.push({
    uuid: uuid,
    memory: { memory: -1, tick: null },
    cpu: { cpu: -1, tick: null },
    requests: [],
    locked: false,
    no_parallel: false,
    clock: 0,
    cores: 0,
    stats: []
  })
  res.send(JSON.stringify(uuid))
})

app.post('/poll', (req, res) => {
  console.log('Poll request received')

  let stats = req.body
  let worker = workers.find((w) => w.uuid === stats.uuid)
  if (!worker) {
    console.log('Cannot find worker ' + stats.uuid)
    return res.send({ id: '-1' })
  }
  worker.cores = stats.cores
  worker.stats.push({ cpu: stats.cpu, memory: stats.memory, req: worker.requests.length, tick: worker.clock })
  worker.clock++
  if (stats.cpu > worker.cpu.cpu) {
    worker.cpu.cpu = stats.cpu
    worker.cpu.tick = worker.clock
  }

  if (stats.memory > worker.memory.memory) {
    worker.memory.memory = stats.memory
    worker.memory.tick = worker.clock
  }

  if (worker.requests.length === 0 && pool.trial.length > 0) {
    let request = pool.trial.shift()
    console.log('Employing worker ' + worker.uuid + ' to profile function ' + request.type)
    let response = {
      id: request.id,
      type: request.type,
      args: request.args,
      lock: request.lock
    }
    worker.locked = true
    worker.requests.push(request)
    return res.send(response)
  }
  let pr = profile(worker.cpu, worker.memory)

  if (pool.standard.length !== 0 && !worker.locked && !worker.no_parallel && worker.clock < 35 && worker.requests.length < worker.cores - 1) {
    let request = null
    order = profileOrder(pr)
    for (let prof of order) {
      if (pool.standard[prof].length !== 0) {
        request = pool.standard[prof].shift()
        pool.standard.length--
        break
      }
    }
    if (request !== null) {
      console.log('Sending task id ' + request.id + 'to worker ' + worker.uuid)
      let response = {
        id: request.id,
        type: request.type,
        args: request.args,
        lock: request.lock
      }
      worker.requests.push(request)
      if (no_parallel) {
        worker.no_parallel = true
      }
      return res.send(response)
    }
  }
  if (worker.requests.length === 0) {
    global_stats.push(worker.stats)
    console.log('Signaling worker ' + worker.uuid + ' to shutdown')
    return res.send({ id: '-1', type: worker.clock.toString() })
  } else {
    return res.send({ id: '0' })
  }
})

var nResults = 0

app.post('/sendResult/:reqID', (req, res) => {
  let response = req.body
  let id = req.params.reqID
  console.log('Received results of task id ' + id + 'from worker ' + response.id)
  res.end()

  let w = workers.find(k => k.uuid === response.id)
  if (w) {
    let i = w.requests.findIndex(k => k.id === id)
    let r = w.requests[i]
    if (w.locked) {
      let p = profile(w.cpu, w.memory)
      let k = registeredFunctions.find(k => k.id === r.type)
      k.profile = p
      console.log('Function ' + k.id + ' has been profiled as ' + p)
    }
    w.requests.splice(i, 1)
    r.response(JSON.stringify(response.data))
    nResults++
    console.log(nResults + ' have been fulfilled so far')
  } else {
    console.log('BIG ERROR: REQUEST MISSING: ' + id)
    console.log(workers)
  }
})

var workerCount = 0
var global_stats = []
app.post('/count', (req, res) => {
  workerCount = 0
  global_stats = []
  res.send()
})


app.get('/count', (req, res) => {
  res.send({ workers: workerCount, stats: global_stats })
})

app.post('/invoke', async (req, res) => {
  let type = req.body.id
  let request = {
    id: crypto.randomUUID(),
    type: type,
    args: req.body.args,
    response: (data) => {
      res.send(data)
    },
    lock: false
  }
  console.log('Received request of type ' + type + ' with id ' + request.id)

  let k = registeredFunctions.find(k => k.id === type)
  if (!k) {
    console.log('New function found, registering type ' + type)
    registeredFunctions.push({
      id: type,
      request: request,
      profile: null
    })
    request.lock = true
  } else {
    request.profile = k.profile
  }

  registerRequest(request)
})

async function launchWorker() {
  console.log('Launching new worker')
  workerCount++
  let res
  try {
    res = await axios.post(config.get('server.trigger'), { 'url': url })
    console.log('Worker ' + res.data + ' shutting down')
    let i = workers.findIndex(k => k.uuid === res.data)
    workers.splice(i, 1)
  } catch (e) {

  }

}

async function registerRequest(request) {
  if (request.lock) {
    pool.trial.push(request)
  } else {
    pool.standard[request.profile].push(request)
    pool.standard.length++
  }
  if (workers.length === 0 || request.lock || pool.standard.length > 3 * workers.length) {
    launchWorker()
  }
}

setInterval(function () {
  if ((workers.length === 0 && pool.standard.length !== 0) || pool.standard.length > 3 * workers.length) {
    launchWorker()
  }
}, 500)

app.listen(port, () => {
  console.log(`Coordinator listening on ${url}, no_parallel mode = ${no_parallel}`)
})
