const express = require('express')
const bodyParser = require('body-parser')
const crypto = require('crypto')
const axios = require('axios')
const config = require('config')
const { time } = require('console')
const args = require('yargs').argv;


const no_parallel = args.no_parallel ?? false
const dynamic_spawn = args.dynamic_spawn ?? false
const no_limit = args.no_limit ?? false
const app = express()
app.use(bodyParser.json({ limit: '50mb', extended: true }))

const port = config.get('server.port')
const url = config.get('server.host')

let workers = []
let globalWorkerCount = 0
let avgThreads = 1
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

function canHandleTask(worker) {
  return pool.standard.length !== 0
    && !worker.locked
    && !worker.no_parallel
    && worker.clock < 35
    && (no_limit ? true : worker.requests.length < worker.cores - 1) //if no limit, dont enforce one task per core
    && worker.cpu.cpu < 95
}

app.post('/poll', (req, res) => {
  if (start_tick !== null) {
    let tick = Date.now() - start_tick
    workerTrace[0].push(tick)
    workerTrace[1].push(globalWorkerCount)
    requestTrace[0].push(tick)
    requestTrace[1].push(pool.standard.length)
    avgThreshold[0].push(tick)
    avgThreshold[1].push(dynamic_spawn ? avgThreads : spawn_threshold)
  }

  let stats = req.body
  let worker = workers.find((w) => w.uuid === stats.uuid)
  if (!worker) {
    console.log('Cannot find worker ' + stats.uuid)
    return res.send({ id: '-1' })
  }
  // if first poll request, recalculate avgThreads
  if (worker.cores === 0) {
    avgThreads = Math.floor(((avgThreads * workers.length - 1) + stats.cores) / workers.length)
  }

  worker.cores = stats.cores
  if (worker.clock === 0 && start_tick !== null) {
    worker.start_tick = Date.now() - start_tick
  }

  //if one task per thread limit is removed, once a function goes over the limit then recalculate the avgThreads
  if (no_limit && worker.requests.length > worker.cores - 1) {
    avgThreads = Math.floor(((avgThreads * workers.length - 1) + worker.requests.length) / workers.length)

  }
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

  if (canHandleTask(worker)) {
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
      worker.stats.push({
        cpu: stats.cpu,
        memory: stats.memory,
        req: worker.requests.length,
        tick: worker.clock,
        res: "req",
        cores: worker.cores,
        global_tick: start_tick !== null ? Date.now() - start_tick : -1
      })
      return res.send(response)

    }
  }
  if (worker.requests.length === 0) {
    global_stats.push(worker.stats)
    console.log('Signaling worker ' + worker.uuid + ' to shutdown')
    worker.stats.push({
      cpu: stats.cpu,
      memory: stats.memory,
      req: worker.requests.length,
      tick: worker.clock,
      res: "shutdown",
      cores: worker.cores,
      global_tick: start_tick !== null ? Date.now() - start_tick : -1
    })

    return res.send({ id: '-1', type: worker.clock.toString() })
  } else {
    worker.stats.push({
      cpu: stats.cpu,
      memory: stats.memory,
      req: worker.requests.length,
      tick: worker.clock,
      res: "resume",
      cores: worker.cores,
      global_tick: start_tick !== null ? Date.now() - start_tick : -1
    })

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
var workerTrace = [[], []]
var requestTrace = [[], []]
var avgThreshold = [[], []]
var global_stats = []
var start_tick = null
app.post('/count', (req, res) => {
  workerCount = 0
  global_stats = []
  workerTrace = [[], []]
  requestTrace = [[], []]
  avgThreshold = [[], []]
  start_tick = Date.now()
  res.send()
})


app.get('/count', (req, res) => {
  res.send({ workers: workerCount, stats: global_stats, workerTrace: workerTrace, requestTrace: requestTrace, thresholdTrace: avgThreshold })
})

var spawn_threshold = 4

app.post('/threshold', (req, res) => {
  spawn_threshold = Number(req.body)
  console.log('Function threshold set to', spawn_threshold)
  res.send()
})

app.post('/resetThreshold', (req, res) => {
  avgThreads = 1
  res.send()
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
  console.log('Launching new worker, workers = ', globalWorkerCount, ' task = ', pool.standard.length)
  workerCount++
  globalWorkerCount++
  let res
  try {
    res = await axios.post(config.get('server.trigger'), { 'url': url })
    console.log('Worker ' + res.data + ' shutting down')
    let i = workers.findIndex(k => k.uuid === res.data)
    workers.splice(i, 1)
    globalWorkerCount--
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
  if ((globalWorkerCount === 0 || request.lock || pool.standard.length > (dynamic_spawn ? avgThreads : spawn_threshold) * globalWorkerCount) && globalWorkerCount < 50) {
    launchWorker()
  }
}

setInterval(function () {
  if (((workers.globalWorkerCount === 0 && pool.standard.length !== 0)
    || pool.standard.length > (dynamic_spawn ? avgThreads : spawn_threshold) * globalWorkerCount)
    && globalWorkerCount < 50) {
    launchWorker()
  }
}, 500)

app.listen(port, () => {
  console.log(`Coordinator listening on ${url}, no_parallel = ${no_parallel}, dynamic_spawn = ${dynamic_spawn}, no_limit = ${no_limit}`)
})
