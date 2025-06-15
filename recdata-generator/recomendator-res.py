import csv
import matplotlib.pyplot as plt

# Заменишь на свои кластеры
clusters = [
    ["user_28", "user_80", "user_82", "user_63", "user_1", "user_35", "user_24", "user_9", "user_18", "user_37", "user_10", "user_6"],
    ["user_39", "user_3", "user_44", "user_42", "user_66", "user_77", "user_15"],
    ["user_81", "user_75", "user_87", "user_58", "user_97", "user_98", "user_64", "user_67"],
    ["user_18", "user_76", "user_37", "user_35", "user_93", "user_85", "user_78", "user_68", "user_63", "user_83", "user_44", "user_24", "user_6", "user_15", "user_9", "user_89", "user_42", "user_72", "user_28", "user_3", "user_39", "user_1", "user_94"],
    ["user_43", "user_54", "user_60", "user_51", "user_14", "user_21"],
    ["user_9", "user_24", "user_1", "user_18", "user_28", "user_35", "user_41", "user_37", "user_84", "user_3", "user_57", "user_70", "user_74", "user_49", "user_15", "user_53", "user_6", "user_10"],
    ["user_21", "user_71", "user_88", "user_65", "user_99", "user_14", "user_52"],
]

# Чтение точек из user_points.csv
points = {}
with open("user_points.csv", "r") as f:
    reader = csv.DictReader(f)
    for row in reader:
        user_id = row["user_id"]
        x = float(row["x"])
        y = float(row["y"])
        points[user_id] = (x, y)

# Визуализация кластеров
plt.figure(figsize=(8, 6))
colors = plt.get_cmap("tab10")

for idx, cluster in enumerate(clusters):
    xs = []
    ys = []
    for user_id in cluster:
        if user_id in points:
            x, y = points[user_id]
            xs.append(x)
            ys.append(y)
    plt.scatter(xs, ys, label=f"Cluster {idx}", alpha=0.8, s=60, edgecolors="k", color=colors(idx % 10))

plt.xlabel("RMSE typeA (0–9)")
plt.ylabel("RMSE typeB (10–19)")
plt.title("Кластеризация пользователей")
plt.grid(True)
plt.legend()
plt.tight_layout()
plt.savefig("clustered_users.png", dpi=300)
plt.close()

print("✅ Кластеры визуализированы в clustered_users.png")
