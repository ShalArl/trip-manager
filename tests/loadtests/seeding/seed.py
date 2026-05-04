import asyncio
import csv
import concurrent.futures
from itertools import islice

import firebase_admin
import httpx
from firebase_admin import auth, credentials

from config.settings import (
    BASE_URL,
    FIREBASE_API_KEY,
    SEED_COUNT,
    SEED_CONCURRENCY,
    USERS_CSV,
    FIREBASE_PROJECT_ID,
)


def init_firebase():
    if not firebase_admin._apps:
        firebase_admin.initialize_app(
            credentials.ApplicationDefault(),
            {
                "projectId": FIREBASE_PROJECT_ID,
                "serviceAccountId": f"trip-manager-runtime@{FIREBASE_PROJECT_ID}.iam.gserviceaccount.com",
            },
        )


def chunks(iterable, n):
    """Yield successive n-sized chunks."""
    it = iter(iterable)
    while batch := list(islice(it, n)):
        yield batch


def bulk_create_firebase_users(count: int) -> list[dict]:
    """
    Erstellt count User via Firebase Admin SDK Bulk-Import.
    Nutzt parallele Threads weil das SDK synchron ist.
    """
    print(f"[INFO] Creating {count} Firebase users via Admin SDK...")

    # ImportUserRecord ohne Password — Login geht über Custom Tokens
    user_records = [
        auth.ImportUserRecord(
            uid=f"loadtest-{i}",
            email=f"user_{i}@loadtest.com",
            display_name=f"Loadtest User {i}",
        )
        for i in range(count)
    ]

    created_users = []

    # Batches von 1000 (Firebase Limit) parallel verarbeiten
    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        futures = []
        for batch in chunks(user_records, 1000):
            futures.append(executor.submit(_import_batch, batch))

        for future in concurrent.futures.as_completed(futures):
            result = future.result()
            created_users.extend(result)

    print(f"[INFO] Created {len(created_users)} Firebase users")
    return created_users


def _import_batch(batch: list[auth.ImportUserRecord]) -> list[dict]:
    """Wird in Thread ausgeführt."""
    result = auth.import_users(batch)

    if result.failure_count > 0:
        for error in result.errors[:5]:
            print(f"[WARN] Import error: {error.reason}")

    # Erfolgreiche User zurückgeben (Index entspricht Reihenfolge in batch)
    successful = []
    failed_indices = {e.index for e in result.errors}

    for i, record in enumerate(batch):
        if i not in failed_indices:
            successful.append({
                "uid": record.uid,
                "email": record.email,
                "name": record.display_name,
            })

    return successful


async def get_id_token_for_uid(client: httpx.AsyncClient, uid: str) -> dict | None:
    try:
        custom_token = auth.create_custom_token(uid)
        if isinstance(custom_token, bytes):
            custom_token = custom_token.decode()
    except Exception as e:
        print(f"[ERROR] create_custom_token failed for {uid}: {e}")
        return None

    try:
        resp = await client.post(
            f"https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key={FIREBASE_API_KEY}",
            json={"token": custom_token, "returnSecureToken": True},
            timeout=15,
        )
    except Exception as e:
        print(f"[ERROR] Custom token exchange request failed for {uid}: {e}")
        return None

    if resp.status_code != 200:
        print(f"[WARN] Custom token exchange failed for {uid}: {resp.status_code} {resp.text[:300]}")
        return None

    data = resp.json()
    return {
        "id_token": data["idToken"],
        "refresh_token": data["refreshToken"],
    }


async def provision_backend_user(client: httpx.AsyncClient, id_token: str, name: str) -> bool:
    try:
        resp = await client.post(
            f"{BASE_URL}/api/users/provision",
            headers={"Authorization": f"Bearer {id_token}"},
            json={"name": name},
            timeout=20,
        )
    except Exception as e:
        print(f"[ERROR] Provision request failed: {e}")
        return False

    if resp.status_code in (200, 201):
        return True

    print(f"[WARN] Provision failed: {resp.status_code} {resp.text[:300]}")
    return False


async def harvest_token_and_provision(
        client: httpx.AsyncClient,
        semaphore: asyncio.Semaphore,
        user: dict,
) -> dict | None:
    async with semaphore:
        token_data = await get_id_token_for_uid(client, user["uid"])
        if not token_data:
            return None

        success = await provision_backend_user(client, token_data["id_token"], user["name"])
        if not success:
            return None

        return {
            "email": user["email"],
            "name": user["name"],
            "firebase_uid": user["uid"],
            "refresh_token": token_data["refresh_token"],
        }


async def main():
    init_firebase()

    # Phase 1: Bulk-Create in Firebase (schnell)
    firebase_users = bulk_create_firebase_users(SEED_COUNT)

    if not firebase_users:
        print("[ERROR] No users created in Firebase")
        return

    # Phase 2: Token harvesting + Backend provisioning (parallel)
    print(f"[INFO] Harvesting tokens and provisioning {len(firebase_users)} users...")

    semaphore = asyncio.Semaphore(SEED_CONCURRENCY)

    async with httpx.AsyncClient(timeout=40) as client:
        tasks = [
            harvest_token_and_provision(client, semaphore, user)
            for user in firebase_users
        ]
        results = await asyncio.gather(*tasks, return_exceptions=True)

    successful = [r for r in results if isinstance(r, dict)]

    if successful:
        with open(USERS_CSV, "w", newline="") as f:
            writer = csv.DictWriter(
                f, fieldnames=["email", "name", "firebase_uid", "refresh_token"]
            )
            writer.writeheader()
            writer.writerows(successful)
        print(f"[DONE] Seeded {len(successful)}/{SEED_COUNT} users → {USERS_CSV}")
    else:
        print("[ERROR] No users seeded successfully")


if __name__ == "__main__":
    asyncio.run(main())