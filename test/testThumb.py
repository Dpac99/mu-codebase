import numpy as np
import matplotlib.pyplot as plt
import requests
import sys
from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import as_completed
import time
from datetime import datetime 


function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'
baseline_url = 'https://uxyihtkanjsorrwibbqezaapxy0nazhw.lambda-url.eu-west-3.on.aws/'
control_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/count'
baseline_matrix_url = 'https://dlogbc5lu6e4nkbuu6uycob64i0kklzf.lambda-url.eu-west-3.on.aws/'

dim = 1024
n_req = 64


output = "Test for {} with {} vectors:\n\tTotal time:\t{}\n\tAverage Time:\t{}\n\tMin Time:\t{}\n\tMax Time:\t{}\n"

matrix1 = np.random.randint(1000, size=(dim, dim))
matrix2 = np.random.randint(1000, size=(dim, dim))
matrix2_trans = matrix2.transpose()


def request_baseline(**args):
    i = args["i"]
    start = time.time()
    requests.post(
        baseline_url, json={
    "input_bucket": "dpac-serverless-thesis",
    "input_key": "big_image.jpg",
    "output_bucket": "dpac-serverless-output",
    "output_key": "logo_ist_thumb_{}.jpg".format(str(i)),
    "width": 100,
    "height": 100,
})
    end = time.time()

    return end-start


def request_func(**args):
    i = args["i"]
    start = time.time()
    body = {"id": "thumbnail", "args": {
    "input_bucket": "dpac-serverless-thesis",
    "input_key": "big_image.jpg",
    "output_bucket": "dpac-serverless-output",
    "output_key": "logo_ist_thumb_{}.jpg".format(str(i)),
    "width": 100,
    "height": 100,
}}
    requests.post(
        function_url, json=body)
    end = time.time()

    return end-start


def analyze(results):
    min = results[0]
    max = results[0]
    avg = 0
    for res in results:
        if res < min:
            min = res
        elif res > max:
            max = res
        avg += res
    return [avg/len(results), min, max]




def request_baseline_matrix(**args):
    m1 = args["m1"]
    m2 = args["m2"]
    start = time.time()
    requests.post(
        baseline_matrix_url, json={"a": m1, "b": m2})
    end = time.time()

    return end-start


def request_func_matrix(**args):
    m1 = args["m1"]
    m2 = args["m2"]
    start = time.time()
    body = {"id": "matrix", "args": {"a": m1, "b": m2}}
    requests.post(
        function_url, json=body)
    end = time.time()

    return end-start


def run(n, func, **data):
    print("Starting scaling test for {} with {}".format(func, n))
    with ThreadPoolExecutor() as pool:
        start = time.time()

        futures = []
        results = []

        futures = [pool.submit(func, **data) for i in range(n)]
        results = [future.result() for future in as_completed(futures)]
        end = time.time()
        stats = analyze(results)
        print(output.format(func, n, end-start,
                stats[0], stats[1], stats[2]), flush=True)


f = open("./matrixStats/6144.json", "w")
run(100, request_baseline_matrix, m1=[matrix1[0].tolist()], m2=matrix2[0][np.newaxis].T.tolist())
requests.post(control_url)
run(100, request_func_matrix, m1=[matrix1[0].tolist()], m2=matrix2[0][np.newaxis].T.tolist())
stats = requests.get(control_url)
f.write(stats.text)
f.close()


