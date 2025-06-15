import json
import csv
import numpy as np
import matplotlib.pyplot as plt

# Файлы
json_file = "metrics.json"
csv_output = "user_vectors.csv"
points_output = "user_points.csv"
image_output = "embedding.png"

# Видеопул: videoID -> индекс
video_index_map = {f"typeA_{i}": i for i in range(10)}
video_index_map.update({f"typeB_{i}": 10 + i for i in range(10)})

# Чтение метрик
with open(json_file, "r") as f:
    metrics = json.load(f)

# userID -> вектор (20 элементов)
user_vectors = {}

for metric in metrics:
    user = metric["userID"]
    video = metric["videoID"]
    value = metric["value"]

    if video not in video_index_map:
        continue

    index = video_index_map[video]

    if user not in user_vectors:
        user_vectors[user] = [0.0] * 20

    user_vectors[user][index] += value

# CSV: полный вектор + rmse
with open(csv_output, "w", newline="") as csvfile, open(points_output, "w", newline="") as pointfile:
    writer = csv.writer(csvfile)
    point_writer = csv.writer(pointfile)

    header = [f"v{i}" for i in range(20)] + ["x_rmse", "y_rmse"]
    writer.writerow(header)
    point_writer.writerow(["user_id", "x", "y"])

    scatter_points = []

    for user, vector in user_vectors.items():
        v = np.array(vector)
        x = np.sqrt(np.mean(v[0:10] ** 2))
        y = np.sqrt(np.mean(v[10:20] ** 2))
        scatter_points.append((x, y))
        writer.writerow(list(v) + [x, y])
        point_writer.writerow([user, x, y])

# Отрисовка точек
xs, ys = zip(*scatter_points)
plt.figure(figsize=(8, 6))
plt.scatter(xs, ys, alpha=0.7, color="blue", edgecolors="k")
plt.xlabel("RMSE typeA (0–9)")
plt.ylabel("RMSE typeB (10–19)")
plt.title("User Metric Embedding (2D RMSE Projection)")
plt.grid(True)
plt.tight_layout()
plt.savefig(image_output, dpi=300)
plt.close()

print("✅ Готово: user_vectors.csv, user_points.csv, embedding.png")
