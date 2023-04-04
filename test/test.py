import numpy as np
import matplotlib as plt
import docker as docker

url = 'http://localhost:1234/invoke'
image = 'test_function'


client = docker.from_env()
containers = []

for i in range(100):
    containers[i] = client.containers.run(image, detach=True, extra-hosts={'host.docker.internal': 'host-gateway'}, ports: {'9000': '8080'})
