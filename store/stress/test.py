import datetime
import io
import random
import string
import threading
import time
import uuid
from collections import defaultdict
from concurrent.futures import ThreadPoolExecutor, as_completed

import matplotlib.pyplot as plt
import requests


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
# ------------------ Конфигурация ------------------
DURATION = 60  # секунд
NUM_THREADS = 10

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

# ------------------ Метрики ------------------
metrics = []
metrics_lock = threading.Lock()

def now_ms():
    return time.time() * 1000

# ------------------ Потоки ------------------

# --- write_worker ---

def write_worker():
    end = time.time() + DURATION
    while time.time() < end:
        target = random.choice(WRITE_TARGETS)
        url = target["url"]
        try:
            if target["type"] == "json":
                body = target["generate"]()
                r = requests.post(url, json=body, timeout=5)
            elif target["type"] == "form":
                files, headers = target["generate"]()
                r = requests.post(url, files=files, headers=headers, timeout=5)
            else:
                continue
            status = r.status_code
            ok = r.ok
        except Exception:
            status = None
            ok = False
        with metrics_lock:
            metrics.append({
                "type": "write",
                "time": now_ms(),
                "url": url,
                "success": ok,
                "status": status
            })


def read_worker():
    end = time.time() + DURATION
    while time.time() < end:
        target = random.choice(READ_TARGETS)
        url = target["url"]
        headers = target.get("headers", {})

        try:
            r = requests.get(url, headers=headers, timeout=5)
            status = r.status_code
            ok = r.ok
        except Exception:
            status = None
            ok = False
        with metrics_lock:
            metrics.append({
                "type": "read",
                "time": now_ms(),
                "url": url,
                "success": ok,
                "status": status
            })

# ------------------ Запуск ------------------
print(f"Запускаем нагрузочный тест на {DURATION} секунд...")

with ThreadPoolExecutor(max_workers=NUM_THREADS) as executor:
    futures = []
    futures.append(executor.submit(write_worker))
    for _ in range(NUM_THREADS - 1):
        futures.append(executor.submit(read_worker))
    for f in as_completed(futures):
        pass

print("Тест завершён. Обрабатываем метрики...")

# ------------------ График ------------------

def plot_success_error_timeline():
    # Группируем данные по секундам с 0 до DURATION-1
    buckets = defaultdict(lambda: {"read_success": 0, "read_error": 0, "write_success": 0, "write_error": 0})

    start_time_s = int(min(m["time"] for m in metrics) // 1000)  # начало измерений в секундах

    for m in metrics:
        t_sec = int(m["time"] // 1000) - start_time_s
        if 0 <= t_sec < DURATION:
            key = f"{m['type']}_{'success' if m['success'] else 'error'}"
            buckets[t_sec][key] += 1

    seconds = list(range(DURATION))
    read_success = [buckets[s]["read_success"] for s in seconds]
    read_error   = [buckets[s]["read_error"] for s in seconds]
    write_success= [buckets[s]["write_success"] for s in seconds]
    write_error  = [buckets[s]["write_error"] for s in seconds]

    plt.figure(figsize=(14, 6))
    plt.plot(seconds, read_success, label="Read success", color="blue")
    plt.plot(seconds, read_error, label="Read error", color="lightblue")
    plt.plot(seconds, write_success, label="Write success", color="green")
    plt.plot(seconds, write_error, label="Write error", color="red")

    plt.xlabel("Время (секунды)")
    plt.ylabel("Количество запросов")
    plt.title("Динамика запросов по секундам")
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.savefig("result.png", dpi=300)
    print("График сохранён в result.png")

plot_success_error_timeline()
