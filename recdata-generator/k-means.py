import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from sklearn.cluster import KMeans

# --- параметры ---
input_csv = "user_vectors.csv"
output_csv = "user_clusters.csv"
output_png = "user_clusters_projected.png"
n_clusters = 7  # количество кластеров

# --- загрузка ---
df = pd.read_csv(input_csv)

# Задаём user_id явно — как user_0, user_1, ...
df["user_id"] = [f"user_{i}" for i in range(len(df))]

# Разделим: вектора, проекция, метки
vector_columns = [f"v{i}" for i in range(20)]
x_column = "x_rmse"
y_column = "y_rmse"

X = df[vector_columns].values

# --- кластеризация ---
kmeans = KMeans(n_clusters=n_clusters, random_state=42, n_init="auto")
df["cluster"] = kmeans.fit_predict(X)

# --- сохраняем кластеры в CSV ---
df[["user_id", "cluster"]].to_csv(output_csv, index=False)

# --- визуализация ---
plt.figure(figsize=(8, 6))
colors = plt.get_cmap("tab10")

for cluster_id in range(n_clusters):
    cluster_df = df[df["cluster"] == cluster_id]
    plt.scatter(
        cluster_df[x_column],
        cluster_df[y_column],
        label=f"Cluster {cluster_id}",
        s=60,
        alpha=0.75,
        edgecolors="k",
        color=colors(cluster_id)
    )

plt.xlabel("RMSE typeA (0–9)")
plt.ylabel("RMSE typeB (10–19)")
plt.title(f"K-Means Clustering (k={n_clusters}) on Projected 2D Space")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig(output_png, dpi=300)
plt.close()

print(f"✅ Всё готово: {output_csv}, {output_png}")
