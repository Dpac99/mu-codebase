const express = require('express')
const bodyParser = require('body-parser')
const crypto = require('crypto')
const axios = require('axios')
const config = require('config')

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

function matchProfile (profile) {
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

function approx (a, b) {
  return a - 10 < b && b < a + 10
}

function profile (cpu, memory) {
  if (cpu.cpu < 20 && memory.memory < 20) {
    return Profiles.A
  }
  if (cpu.cpu > 50 && memory.memory > 50) {
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

let registeredFunctions = []

app.get('/', (req, res) => {
  res.send('Hello World2!')
})

app.post('/register', (req, res) => {
  var uuid = crypto.randomUUID()
  console.log('New container registered with id ' + uuid)
  workers.push({
    uuid: uuid,
    memory: { memory: -1, tick: null },
    cpu: { cpu: -1, tick: null },
    requests: [],
    locked: false,
    clock: 0
  })
  res.send(JSON.stringify(uuid))
})

app.post('/poll', (req, res) => {
  console.log('Poll request received')

  let stats = req.body
  console.log(stats)
  let worker = workers.find((w) => w.uuid === stats.uuid)
  if (!worker) {
    return res.send({ id: '-1' })
  }
  worker.clock++
  console.log('Clock = ' + worker.clock)
  if (stats.cpu > worker.cpu.cpu) {
    worker.cpu.cpu = stats.cpu
    worker.cpu.tick = worker.clock
  }

  if (stats.memory > worker.memory.memory) {
    worker.memory.memory = stats.memory
    worker.memory.tick = worker.clock
  }
  console.log('Nworkers: ' + workers.length)
  console.log('Ntasks: ' + pool.standard.length)

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

  if (pool.standard.length !== 0 && !worker.locked && worker.clock < 12) {
    let request = null
    if (worker.requests.length === 0 && pool.standard.D.length !== 0) {
      request = pool.standard.D.shift()
      pool.standard.length--
    } else {
      let pr = profile(worker.cpu, worker.memory)
      console.log('Worker ' + worker.uuid + ' has profile ' + pr)
      if (pool.standard[pr].length !== 0) {
        request = pool.standard[pr].shift()
        pool.standard.length--
      } else if (pool.standard[matchProfile(pr)].length !== 0) {
        request = pool.standard[matchProfile(pr)].shift()
        pool.standard.length--
      } else if (pr !== Profiles.A && pool.standard[Profiles.A].length !== 0) {
        request = pool.standard[Profiles.A].shift()
        pool.standard.length--
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
      return res.send(response)
    }
  }
  if (worker.requests.length === 0) {
    console.log('Signaling worker ' + worker.uuid + ' to shutdown')
    return res.send({ id: '-1' })
  } else {
    return res.send({ id: '0' })
  }
})

app.post('/sendResult/:reqID', (req, res) => {
  let response = req.body
  let id = req.params.reqID
  console.log('Received results of task id ' + id + 'from worker ' + response.id)
  res.end()

  let w = workers.find(k => k.uuid === response.id)
  if (w) {
    let i = w.requests.findIndex(k => k.id === id)
    let r = w.requests[i]
    console.log(w)
    console.log(i)
    if (w.locked) {
      let p = profile(w.cpu, w.memory)
      let k = registeredFunctions.find(k => k.id === r.type)
      k.profile = p
      console.log('Function ' + k.id + ' has been profiled as ' + p)
    }
    w.requests.splice(i, 1)
    r.response(JSON.stringify(response.data))
  }
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

async function registerRequest (request) {
  if (request.lock) {
    pool.trial.push(request)
  } else {
    pool.standard[request.profile].push(request)
    pool.standard.length++
  }
  console.log('Nworkers: ' + workers.length)
  console.log('Ntasks: ' + pool.standard.length)
  if (workers.length === 0 || request.lock || pool.standard.length - 5 > workers.length) {
    console.log('Launching new worker')
    let res
    try {
      res = await axios.post(config.get('server.trigger'), { 'url': url })
    } catch (e) {
      console.log(e)
    }

    console.log('Worker ' + res.data + ' shutting down')
    let i = workers.findIndex(k => k.id === res.data)
    workers.splice(i, 1)
  }
}

app.listen(port, () => {
  console.log(`Coordinator listening on ${url}`)
})
