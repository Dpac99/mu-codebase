import numpy as np
import matplotlib.pyplot as plt
import json
import math

filename = "stats2/{}_v2.json"
config = ["4096", "6144", "8192", "10240"]
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



plt.style.use('_mpl-gallery')



for file in config:
    min_global_tick = -1
    max_global_tick = -1
    f = open(filename.format(file))
    data = json.load(f)
    f.close()

    x = np.linspace(min_global_tick, max_global_tick)


    fig, axis = plt.subplots(2,2)

    axis[0,0].plot(data["workerTrace"][0], data["workerTrace"][1], label = "Number of workers")
    axis[0,0].plot(data["requestTrace"][0], data["requestTrace"][1], label = "Number of requests")
    axis[0,0].set_title("Number of workers and requests")
    axis[0,0].set_xlabel("Time since first request (ms)")
    axis[0,0].set_ylabel("Total number")
    axis[0,0].legend()

    axis[0,1].set_title("Trace of workers CPU")
    axis[0,1].set_xlabel("Time since first request (ms)")
    axis[0,1].set_ylabel("CPU Used (%)")
    axis[0,1].legend()

    axis[1,0].set_title("Trace of workers memory")
    axis[1,0].set_xlabel("Time since first request (ms)")
    axis[1,0].set_ylabel("Memory Used (%)")
    axis[1,0].legend()

    axis[1,1].set_title("Trace of workers requests")
    axis[1,1].set_xlabel("Time since first request (ms)")
    axis[1,1].set_ylabel("Number of requests")
    axis[1,1].legend()


    
    n_workers = data["workers"]
    for worker in data["stats"]:
        ticks = []
        cpu = []
        memory = []
        req = []
        for t in worker:
            cpu.append(t["cpu"])
            memory.append(t["memory"])
            req.append(t["req"])
            k = t["global_tick"]
            if min_global_tick == -1 or k < min_global_tick:
                min_global_tick = k
            if k > max_global_tick:
                max_global_tick = k
            ticks.append(k)
        axis[0,1].plot(ticks, cpu)
        axis[1,0].plot(ticks, memory)
        axis[1,1].plot(ticks, req)
    plt.show()


    



# def test():
#     avgcpu = []
#     avgmem = []
#     avgreqs = []

#     for tick in range(max_ticks):
#         avg_cpu = 0
#         avg_mem = 0
#         avg_reqs = 0
#         for worker in range(len(cpus)):
#             if tick < len(cpus[worker]):
#                 avg_cpu += cpus[worker][tick]
#                 avg_mem += memories[worker][tick]
#                 avg_reqs += reqs[worker][tick]
#             else:
#                 avg_cpu += 0
#                 avg_mem += 0
#                 avg_reqs += 0
#         avg_cpu /= n_workers
#         avg_mem /= n_workers
#         avg_reqs =  math.ceil(avg_reqs / n_workers)
#         avgcpu.append(avg_cpu)
#         avgmem.append(avg_mem)
#         avgreqs.append(avg_reqs)

#  stats[file]["cpu"] = avgcpu
#     stats[file]["memory"] = avgmem
#     stats[file]["requests"] = avgreqs

