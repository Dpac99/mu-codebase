import numpy as np
import matplotlib.pyplot as plt
import requests
from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import as_completed
import time
from datetime import datetime

function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'
baseline_url = 'https://dlogbc5lu6e4nkbuu6uycob64i0kklzf.lambda-url.eu-west-3.on.aws/'
control_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/count'

dim = 4096
n_req = 256

output = "Test for {} with {} vectors:\n\tTotal time:\t{}\n\tAverage Time:\t{}\n\tMin Time:\t{}\n\tMax Time:\t{}\n"

matrix1 = np.random.randint(1000, size=(dim, dim))
matrix2 = np.random.randint(1000, size=(dim, dim))
matrix2_trans = matrix2.transpose()


def request_baseline(m1, m2):

    start = time.time()
    requests.post(
        baseline_url, json={"a": m1.tolist(), "b": m2.tolist()})
    end = time.time()

    return end-start


def request_func(m1, m2):
    start = time.time()
    body = {"id": "matrix", "args": {"a": m1.tolist(), "b": m2.tolist()}}
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


def runSequential(n):
    start = time.time()
    c = []
    for i in range(n):
        row = []
        for j in range(n):
            cell = 0
            for k in range(n):
                cell += matrix1[i][k] * matrix2[k][j]
            row.append(cell)
        c.append(row)
    end = time.time()
    return end-start


def run(n, func):
    print("Starting test for {} with {}".format(func, n))

    with ThreadPoolExecutor() as pool:
        m_size = n//n_req
        start = time.time()

        futures = []
        results = []

        futures = [pool.submit(func,
                               matrix1[i*m_size:(i+1)*m_size],
                               matrix2_trans[i*m_size:(i+1)*m_size].transpose())
                   for i in range(n_req)]
        results = [future.result() for future in as_completed(futures)]
        end = time.time()
        stats = analyze(results)
        print(output.format(func, n, end-start,
              stats[0], stats[1], stats[2]), flush=True)
        return [stats, end-start, ]


# sequentialTimes = []

# sequentialTimes.append(runSequential(256))
# sequentialTimes.append(runSequential(512))
# sequentialTimes.append(runSequential(1024))
# sequentialTimes.append(runSequential(2048))
# sequentialTimes.append(runSequential(4096))
# sequentialTimes.append(runSequential(8192))

values = [256, 512, 1024, 2048, 4096]

# baselineStats = []

# for value in values:
#     baselineStats.append(run(value, request_baseline))


# baselineTotals = [x[1] for x in baselineStats]
# baselineAverages = [x[0][0] for x in baselineStats]
# baselineMins = [x[0][1] for x in baselineStats]
# baselineMaxs = [x[0][2] for x in baselineStats]

functionStats = []
functionInvocations = []

for value in values:
    requests.post(control_url)
    functionStats.append(run(value, request_func))
    functionInvocations.append(int(requests.get(control_url).text))


functionTotals = [x[1] for x in functionStats]
functionAverages = [x[0][0] for x in functionStats]
functionMins = [x[0][1] for x in functionStats]
functionMaxs = [x[0][2] for x in functionStats]


figure, axis = plt.subplots(2, 3)


# axis[0, 0].plot(values, sequentialTimes, label="sequential")
# axis[0, 0].plot(values, baselineTotals, label="baseline")
axis[0, 0].plot(values, functionTotals, label="solution")
axis[0, 0].set_title("Total run time")

# axis[0, 1].plot(values, baselineAverages, label="baseline")
axis[0, 1].plot(values, functionAverages, label="solution")
axis[0, 1].set_title("Average run time")

# axis[1, 0].plot(values, baselineMins, label="baseline")
axis[1, 0].plot(values, functionMins, label="solution")
axis[1, 0].set_title("Minimum run times")

# axis[1, 1].plot(values, baselineMaxs, label="baseline")
axis[1, 1].plot(values, functionMaxs, label="solution")
axis[1, 1].set_title("Maximum run times")

axis[1, 2].plot(values, values, lable="baseline")
axis[1, 2].plot(values, functionInvocations, label="solution")

now = datetime.now()
try:
    plt.savefig("{}.png".format(now.strftime("%m/%d/%Y, %H:%M:%S")))
except:
    plt.savefig("plots.png")
plt.show()
