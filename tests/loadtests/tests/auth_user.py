import csv
import random
from locust import task, between
from tests.base import BaseUser, USERS


class AuthUser(BaseUser):
    weight = 5
    wait_time = between(5, 10)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        user = random.choice(USERS)
        self.email = user["email"]
        self.password = user["password"]

    def on_start(self):
        self._login()

    def _login(self):
        resp = self.client.post("/auth/login", json={
            "email": self.email,
            "password": self.password
        })
        if resp.status_code == 200:
            self.client.headers["Authorization"] = f"Bearer {resp.json()['token']}"

    @task
    def login_logout(self):
        self.client.headers.pop("Authorization", None)
        self._login()