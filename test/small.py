import requests
import numpy as np
from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import as_completed
import time


m1 = np.random.rand(2, 2)
m2 = np.random.rand(2, 2)

dim = 64
function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'
matrix1 = np.random.randint(1000, size=(dim, dim))
matrix2 = np.random.randint(1000, size=(dim, dim))


def post(body, i):
    start = time.time()
    x = requests.post(function_url, json=body)
    end = time.time()
    print("Request {} done".format(i), flush=True)
    return end-start


post({"id": "matrix", "args": {
            "a": matrix1.tolist(), "b": matrix2.tolist()}}, -1)
post({"id": "thumbnail", "args": {
    "input_bucket": "dpac-serverless-thesis",
    "input_key": "big_image.jpg",
    "output_bucket": "dpac-serverless-output",
    "output_key": "big_small.jpg",
    "width": 100,
    "height": 100,
}}, -1)


