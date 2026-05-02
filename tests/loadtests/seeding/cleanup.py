# cleanup.py
import asyncio
import re

import firebase_admin
from firebase_admin import auth, credentials

from config.settings import FIREBASE_PROJECT_ID

# Pattern für Test-User-Emails (passe an dein Naming an)
TEST_USER_EMAIL_PATTERN = re.compile(r"^user_\d+@loadtest\.com$")


def init_firebase():
    """Initialisiert Firebase Admin SDK via ADC."""
    if not firebase_admin._apps:
        cred = credentials.ApplicationDefault()
        firebase_admin.initialize_app(cred, {
            "projectId": FIREBASE_PROJECT_ID,
        })


def list_test_users():
    """Listet alle User in Firebase und filtert Test-User."""
    test_users = []
    page = auth.list_users()
    while page:
        for user in page.users:
            if user.email and TEST_USER_EMAIL_PATTERN.match(user.email):
                test_users.append(user.uid)
        page = page.get_next_page()
    return test_users


def delete_users_batch(uids: list[str]):
    """Löscht User in Batches von 1000 (Firebase-Limit)."""
    batch_size = 1000
    total_deleted = 0
    total_failed = 0

    for i in range(0, len(uids), batch_size):
        batch = uids[i:i + batch_size]
        result = auth.delete_users(batch)
        total_deleted += result.success_count
        total_failed += result.failure_count

        if result.errors:
            for err in result.errors[:5]:  # nur die ersten 5 zeigen
                print(f"[WARN] Failed to delete uid={batch[err.index]}: {err.reason}")

    return total_deleted, total_failed


def main():
    init_firebase()

    print("[INFO] Listing test users in Firebase...")
    uids = list_test_users()
    print(f"[INFO] Found {len(uids)} test users")

    if not uids:
        print("[INFO] Nothing to delete")
        return

    confirm = input(f"Delete {len(uids)} users? Type 'DELETE' to confirm: ")
    if confirm != "DELETE":
        print("[INFO] Aborted")
        return

    deleted, failed = delete_users_batch(uids)
    print(f"[INFO] Deleted {deleted}, failed {failed}")


if __name__ == "__main__":
    main()