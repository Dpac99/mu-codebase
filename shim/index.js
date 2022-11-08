const express = require('express')
const app = express()
const port = 1234

let workers = []

app.get('/', (req, res) => {
  res.send('Hello World2!')
})

app.post('/register', (req, res) => {
  var ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress
  workers.push(ip)
  console.log(workers)
  res.send()
})

app.get('/poll', (req, res) => {
  var ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress

  if (workers.includes(ip)) {
    console.log('received poll request')
    res.send('ok')
  }
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
