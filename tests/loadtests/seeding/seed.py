import asyncio
import csv
import httpx
from config.settings import BASE_URL, SEED_COUNT, SEED_CONCURRENCY, USERS_CSV

from seeding.generators import generate_user


async def register_user(client: httpx.AsyncClient, semaphore: asyncio.Semaphore, i: int) -> dict | None:

    user = generate_user(i)

    async with semaphore:
        reg = await client.post("/auth/register", json={
            "email": user["email"],
            "password": user["password"],
            "name": user["name"],
        })

        if reg.status_code != 201:
            print(f"[WARN] Register failed for {user['email']}: {reg.status_code}")
            return None

        token = reg.json()["token"]

        return {
            "email": user["email"],
            "password": user["password"],
            "name": user["name"],
            "token": token
        }


async def seed():
    semaphore = asyncio.Semaphore(SEED_CONCURRENCY)

    async with httpx.AsyncClient(base_url=BASE_URL, timeout=30) as client:
        tasks = [
            register_user(client, semaphore, i)
            for i in range(SEED_COUNT)
        ]

        results = await asyncio.gather(*tasks)

    successful = [r for r in results if r is not None]

    if len(successful) == 0:
        print(f"[ERROR] Failed to register users! Is the backend up and running?")

    print(f"[INFO] Registered {len(successful)} users")

    with open(USERS_CSV, "w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=["email", "password", "name", "token"])
        writer.writeheader()
        writer.writerows(successful)

    print(f"Seeded {len(successful)}/{SEED_COUNT} users → {USERS_CSV}")


def main():
    asyncio.run(seed())


if __name__ == "__main__":
    main()
