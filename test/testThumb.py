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

dim = 1024
n_req = 64


output = "Test for {} with {} vectors:\n\tTotal time:\t{}\n\tAverage Time:\t{}\n\tMin Time:\t{}\n\tMax Time:\t{}\n"

matrix1 = np.random.randint(1000, size=(dim, dim))
matrix2 = np.random.randint(1000, size=(dim, dim))
matrix2_trans = matrix2.transpose()


def request_baseline(m1, m2, i):

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


def request_func(m1, m2, i):
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



def run(n, func):
    print("Starting scaling test for {} with {}".format(func, n))
    with ThreadPoolExecutor() as pool:
        start = time.time()

        futures = []
        results = []

        futures = [pool.submit(func,
                                matrix1[i],
                                matrix2_trans[i].transpose(), i)
                    for i in range(n)]
        results = [future.result() for future in as_completed(futures)]
        end = time.time()
        stats = analyze(results)
        print(output.format(func, n, end-start,
                stats[0], stats[1], stats[2]), flush=True)
        return [stats, end-start, ]

f = open("2048_v3.json", "w")
# run(100, request_baseline)
requests.post(control_url)
run(100, request_func)
stats = requests.get(control_url)
f.write(stats.text)
f.close()

# values = [64, 128, 256, 512, 1024]

# baselineStats = []

# for value in values:
#     baselineStats.append(run(value, request_baseline))


# baselineTotals = [x[1] for x in baselineStats]
# baselineAverages = [x[0][0] for x in baselineStats]
# baselineMins = [x[0][1] for x in baselineStats]
# baselineMaxs = [x[0][2] for x in baselineStats]

# functionStats = []
# functionInvocations = []

# for value in values:
#     requests.post(control_url)
#     functionStats.append(run(value, request_func))
#     functionInvocations.append(int(requests.get(control_url).text))


# functionTotals = [x[1] for x in functionStats]
# functionAverages = [x[0][0] for x in functionStats]
# functionMins = [x[0][1] for x in functionStats]
# functionMaxs = [x[0][2] for x in functionStats]


# figure, axis = plt.subplots(2, 2)


# # axis[0, 0].plot(values, sequentialTimes, label="sequential")
# axis[0, 0].plot(values, baselineTotals, label="baseline", marker='|')
# axis[0, 0].plot(values, functionTotals, label="solution", marker='|')
# axis[0, 0].set_title("Total run time")
# axis[0,0].set_ylabel("Time in Seconds")
# axis[0,0].set_xlabel("Size of data")
# axis[0,0].legend()

# axis[0, 1].errorbar(values, baselineAverages, yerr=[baselineMins, baselineMaxs], label="baseline", marker='|')
# axis[0, 1].errorbar(values, functionAverages, yerr=[functionMins, functionMaxs], label="solution", marker='|')
# axis[0, 1].set_title("Average run time")
# axis[0,1].set_ylabel("Time in Seconds")
# axis[0,1].set_xlabel("Size of data")
# axis[0,1].legend()


# axis[1, 0].plot(values, values, label="baseline", marker='|')
# axis[1, 0].plot(values, functionInvocations, label="solution", marker='|')
# axis[1, 0].set_title("Number of invocations")
# axis[1,0].set_ylabel("Lambdas Invoked")
# axis[1,0].set_xlabel("Size of data")
# axis[1,0].legend()

# now = datetime.now()
# try:
#     plt.savefig("../images/{}.png".format(now.strftime("%d-%m-%Y_%H:%M:%S")), dpi=200)
# except:
#     plt.savefig("plots.png")
# plt.show()
