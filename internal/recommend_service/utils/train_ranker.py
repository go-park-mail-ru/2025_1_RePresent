import numpy as np
from catboost import CatBoostRanker, Pool
from sklearn.model_selection import train_test_split
from sklearn.metrics import ndcg_score
from loguru import logger


def calculate_ndcg(ranker, X, y, group_id):
    unique_groups = np.unique(group_id)
    scores = []
    for g in unique_groups:
        mask = group_id == g
        x = X[mask]
        y_true = y[mask]
        if len(y_true) < 2:
            continue
        y_pred = ranker.predict(x)
        scores.append(ndcg_score([y_true], [y_pred]))
    return round(np.mean(scores), 4)


# Примеры данных
X = np.array(
    [
        [0.95, 15, 12, 60, 0.3 * 0.95 + 15 * 0.7],
        [0.90, 12, 10, 50, 0.8 * 0.90 + 12 * 0.7],
        [0.88, 18, 12, 40, 0.6 * 0.88 + 18 * 0.7],
        [0.75, 10, 12, 40, 0.5 * 0.75 + 10 * 0.7],
        [0.70, 8, 8, 30, 0.4 * 0.70 + 8 * 0.7],
        [0.65, 14, 10, 50, 0.5 * 0.65 + 14 * 0.7],
        [0.4, 20, 15, 100, 0.2 * 0.4 + 20 * 0.7],
        [0.3, 18, 5, 20, 0.1 * 0.3 + 18 * 0.7],
        [0.25, 15, 4, 8, 0.1 * 0.25 + 15 * 0.7],
        [0.92, 20, 12, 60, 0.7 * 0.92 + 20 * 0.7],
    ]
)

y = np.array([5, 5, 4, 4, 3, 3, 2, 2, 1, 1])
group_id = np.array([0] * len(X))

X_train, X_eval, y_train, y_eval, g_train, g_eval = train_test_split(
    X, y, group_id, test_size=0.2, random_state=42
)

train_pool = Pool(data=X_train, label=y_train, group_id=g_train)
eval_pool = Pool(data=X_eval, label=y_eval, group_id=g_eval)

ranker = CatBoostRanker(
    iterations=300,
    learning_rate=0.05,
    loss_function="YetiRank",
    verbose=20,
    random_seed=42,
    l2_leaf_reg=3.0,
    depth=6,
    bagging_temperature=0.8,
    use_best_model=True,
    early_stopping_rounds=20,
)

ranker.fit(train_pool, eval_set=eval_pool)

metrics = calculate_ndcg(ranker, X_eval, y_eval, g_eval)
logger.info(f"NDCG after training: {metrics}")

ranker.save_model("reluma.cbm")
print("Model saved as reluma.cbm")
