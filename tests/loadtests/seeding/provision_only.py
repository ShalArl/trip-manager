#!/usr/bin/env python3
import asyncio
import csv
import logging
import httpx
import requests
from config.settings import FIREBASE_API_KEY, USERS_CSV, BASE_URL, SEED_CONCURRENCY

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler('data/provision.log')
    ]
)
logger = logging.getLogger(__name__)


async def provision_all():
    try:
        with open(USERS_CSV) as f:
            users = list(csv.DictReader(f))
    except FileNotFoundError:
        logger.error(f"Users CSV file not found: {USERS_CSV}")
        return
    except Exception as e:
        logger.error(f"Error reading users CSV: {e}")
        return

    logger.info(f"Provisioning {len(users)} users against {BASE_URL}...")
    logger.info(f"Concurrency limit: {SEED_CONCURRENCY}")

    # Semaphore zur Kontrolle der parallelen Requests
    semaphore = asyncio.Semaphore(SEED_CONCURRENCY)

    async def provision_user(client, user):
        async with semaphore:
            user_email = user.get("email", "unknown")
            user_name = user.get("name", "unknown")

            try:
                # Token refreshen
                resp = requests.post(
                    f"https://securetoken.googleapis.com/v1/token?key={FIREBASE_API_KEY}",
                    data={"grant_type": "refresh_token", "refresh_token": user["refresh_token"]},
                    timeout=15
                )
                if resp.status_code != 200:
                    logger.warning(f"Token refresh failed for {user_email}: {resp.status_code}")
                    return False

                id_token = resp.json()["id_token"]
                logger.debug(f"Token refreshed for {user_email}")

            except Exception as e:
                logger.error(f"Token refresh error for {user_email}: {e}")
                return False

            try:
                r = await client.post(
                    f"{BASE_URL}/api/users/provision",
                    headers={"Authorization": f"Bearer {id_token}", "Content-Type": "application/json"},
                    json={"name": user_name},
                    timeout=20
                )

                if r.status_code in (200, 201):
                    logger.info(f"Successfully provisioned user: {user_email}")
                    return True
                else:
                    logger.warning(f"Provision failed for {user_email}: {r.status_code} - {r.text[:200]}")
                    return False

            except Exception as e:
                logger.error(f"Provision request error for {user_email}: {e}")
                return False

    async with httpx.AsyncClient() as client:
        results = await asyncio.gather(*[provision_user(client, u) for u in users])

    successful = sum(1 for r in results if r)
    failed = len(users) - successful

    logger.info(f"Provisioning complete: {successful}/{len(users)} successful")
    if failed > 0:
        logger.warning(f"{failed} users failed to provision")


if __name__ == "__main__":
    try:
        asyncio.run(provision_all())
        logger.info("Provisioning task completed successfully")
    except KeyboardInterrupt:
        logger.info("Provisioning interrupted by user")
    except Exception as e:
        logger.error(f"Fatal error during provisioning: {e}", exc_info=True)

