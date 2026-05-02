import os
from dotenv import load_dotenv


load_dotenv()

BASE_URL = os.getenv("BASE_URL", "http://localhost:8000/api")
SEED_COUNT = int(os.getenv("SEED_COUNT", 5000))
SEED_CONCURRENCY = int(os.getenv("SEED_CONCURRENCY", 1))
FIREBASE_API_KEY = os.getenv("FIREBASE_API_KEY", "your-firebase-api-key")
FIREBASE_PROJECT_ID = os.getenv("FIREBASE_PROJECT_ID", "your-firebase-project-id")
USERS_CSV = "data/users.csv"
