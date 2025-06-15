import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from sklearn.cluster import KMeans
from sklearn.decomposition import PCA

# --- параметры ---
input_csv = "user_vectors.csv"
output_csv = "user_clusters_pca.csv"
output_png = "user_clusters_pca.png"
n_clusters = 7

# --- загрузка ---
df = pd.read_csv(input_csv)
df["user_id"] = [f"user_{i}" for i in range(len(df))]
vector_columns = [f"v{i}" for i in range(20)]
X = df[vector_columns].values

# --- PCA ---
pca = PCA(n_components=2)
X_2d = pca.fit_transform(X)

df["pca_x"] = X_2d[:, 0]
df["pca_y"] = X_2d[:, 1]

# --- кластеризация ---
kmeans = KMeans(n_clusters=n_clusters, random_state=42, n_init="auto")
df["cluster"] = kmeans.fit_predict(X)

# --- CSV ---
df[["user_id", "cluster", "pca_x", "pca_y"]].to_csv(output_csv, index=False)

# --- Визуализация ---
plt.figure(figsize=(8, 6))
colors = plt.get_cmap("tab10")

for cluster_id in range(n_clusters):
    cluster_df = df[df["cluster"] == cluster_id]
    plt.scatter(
        cluster_df["pca_x"],
        cluster_df["pca_y"],
        label=f"Cluster {cluster_id}",
        s=60,
        alpha=0.75,
        edgecolors="k",
        color=colors(cluster_id)
    )

plt.xlabel("PCA 1")
plt.ylabel("PCA 2")
plt.title(f"K-Means Clustering (k={n_clusters}) via PCA")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig(output_png, dpi=300)
plt.close()

print(f"✅ Готово: {output_csv}, {output_png}")
