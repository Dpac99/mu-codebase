import numpy as np
import matplotlib.pyplot as plt
import json
import math

filename = "stats2/{}.json"
config = ["2048","4096","6144","8192","10240"]
stats = {
    "2048":{
        "cpu": [],
        "memory": [],
        "requests": [],
    },
    "4096":{
        "cpu": [],
        "memory": [],
        "requests": [],
    },
    "6144":{
        "cpu": [],
        "memory": [],
        "requests": [],
    },
    "8192":{
        "cpu": [],
        "memory": [],
        "requests": [],
    },
    "10240":{
        "cpu": [],
        "memory": [],
        "requests": [],
    },
}
max_ticks = 0
for file in config:
    f = open(filename.format(file))
    data = json.load(f)
    f.close()
    cpus = []
    memories = []
    reqs = []
    n_workers = data["workers"]
    for worker in data["stats"]:
        if len(worker) > max_ticks:
            max_ticks = len(worker)
        cpu = []
        memory = []
        req = []
        for tick in worker:
            cpu.append(tick["cpu"])
            memory.append(tick["memory"])
            req.append(tick["req"])
        cpus.append(cpu)
        memories.append(memory)
        reqs.append(req)

    avgcpu = []
    avgmem = []
    avgreqs = []

    for tick in range(max_ticks):
        avg_cpu = 0
        avg_mem = 0
        avg_reqs = 0
        for worker in range(len(cpus)):
            if tick < len(cpus[worker]):
                avg_cpu += cpus[worker][tick]
                avg_mem += memories[worker][tick]
                avg_reqs += reqs[worker][tick]
            else:
                avg_cpu += 0
                avg_mem += 0
                avg_reqs += 0
        avg_cpu /= n_workers
        avg_mem /= n_workers
        avg_reqs =  math.ceil(avg_reqs / n_workers)
        avgcpu.append(avg_cpu)
        avgmem.append(avg_mem)
        avgreqs.append(avg_reqs)

    stats[file]["cpu"] = avgcpu
    stats[file]["memory"] = avgmem
    stats[file]["requests"] = avgreqs

plt.style.use('_mpl-gallery')

fig, ax = plt.subplots()

ax.plot(range())



