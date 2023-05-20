const axios = require('axios')

const function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'
const baseline_url = 'https://uxyihtkanjsorrwibbqezaapxy0nazhw.lambda-url.eu-west-3.on.aws/'



async function request_func(i) {
    let body = {
        "id": "thumbnail", "args": {
            "input_bucket": "dpac-serverless-thesis",
            "input_key": "big_image.jpg",
            "output_bucket": "dpac-serverless-output",
            "output_key": `logo_ist_thumb_${i}.jpg`,
            "width": 100,
            "height": 100,
        }
    }
    axios.post(function_url, body)

}

async function request_baseline(i) {
    let body = {
        "input_bucket": "dpac-serverless-thesis",
        "input_key": "big_image.jpg",
        "output_bucket": "dpac-serverless-output",
        "output_key": `logo_ist_thumb_${i}.jpg`,
        "width": 100,
        "height": 100,
    }
    axios.post(baseline_url, body)
}

start = Date.now()
for (let i = 0; i < 128; i++) {
    request_baseline(i)
}
end = Date.now() - start
console.log(`Time for baseline: ${end} ms`)


start = Date.now()
for (let i = 0; i < 128; i++) {
    request_func(i)
}
end = Date.now() - start
console.log(`Time for function: ${end} ms`)
