import random

import httpx
from locust import task, between
from seeding.generators import generate_trip, generate_location, generate_activity, generate_comment
from tests.base import BaseUser, with_auth

from tests.media_user import IMAGES


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

    def _random_activity_id(self, trip_id: str, location_id: str) -> str | None:
        ids = self.activity_ids.get(trip_id, {}).get(location_id, [])
        return random.choice(ids) if ids else None

    # -- Tasks --

    @task(3)
    @with_auth
    def create_full_trip(self):
        # Trip
        resp = self.client.post("/trips", json=generate_trip())
        if resp.status_code != 201:
            return
        trip_id = resp.json()["id"]
        self.trip_ids.append(trip_id)
        self.location_ids[trip_id] = []
        self.activity_ids[trip_id] = {}

        # Location
        loc_resp = self.client.post(f"/trips/{trip_id}/locations", json=generate_location())
        if loc_resp.status_code != 201:
            return
        location_id = loc_resp.json()["id"]
        self.location_ids[trip_id].append(location_id)
        self.activity_ids[trip_id][location_id] = []

        # Activity
        act_resp = self.client.post(
            f"/trips/{trip_id}/locations/{location_id}/activities",
            json=generate_activity(location_id)
        )
        if act_resp.status_code == 201:
            activity_id = act_resp.json()["id"]
            self.activity_ids[trip_id][location_id].append(activity_id)

        # add comment to own trip
        comment_resp = self.client.post(f"/trips/{trip_id}/comments", json=generate_comment())
        if comment_resp == 201:
            comment_id = comment_resp.json()["id"]
            self.comment_ids[trip_id].append(comment_id)

    @task(2)
    @with_auth
    def add_location(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return

        payload = generate_location()
        print(f"[DEBUG] Sending: {payload}")  # ← was schicken wir?

        resp = self.client.post(f"/trips/{trip_id}/locations", json=payload)

        if resp.status_code != 201:
            print(f"[DEBUG] {resp.status_code}: {resp.text[:500]}")  # ← was sagt Backend?
            return

        location_id = resp.json()["id"]
        self.location_ids[trip_id].append(location_id)
        self.activity_ids[trip_id][location_id] = []


    @task(3)
    @with_auth
    def add_image_to_location(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return

        # 1. Presigned URL holen
        presigned_resp = self.client.post("/uploads/presigned", json={
            "fileName": "location-image.jpg",
            "mediaType": "location-image",
        })
        if presigned_resp.status_code != 200:
            return

        data = presigned_resp.json()
        upload_url = data["presignedUrl"]
        image_key = data["objectKey"]

        image_bytes = random.choice(IMAGES)
        with httpx.Client() as client:
            client.put(upload_url, content=image_bytes, headers={"Content-Type": "image/jpeg"})

        # 3. Bei Location registrieren
        self.client.post(
            f"/trips/{trip_id}/locations/{location_id}/images",
            json={"imageKey": image_key},
        )


    @task(2)
    @with_auth
    def add_activity(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return
        resp = self.client.post(
            f"/trips/{trip_id}/locations/{location_id}/activities",
            json=generate_activity(location_id)
        )
        if resp.status_code == 201:
            print(f"[DEBUG] {resp.status_code}: {resp.text[:500]}")
            self.activity_ids[trip_id][location_id].append(resp.json()["id"])

    @task(3)
    @with_auth
    def write_comment(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return

        resp = self.client.post(
            f"/trips/{trip_id}/comments",
            json={"text": "Test comment"},
        )

        if resp.status_code in (200, 201):
            print(f"[DEBUG] {resp.status_code}: {resp.text[:500]}")  # ← was sagt Backend?
            self.comment_ids.setdefault(trip_id, []).append(resp.json()["id"])

    @task(1)
    @with_auth
    def write_and_delete_comment(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return

        resp = self.client.post(
            f"/trips/{trip_id}/comments",
            json={"text": "Test comment"},
        )

        if resp not in(200, 201):
            print(f"[DEBUG] {resp.status_code}: {resp.text[:500]}")  # ← was sagt Backend?
            return

        comment_ids = self.comment_ids.get(trip_id, [])

        if not comment_ids:
            return

        comment_id = random.choice(comment_ids)
        resp = self.client.delete(f"/trips/{trip_id}/comments/{comment_id}")
        if resp.status_code == 204:
            self.comment_ids[trip_id].remove(comment_id)

    @task(1)
    @with_auth
    def update_trip(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        self.client.put(f"/trips/{trip_id}", json=generate_trip())

    @task(1)
    def update_location(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return
        self.client.put(f"/trips/{trip_id}/locations/{location_id}", json=generate_location())

    @task(1)
    @with_auth
    def update_activity(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return
        activity_id = self._random_activity_id(trip_id, location_id)
        if not activity_id:
            return
        self.client.put(
            f"/trips/{trip_id}/locations/{location_id}/activities/{activity_id}",
            json=generate_activity(location_id)
        )

    @task(1)
    @with_auth
    def delete_trip(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        resp = self.client.delete(f"/trips/{trip_id}")
        if resp.status_code == 204:
            self.trip_ids.remove(trip_id)
            self.location_ids.pop(trip_id, None)
            self.activity_ids.pop(trip_id, None)

    @task(1)
    @with_auth
    def delete_location(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return
        resp = self.client.delete(f"/trips/{trip_id}/locations/{location_id}")
        if resp.status_code == 204:
            self.location_ids[trip_id].remove(location_id)
            self.activity_ids[trip_id].pop(location_id, None)

    @task(1)
    @with_auth
    def delete_activity(self):
        trip_id = self._random_trip_id()
        if not trip_id:
            return
        location_id = self._random_location_id(trip_id)
        if not location_id:
            return
        activity_id = self._random_activity_id(trip_id, location_id)
        if not activity_id:
            return
        resp = self.client.delete(
            f"/trips/{trip_id}/locations/{location_id}/activities/{activity_id}"
        )
        if resp.status_code == 204:
            self.activity_ids[trip_id][location_id].remove(activity_id)