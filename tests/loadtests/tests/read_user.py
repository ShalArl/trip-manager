import random

from locust import task, between
from tests.base import BaseUser, with_auth

from seeding.generators import generate_city


class ReadUser(BaseUser):
    weight = 75
    wait_time = between(1, 3)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.trip_ids: list[str] = []
        self.location_ids: dict[str, list[str]] = {}

    def on_start(self):
        super().on_start()
        # Lade existierende Trips beim Start
        resp = self.client.get("/trips")
        if resp.status_code == 200:
            trips = resp.json()
            if isinstance(trips, list):
                self.trip_ids = [t["id"] for t in trips]
            elif isinstance(trips, dict) and "data" in trips:
                self.trip_ids = [t["id"] for t in trips["data"]]

        for trip_id in self.trip_ids[:5]:  # nur erste 5 um on_start nicht zu verlangsamen
            loc_resp = self.client.get(f"/trips/{trip_id}/locations")
            if loc_resp.status_code == 200:
                locs = loc_resp.json()
                if isinstance(locs, list):
                    self.location_ids[trip_id] = [l["id"] for l in locs]
                elif isinstance(locs, dict) and "data" in locs:
                    self.location_ids[trip_id] = [l["id"] for l in locs["data"]]

    @task(3)
    def get_trips(self):
        self.client.get("/trips")

    @task(2)
    def get_recent_trips(self):
        self.client.get("/trips/recent")

    @task(2)
    def get_trip_details(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/trips/{trip_id}")

    @task(2)
    def get_locations(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/trips/{trip_id}/locations")

    @task(1)
    def read_comments(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/social/{trip_id}/comments")

    @task(1)
    def read_likes(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/social/{trip_id}/likes")

    @task(2)
    def search_trips(self):
        city = generate_city()
        limit = random.randint(10, 50)
        self.client.get(f"/trips/search?q={city}&limit={limit}&offset=0")

    @task(1)
    def get_feed(self):
        self.client.get("/feed")

    @task(1)
    @with_auth
    def get_personal_feed(self):
        self.client.get("/feed/personal")

    @task(1)
    @with_auth
    def get_newsletter(self):
        self.client.get("/newsletter")