import random

from locust import task, between
from tests.base import BaseUser

from seeding.generators import generate_city


class ReadUser(BaseUser):
    weight = 0
    wait_time = between(1, 3)


    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.trip_ids: list[str] = []
        self.location_ids: dict[str, list[str]] = {}
        self.activity_ids: dict[str, dict[str, list[str]]] = {}

    def on_start(self):
        super().on_start()

        for trip_id in self.trip_ids:
            loc_resp = self.client.get(f"/trips/{trip_id}/locations")
            if loc_resp.status_code != 200:
                continue

            location_ids = [l["id"] for l in loc_resp.json()]
            self.location_ids[trip_id] = location_ids

            self.activity_ids[trip_id] = {}
            for location_id in location_ids:
                act_resp = self.client.get(f"/trips/{trip_id}/locations/{location_id}/activities")
                if act_resp.status_code == 200:
                    self.activity_ids[trip_id][location_id] = [a["id"] for a in act_resp.json()]

    @task
    def get_trips(self):
        self.client.get("/trips")

    @task
    def get_trip_details(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/trips/{trip_id}")

    @task
    def get_locations(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/trips/{trip_id}/locations")

    @task
    def get_activities(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        location_ids = self.location_ids.get(trip_id, [])
        if not location_ids:
            return
        location_id = random.choice(location_ids)
        self.client.get(f"/trips/{trip_id}/locations/{location_id}/activities")

    @task
    def read_comments(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/trips/{trip_id}/comments")

    @task
    def read_likes(self):
        if not self.trip_ids:
            return
        trip_id = random.choice(self.trip_ids)
        self.client.get(f"/trips/{trip_id}/likes")

    @task
    def search_trips(self):
        city = generate_city()
        limit = random.randint(10, 50)
        self.client.get(f"/trips/search?q={city}&limit={limit}&offset=0")