import time

count = 0

while count < 1000:
    print(count, '\n', count, flush=True)
    count += 1
    time.sleep(1)