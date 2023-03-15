const express = require('express')
const bodyParser = require('body-parser')
const crypto = require('crypto')
const axios = require('axios')
// var Docker = require('dockerode')

const app = express()
app.use(bodyParser.json({ extended: true }))
const port = 1234

// var docker = new Docker()
// let image

// async function loadImages () {
//   let image = await docker.listImages().then((data) => {
//     return data.find((image) => {
//       if (image && image.RepoTags) {
//         return image.RepoTags[0] === 'test_function:latest'
//       }
//     })
//   })
// }

let workers = []

app.get('/', (req, res) => {
  res.send('Hello World2!')
})

app.post('/register', (req, res) => {
  var uuid = crypto.randomUUID()
  workers.push({
    'UUID': uuid,
    'Memory': null,
    'CPU': null,
    'nTasks': 1
  })
  res.send(JSON.stringify(uuid))
})

app.post('/poll', async (req, res) => {
  let stats = req.body
  let index = workers.findIndex((w) => w.UUID === stats.UUID)
  workers[index].Memory = stats.Memory
  workers[index].CPU = stats.CPU

  console.log(workers[index])

  res.send()
})

app.post('/invoke', async (req, res) => {
  await axios.post('http://localhost:9000/2015-03-31/functions/function/invocations', req.body)

  res.send()
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
