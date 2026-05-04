import random
from tests.base import BaseUser, USERS, with_auth
from locust import between, task

from seeding.generators import get_bio, get_name, get_password

class AccountUser(BaseUser):
    weight = 3
    wait_time = between(1, 3)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        user = random.choice(USERS)
        self.current_password = user["password"]

    def on_start(self):
        super().on_start()

    @task
    def get_profile(self):
        self.client.get("/users/me")

    @task
    def update_profile(self):
        self.client.put("/users/me", json={
            "name": get_name(),
            "bio": get_bio(),
        })

    @task
    def change_password(self):
        self.client.put("/users/me/password", json={
            "currentPassword": self.current_password,
            "newPassword": get_password()
        })