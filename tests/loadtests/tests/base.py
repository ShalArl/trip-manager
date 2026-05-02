import csv
import random
import time
import requests

from locust import HttpUser, between
from config.settings import FIREBASE_API_KEY, USERS_CSV

from functools import wraps

FIREBASE_REFRESH_URL = (
    f"https://securetoken.googleapis.com/v1/token?key={FIREBASE_API_KEY}"
)


def with_auth(func):
    """Decorator that ensures the user has a valid token before executing the task."""

    @wraps(func)
    def wrapper(self, *args, **kwargs):
        self._ensure_token()
        return func(self, *args, **kwargs)

    return wrapper


# Buffer in Sekunden bevor das Token abläuft
TOKEN_REFRESH_BUFFER = 120  # 2 Minuten

with open(USERS_CSV) as f:
    USERS = list(csv.DictReader(f))


def refresh_id_token(refresh_token: str) -> tuple[str, int]:
    """Holt neuen ID-Token via Refresh-Token.

    Returns:
        (id_token, expires_at_unix_timestamp)
    """
    resp = requests.post(
        FIREBASE_REFRESH_URL,
        data={
            "grant_type": "refresh_token",
            "refresh_token": refresh_token,
        },
        timeout=10,
    )
    resp.raise_for_status()
    data = resp.json()
    expires_in = int(data["expires_in"])
    return data["id_token"], int(time.time()) + expires_in - TOKEN_REFRESH_BUFFER


class BaseUser(HttpUser):
    abstract = True
    wait_time = between(1, 3)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        user = random.choice(USERS)
        self.email: str = user["email"]
        self.refresh_token: str = user["refresh_token"]
        self.firebase_uid: str = user["firebase_uid"]
        self.id_token: str | None = None
        self.token_expires_at: int = 0

    def _ensure_token(self):
        """Refresh ID-Token wenn er abläuft. Vor jedem Backend-Call aufrufen."""
        if time.time() >= self.token_expires_at:
            try:
                self.id_token, self.token_expires_at = refresh_id_token(self.refresh_token)
                self.client.headers["Authorization"] = f"Bearer {self.id_token}"
            except requests.HTTPError as e:
                print(f"[WARN] Token refresh failed for {self.email}: {e}")

    def on_start(self):
        self._ensure_token()
