import io
import string
import time
import random
import threading
import uuid

import requests
import matplotlib.pyplot as plt
from collections import defaultdict

DURATION = 60  # секунд
RATE_PER_SECOND = 10  # всего запросов в секунду
WRITE_RATIO = 0.1     # 1 из 10 запросов на запись
READ_RATIO = 0.9      # 9 из 10 на чтение

metrics_lock = threading.Lock()
metrics = []

# Пример read targets с заголовками
READ_TARGETS = [
    {
        "url": "http://localhost:8080/authors/681f98e307ea0ea1385f3865"
    },
    {
        "url": "http://localhost:8080/file/681faa25e56ffb08a70ae230/681faa25e56ffb08a70ae231",
        "headers": {"user_id": "681f98e307ea0ea1385f3865", "Accept": "multipart/form-data"}
    },
    {
        "url": "GET http://localhost:8080/meta/681fd5c39725629cd0b4e797",
    },
    {
        "url": "http://localhost:8080/preview/681fd5c39725629cd0b4e797",
    },
    {
        "url": "http://localhost:8080/search/videos?q=test",
    }
]

# Пример write targets
def random_username():
    return "user_" + ''.join(random.choices(string.ascii_lowercase, k=6))

def random_size():
    return {"sizeInBytes": random.randint(512, 4096)}

def random_user_id():
    return uuid.uuid4().hex

def random_file_upload():
    file_id = uuid.uuid4().hex
    file_name = "rand_" + ''.join(random.choices(string.ascii_lowercase, k=5)) + ".txt"
    file_content = io.BytesIO(random.randbytes(random.randint(100, 1000)))

    files = {
        file_id: (file_name, file_content)
    }

    headers = {
        "user_id": random_user_id()
    }

    return files, headers

WRITE_TARGETS = [
    {
        "url": "http://localhost:8080/authors/create",
        "type": "json",
        "generate": lambda: {"username": random_username()}
    },
    {
        "url": "http://localhost:8080/file/write/plan",
        "type": "json",
        "generate": random_size
    },
    {
        "url": "http://localhost:8080/file/write/" + uuid.uuid4().hex,
        "type": "form",
        "generate": random_file_upload
    }
]

def now_ms():
    return int(time.time() * 1000)

def do_read():
    target = random.choice(READ_TARGETS)
    url = target["url"]
    headers = target.get("headers", {})
    try:
        r = requests.get(url, headers=headers, timeout=5)
        success = r.ok
    except Exception:
        success = False
    with metrics_lock:
        metrics.append({"type": "read", "time": now_ms(), "success": success})

def do_write():
    target = random.choice(WRITE_TARGETS)
    url = target["url"]
    try:
        if target["type"] == "json":
            body = target["generate"]()
            r = requests.post(url, json=body, timeout=5)
        else:
            r = None
        success = r.ok if r else False
    except Exception:
        success = False
    with metrics_lock:
        metrics.append({"type": "write", "time": now_ms(), "success": success})

def main_loop():
    start = time.time()
    end = start + DURATION
    interval = 1.0 / RATE_PER_SECOND  # время между запросами

    while time.time() < end:
        # Выбираем тип запроса с вероятностью
        t = random.random()
        if t < WRITE_RATIO:
            threading.Thread(target=do_write).start()
        else:
            threading.Thread(target=do_read).start()
        time.sleep(interval)

def plot_errors_over_time():
    if not metrics:
        print("Нет данных для графика")
        return
    start_time = min(m["time"] for m in metrics)
    counts = defaultdict(lambda: {
        "read_total":0, "read_errors":0,
        "write_total":0, "write_errors":0
    })
    for m in metrics:
        sec = (m["time"] - start_time) // 1000
        if m["type"] == "read":
            counts[sec]["read_total"] += 1
            if not m["success"]:
                counts[sec]["read_errors"] += 1
        else:
            counts[sec]["write_total"] += 1
            if not m["success"]:
                counts[sec]["write_errors"] += 1

    seconds = sorted(counts.keys())
    read_total_cum = []
    read_errors_cum = []
    write_total_cum = []
    write_errors_cum = []

    rtotal = rerror = wtotal = werror = 0
    for s in seconds:
        rtotal += counts[s]["read_total"]
        rerror += counts[s]["read_errors"]
        wtotal += counts[s]["write_total"]
        werror += counts[s]["write_errors"]

        read_total_cum.append(rtotal)
        read_errors_cum.append(rerror)
        write_total_cum.append(wtotal)
        write_errors_cum.append(werror)

    plt.figure(figsize=(12,6))
    plt.plot(seconds, read_total_cum, label="Всего чтение", color="blue")
    plt.plot(seconds, read_errors_cum, label="Ошибки чтения", color="lightblue")
    plt.plot(seconds, write_total_cum, label="Всего запись", color="green")
    plt.plot(seconds, write_errors_cum, label="Ошибки записи", color="red")

    plt.xlabel("Время, секунды")
    plt.ylabel("Накопленное число запросов")
    plt.title("Накопление ошибок и общего числа запросов по типам")
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.savefig("result_sum.png", dpi=300)
    print("График сохранён в result_sum.png")

if __name__ == "__main__":
    print("Старт теста...")
    main_loop()
    # Немного ждем пока все потоки завершатся
    time.sleep(2)
    print("Построение графика...")
    plot_errors_over_time()
    print("Готово.")
