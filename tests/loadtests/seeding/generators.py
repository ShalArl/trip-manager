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

"""
type CreateUserRequest struct {
	Email openapi_types.Email `json:"email"`
	Name  string              `json:"name"`

	// Password Minimum 8 characters
	Password string `json:"password"`
}
"""
def generate_user(i: int) -> dict:
    return {
        "email": f"user_{i}@loadtest.com",
        "password": "SuperStrongPW1",
        "name": fake.name(),
    }

"""
type CreateTripRequest struct {
    Description      *string            `json:"description,omitempty"`
    EndDate          openapi_types.Date `json:"endDate"`
    ShortDescription string             `json:"shortDescription"`
    StartDate        openapi_types.Date `json:"startDate"`
    Title            string             `json:"title"`
}
"""
def generate_trip() -> dict:
    start = fake.date_between(start_date="today", end_date="+1y")
    end = fake.date_between(start_date=start, end_date="+2y")

    # for better titles
    title_prefix = ["Welcome to",
                    "My trip in",
                    "Experiencing",
                    "Holidays in",
                    "With friends in",
                    "Checking out",
                    "Wonders of"]

    title = random.choice(title_prefix) + " " + fake.country()
    return {
        "title": title,
        "shortDescription": fake.text(max_nb_chars=80),
        "description": fake.text(max_nb_chars=400),
        "startDate": start.isoformat(),
        "endDate": end.isoformat(),
    }


"""
type CreateLocationRequest struct {
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Latitude  *float32 `json:"latitude,omitempty"`
	Longitude *float32 `json:"longitude,omitempty"`
	Name      string   `json:"name"`
	Notes     *string  `json:"notes,omitempty"`

	// Sequence Visit sequence number
	Sequence *int `json:"sequence,omitempty"`
}
"""

def generate_location() -> dict:
    start_date = fake.date_between(start_date="today", end_date="+1y")
    end_date = fake.date_between(
        start_date=start_date + timedelta(days=1),
        end_date=start_date + timedelta(days=10),
    )
    return {
        "title": fake.sentence(nb_words=3),  # ← neu
        "city": fake.city(),
        "country": fake.country(),
        "shortDescription": fake.sentence(nb_words=6),
        "dateFrom": start_date.isoformat(),
        "dateTo": end_date.isoformat(),
        "latitude": float(fake.latitude()),
        "longitude": float(fake.longitude()),
        "name": fake.sentence(nb_words=2),
        "notes": fake.sentence(nb_words=10),
        "sequence": random.randint(1, 10),
    }

"""
// CreateActivityRequest defines model for CreateActivityRequest.
type CreateActivityRequest struct {
	Category    *CreateActivityRequestCategory `json:"category,omitempty"`
	Cost        *float32                       `json:"cost,omitempty"`
	Currency    *string                        `json:"currency,omitempty"`
	Date        openapi_types.Date             `json:"date"`
	Description *string                        `json:"description,omitempty"`
	EndTime     *string                        `json:"endTime,omitempty"`
	LocationId  openapi_types.UUID             `json:"locationId"`
	Name        string                         `json:"name"`
	StartTime   *string                        `json:"startTime,omitempty"`
}
"""

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


"""
type PresignedURLRequest struct {
	// FileName Name of the file to upload (e.g., "avatar.jpg")
	FileName string `json:"fileName"`

	// MediaType Type of media being uploaded
	MediaType PresignedURLRequestMediaType `json:"mediaType"`
}
"""

def generate_image_upload_request() -> dict:
    return {
        "fileName": f"{fake.word()}.jpg",
        "mediaType": "image/jpeg",
    }


# Create some images for file upload testing
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