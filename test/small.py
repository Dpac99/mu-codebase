import requests
import numpy as np


m1 = np.random.rand(2, 2)
m2 = np.random.rand(2, 2)

dim = 64

matrix1 = np.random.randint(1000, size=(dim, dim))
matrix2 = np.random.randint(1000, size=(dim, dim))

function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'

x = requests.post(
    function_url, json={"id": "matrix", "args": {"a": matrix1.tolist(), "b": matrix2.tolist()}})

print(x.json())
