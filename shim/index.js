const express = require('express')
const app = express()
const port = 1234
const crypto = require('crypto');
var Docker = require('dockerode')
var docker = new Docker();
let image

async function loadImages() {
  image = await docker.listImages().then((data) => {
    return data.find((image) => image.RepoTags[0] === 'test_function:latest')
  })
  console.log(image)
}


let workers = []

app.get('/', (req, res) => {
  res.send('Hello World2!')
})

app.post('/register', (req, res) => {
  var ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress
  var uuid = crypto.randomUUID()
  console.log(uuid)
  workers.push(ip)
  console.log(workers)
  res.send(JSON.stringify(uuid))
})

app.post('/poll', (req, res) => {
  var ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress
  console.log(req)
  if (workers.includes(ip)) {
    console.log('received poll request')
    res.send('ok')
  }
})

app.listen(port, () => {
  loadImages()
  console.log(`Example app listening on port ${port}`)
})
