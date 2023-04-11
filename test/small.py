import requests
import numpy as np


m1 = np.random.rand(2, 2)
m2 = np.random.rand(2, 2)

function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'

x = requests.post(
    function_url, json={"id": "matrix", "args": {"a": m1.tolist(), "b": m2.tolist()}})

print(x.json())
