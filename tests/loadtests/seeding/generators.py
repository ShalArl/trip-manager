import io

import random
from datetime import timedelta

from faker import Faker
from PIL import Image

fake = Faker()

####################################################################

# General stuff

def get_bio() -> str:
    return fake.text(max_nb_chars=400)


def get_name() -> str:
    return fake.name()


def get_password() -> str:
    return "SuperStrongPW1"

####################################################################

def generate_user(i: int) -> dict:
    return {
        "email": f"user_{i}@loadtest.com",
        "password": "SuperStrongPW1",
        "name": fake.name(),
    }

def generate_trip() -> dict:
    start = fake.date_between(start_date="today", end_date="+1y")
    end = fake.date_between(start_date=start, end_date="+2y")

    title_prefix = ["Welcome to", "My trip in", "Experiencing", "Holidays in",
                    "With friends in", "Checking out", "Wonders of"]

    title = random.choice(title_prefix) + " " + fake.country()
    return {
        "title": title,
        "shortDescription": fake.text(max_nb_chars=80),
        "description": fake.text(max_nb_chars=400),
        "startDate": start.isoformat(),
        "endDate": end.isoformat(),
    }


EXAMPLE_LOCATIONS = [
    {"title": "Eiffel Tower",      "name": "Eiffel Tower",      "city": "Paris",      "country": "France",         "countryCode": "FR", "shortDescription": "Iconic iron lattice tower on the Champ de Mars.", "latitude": 48.8584,  "longitude": 2.2945},
    {"title": "Colosseum",         "name": "Colosseum",         "city": "Rome",        "country": "Italy",          "countryCode": "IT", "shortDescription": "Ancient amphitheatre in the centre of Rome.",     "latitude": 41.8902,  "longitude": 12.4922},
    {"title": "Brandenburg Gate",  "name": "Brandenburg Gate",  "city": "Berlin",      "country": "Germany",        "countryCode": "DE", "shortDescription": "18th-century neoclassical monument in Berlin.",   "latitude": 52.5163,  "longitude": 13.3777},
    {"title": "Sagrada Familia",   "name": "Sagrada Familia",   "city": "Barcelona",   "country": "Spain",          "countryCode": "ES", "shortDescription": "Large unfinished Roman Catholic basilica.",       "latitude": 41.4036,  "longitude": 2.1744},
    {"title": "Acropolis",         "name": "Acropolis",         "city": "Athens",      "country": "Greece",         "countryCode": "GR", "shortDescription": "Ancient citadel located on a rocky outcrop.",     "latitude": 37.9715,  "longitude": 23.7267},
    {"title": "Anne Frank House",  "name": "Anne Frank House",  "city": "Amsterdam",   "country": "Netherlands",    "countryCode": "NL", "shortDescription": "Museum dedicated to wartime diarist Anne Frank.", "latitude": 52.3752,  "longitude": 4.8840},
    {"title": "Charles Bridge",    "name": "Charles Bridge",    "city": "Prague",      "country": "Czech Republic", "countryCode": "CZ", "shortDescription": "Historic bridge crossing the Vltava river.",      "latitude": 50.0865,  "longitude": 14.4114},
    {"title": "Wawel Castle",      "name": "Wawel Castle",      "city": "Krakow",      "country": "Poland",         "countryCode": "PL", "shortDescription": "A fortified architectural complex in Krakow.",    "latitude": 50.0540,  "longitude": 19.9355},
    {"title": "Schönbrunn Palace", "name": "Schönbrunn Palace", "city": "Vienna",      "country": "Austria",        "countryCode": "AT", "shortDescription": "Former imperial summer residence in Vienna.",     "latitude": 48.1845,  "longitude": 16.3122},
    {"title": "Old Town Square",   "name": "Old Town Square",   "city": "Tallinn",     "country": "Estonia",        "countryCode": "EE", "shortDescription": "Well-preserved medieval old town in Tallinn.",   "latitude": 59.4370,  "longitude": 24.7536},
]

def generate_location() -> dict:
    base = random.choice(EXAMPLE_LOCATIONS).copy()
    start_date = fake.date_between(start_date="today", end_date="+1y")
    end_date = fake.date_between(
        start_date=start_date + timedelta(days=1),
        end_date=start_date + timedelta(days=10),
    )
    base["dateFrom"] = start_date.isoformat()
    base["dateTo"] = end_date.isoformat()
    base["notes"] = fake.sentence(nb_words=10)
    base["sequence"] = random.randint(1, 10)
    return base


def generate_activity(location_id: str) -> dict:
    return {
        "name": fake.sentence(nb_words=3),
        "date": fake.date_between(start_date="today", end_date="+1y").isoformat(),
        "locationId": location_id,
        "description": fake.text(max_nb_chars=200),
        "cost": round(random.uniform(10, 500), 2),
        "currency": random.choice(["USD", "EUR", "GBP", "JPY"]),
        "category": random.choice(["Sightseeing", "Adventure", "Relaxation", "Cultural"]),
    }


def generate_image_upload_request() -> dict:
    return {
        "fileName": f"{fake.word()}.jpg",
        "mediaType": "image/jpeg",
    }


def generate_fake_images(count=10, width=800, height=600) -> list[bytes]:
    images = []
    for _ in range(count):
        img = Image.new("RGB", (width, height), color=(
            random.randint(0, 255),
            random.randint(0, 255),
            random.randint(0, 255),
        ))
        buf = io.BytesIO()
        img.save(buf, format="JPEG")
        images.append(buf.getvalue())
    return images


def generate_comment() -> dict:
    return {
        "text": fake.sentence(nb_words=20),
    }


def generate_city() -> str:
    return fake.city()