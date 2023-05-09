import requests

function_url = "http://localhost:9000/2015-03-31/functions/function/invocations"


def post(body):
    x = requests.post(function_url, json=body)
    print("Request done", flush=True)


post({
    "input_bucket": "dpac-serverless-thesis",
    "input_key": "logo_ist.jpg",
    "output_bucket": "dpac-serverless-output",
    "output_key": "logo_ist_thumb.jpg",
    "width": 100,
    "height": 100,
})
