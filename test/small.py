import requests
import numpy as np
from concurrent.futures import ThreadPoolExecutor


m1 = np.random.rand(2, 2)
m2 = np.random.rand(2, 2)

dim = 64
function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'
matrix1 = np.random.randint(1000, size=(dim, dim))
matrix2 = np.random.randint(1000, size=(dim, dim))


def post(body, i):
    print("Request {} starting ...".format(i))
    x = requests.post(function_url, json=body)
    print("Request {} done".format(i))
    return x.text


post({"id": "matrix", "args": {
            "a": matrix1.tolist(), "b": matrix2.tolist()}}, -1)

futures = []


with ThreadPoolExecutor() as pool:
    for i in range(15):
        futures.append(pool.submit(post, {"id": "matrix", "args": {
            "a": matrix1.tolist(), "b": matrix2.tolist()}}, i))
    for i in range(15):
        futures[i].result()
