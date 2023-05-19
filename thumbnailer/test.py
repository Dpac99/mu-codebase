import requests

function_url = "https://uxyihtkanjsorrwibbqezaapxy0nazhw.lambda-url.eu-west-3.on.aws/"


def post(body):
    x = requests.post(function_url, json=body)
    print("Request done", flush=True)


post({
    "input_bucket": "dpac-serverless-thesis",
    "input_key": "logo_ist.jpg",
    "output_bucket": "dpac-serverless-output",
    "output_key": "logo_ist_thumb_base.jpg",
    "width": 100,
    "height": 100,
})
