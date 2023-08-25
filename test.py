import time

count = 0

while count < 1000:
    print(count, flush=True)
    count += 1
    time.sleep(1)
    if count == 10:
        raise Exception('Test exception')