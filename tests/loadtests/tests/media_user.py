import random
import httpx

from tests.base import BaseUser, with_auth
from locust import between, task
from seeding.generators import generate_fake_images

IMAGES = generate_fake_images(count=10)


class MediaUser(BaseUser):
    weight = 10
    wait_time = between(3, 8)

    def on_start(self):
        super().on_start()

    @task
    @with_auth
    def upload_image(self):
        resp = self.client.post("/uploads/presigned", json={
            "fileName": "image",
            "mediaType": "avatar"
        })

        if resp.status_code != 200:
            return

        data = resp.json()
        upload_url = data["presignedUrl"]

        image_bytes = random.choice(IMAGES)
        # mit Absicht via httpx, damit die performance von storage nicht die metrik des backends beeinflusst
        with httpx.Client() as client:
            client.put(upload_url, content=image_bytes, headers={
                "Content-Type": "image/jpg",
            })