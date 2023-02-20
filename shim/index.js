const express = require('express')
const bodyParser = require('body-parser');
const crypto = require('crypto');
var Docker = require('dockerode')


const app = express()
app.use(bodyParser.json({ extended: true }));
const port = 1234

var docker = new Docker();
let image

async function loadImages() {
  image = await docker.listImages().then((data) => {
    return data.find((image) => image.RepoTags[0] === 'test_function:latest')
  })
  console.log(image)
}


let workers = []
let pool = []

app.get('/', (req, res) => {
  res.send('Hello World2!')
})

app.post('/register', (req, res) => {
  var uuid = crypto.randomUUID()
  console.log(uuid)
  console.log(req.body)
  workers.push({
    "UUID": uuid,
    "Memory": null,
    "CPU": null,
    "nTasks": 1,
  })
  console.log(workers)
  res.send(JSON.stringify(uuid))
})

app.post('/poll', (req, res) => {
  let stats = req.body
  let index = workers.findIndex((w) => w.UUID === stats.UUID)
  workers[index].Memory = stats.Memory
  workers[index].CPU = stats.CPU
})

app.post('/invoke', (req, res) => {

})

app.listen(port, () => {
  loadImages()
  console.log(`Example app listening on port ${port}`)
})
