import csv
import random
from locust import HttpUser, between
from config.settings import USERS_CSV

with open(USERS_CSV) as f:
    USERS = list(csv.DictReader(f))


class BaseUser(HttpUser):
    abstract = True
    wait_time = between(1, 3)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.email: str = ""

    def on_start(self):
        user = random.choice(USERS)
        self.client.headers["Authorization"] = f"Bearer {user['token']}"
        self.email = user["email"]