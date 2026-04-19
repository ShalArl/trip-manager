import os

BASE_URL = os.getenv("BASE_URL", "http://localhost:8000/api")
SEED_COUNT = int(os.getenv("SEED_COUNT", 5000))
SEED_CONCURRENCY = int(os.getenv("SEED_CONCURRENCY", 50))
USERS_CSV = "data/users.csv"