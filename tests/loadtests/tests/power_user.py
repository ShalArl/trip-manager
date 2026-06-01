import random

from locust import task, between
from seeding.generators import generate_trip, generate_location, generate_comment
from tests.base import BaseUser, with_auth


class PowerUser(BaseUser):
    weight = 20
    wait_time = between(1, 3)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.trip_ids: list[str] = []
        self.location_ids: dict[str, list[str]] = {}
        self.activity_ids: dict[str, dict[str, list[str]]] = {}
        self.comment_ids: dict[str, list[str]] = {}

    def on_start(self):
        super().on_start()

    # -- Helpers --

    def _random_trip_id(self) -> str | None:
        return random.choice(self.trip_ids) if self.trip_ids else None

    def _random_location_id(self, trip_id: str) -> str | None:
        ids = self.location_ids.get(trip_id, [])
        return random.choice(ids) if ids else None

    # -- Tasks --

    @task(3)
    @with_auth
    def create_full_trip(self):
        resp = self.client.post("/trips", json=generate_trip())
        if resp.status_code != 201:
            print(f"[DEBUG] create trip error {resp.status_code}: {resp.text[:500]}")
            return
        trip_id = resp.json()["id"]
        self.trip_ids.append(trip_id)
        self.location_ids[trip_id] = []
        self.activity_ids[trip_id] = {}

        loc_resp = self.client.post(f"/trips/{trip_id}/locations", json=generate_location())
        if loc_resp.status_code != 201:
            print(f"[DEBUG] create location error {loc_resp.status_code}: {loc_resp.text[:500]}")
            return
        location_id = loc_resp.json()["id"]
        self.location_ids[trip_id].append(location_id)
        self.activity_ids[trip_id][location_id] = []

        comment_resp = self.client.post(f"/social/{trip_id}/comments", json=generate_comment())
        if comment_resp.status_code == 201:
            self.comment_ids.setdefault(trip_id, []).append(comment_resp.json()["id"])
        else:
            print(f"[DEBUG] create comment error {comment_resp.status_code}: {comment_resp.text[:500]}")

    @task(2)
    @with_auth
    def add_location(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        resp = self.client.post(f"/trips/{trip_id}/locations", json=generate_location())
        if resp.status_code != 201:
            print(f"[DEBUG] add_location {resp.status_code}: {resp.text[:500]}")
            return
        location_id = resp.json()["id"]
        self.location_ids[trip_id].append(location_id)
        self.activity_ids[trip_id][location_id] = []

    @task(3)
    @with_auth
    def write_comment(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        resp = self.client.post(f"/social/{trip_id}/comments", json=generate_comment())
        if resp.status_code in (200, 201):
            self.comment_ids.setdefault(trip_id, []).append(resp.json()["id"])
        else:
            print(f"[DEBUG] write_comment error {resp.status_code}: {resp.text[:500]}")

    @task(1)
    @with_auth
    def write_and_delete_comment(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        resp = self.client.post(f"/social/{trip_id}/comments", json=generate_comment())
        if resp.status_code not in (200, 201):
            print(f"[DEBUG] write_and_delete_comment error {resp.status_code}: {resp.text[:500]}")
            return
        comment_id = resp.json()["id"]
        self.comment_ids.setdefault(trip_id, []).append(comment_id)

        comment_ids = self.comment_ids.get(trip_id, [])
        if not comment_ids:
            return
        comment_id = random.choice(comment_ids)
        del_resp = self.client.delete(f"/social/{trip_id}/comments/{comment_id}")
        if del_resp.status_code == 204:
            self.comment_ids[trip_id].remove(comment_id)
        else:
            print(f"[DEBUG] delete_comment error {del_resp.status_code}: {del_resp.text[:500]}")

    @task(2)
    @with_auth
    def toggle_like(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        with self.client.post(f"/social/{trip_id}/likes", catch_response=True) as resp:
            if resp.status_code in (200, 201):
                resp.success()
            elif resp.status_code == 409:
                # Already liked → unlike (Toggle wie im Frontend)
                resp.success()
                self.client.delete(f"/social/{trip_id}/likes")
            else:
                resp.failure(f"toggle_like error {resp.status_code}: {resp.text[:200]}")

    @task(1)
    @with_auth
    def update_trip(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        resp = self.client.put(f"/trips/{trip_id}", json=generate_trip())
        if resp.status_code not in (200, 204):
            print(f"[DEBUG] update_trip error {resp.status_code}: {resp.text[:500]}")

    @task(1)
    @with_auth
    def update_location(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return
        with self.client.put(f"/locations/{trip_id}/{location_id}", json=generate_location(), catch_response=True) as resp:
            if resp.status_code in (200, 204):
                resp.success()
            elif resp.status_code == 404:
                resp.success()  # Location schon gelöscht - kein Fehler
                if location_id in self.location_ids.get(trip_id, []):
                    self.location_ids[trip_id].remove(location_id)
                    self.activity_ids[trip_id].pop(location_id, None)
            else:
                resp.failure(f"update_location error {resp.status_code}: {resp.text[:200]}")

    @task(1)
    @with_auth
    def delete_trip(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        resp = self.client.delete(f"/trips/{trip_id}")
        if resp.status_code in (204, 404):
            if trip_id in self.trip_ids:
                self.trip_ids.remove(trip_id)
            self.location_ids.pop(trip_id, None)
            self.activity_ids.pop(trip_id, None)
        else:
            print(f"[DEBUG] delete_trip error {resp.status_code}: {resp.text[:500]}")

    @task(1)
    @with_auth
    def delete_location(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return
        with self.client.delete(f"/locations/{trip_id}/{location_id}", catch_response=True) as resp:
            if resp.status_code in (204, 404):
                resp.success()  # 404 = schon gelöscht, kein Fehler
                if location_id in self.location_ids.get(trip_id, []):
                    self.location_ids[trip_id].remove(location_id)
                    self.activity_ids[trip_id].pop(location_id, None)
            else:
                resp.failure(f"delete_location error {resp.status_code}: {resp.text[:200]}")