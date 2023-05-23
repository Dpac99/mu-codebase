const axios = require('axios')
const fs = require('fs')

const function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'
const baseline_url = 'https://uxyihtkanjsorrwibbqezaapxy0nazhw.lambda-url.eu-west-3.on.aws/'
const control_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/count'
const thread_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/threshold'




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
    return axios.post(function_url, body)

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
    return axios.post(baseline_url, body)
}

function sleep(ms) {
    return new Promise((resolve) => {
        setTimeout(resolve, ms);
    });
}

async function baseline() {
    let start = Date.now()
    let promises = []

    for (let k = 0; k < 5; k++) {
        for (let i = 0; i < 5; i++) {
            promises.push(request_baseline(i))
        }
        await sleep(1000)
    }

    for (let k = 0; k < 5; k++) {
        for (let i = 0; i < 20; i++) {
            promises.push(request_baseline(i))
        }
        await sleep(1000)
    }

    for (let k = 0; k < 5; k++) {
        for (let i = 0; i < 5; i++) {
            promises.push(request_baseline(i))
        }
        await sleep(1000)
    }
    await Promise.all(promises)
    let end = Date.now() - start
    console.log(`Time for baseline: ${end} ms`)

}

async function init() {
    await request_func('small')
}


async function solution() {
    // await init()
    // await axios.post(control_url, {})
    let ncores = Number(process.argv[2])
    await axios.post(thread_url, ncores)
    let start = Date.now()
    let promises = []

    for (let k = 0; k < 5; k++) {
        for (let i = 0; i < 5; i++) {
            promises.push(request_func(i))
        }
        await sleep(1000)
    }

    for (let k = 0; k < 5; k++) {
        for (let i = 0; i < 20; i++) {
            promises.push(request_func(i))
        }
        await sleep(1000)
    }

    for (let k = 0; k < 5; k++) {
        for (let i = 0; i < 5; i++) {
            promises.push(request_func(i))
        }
        await sleep(1000)
    }
    await Promise.all(promises)
    let end = Date.now() - start
    console.log(`Time for func: ${end} ms`)
    // let res = await axios.get(control_url, {})
    // fs.writeFile('stats2/10240_v2.json', JSON.stringify(res.data), function (err) {
    //     if (err) throw err;
    //     console.log('Saved!');
    // })
}



// baseline()
solution()