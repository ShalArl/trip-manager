"""
Generic Load Tests for Trip Manager API

Dieses Modul definiert flexible, konfigurierbare Load-Test-Szenarien.
User-Gruppen und Task-Verteilungen können in scenario.py definiert werden.
"""

import os
import sys
from locust import HttpUser, TaskSet, task, between, events
from datetime import datetime, timedelta
from scenario import SCENARIOS, SCENARIO_LIGHT


# Hole Scenario aus Umgebungsvariable oder CLI
SCENARIO_NAME = os.getenv("SCENARIO", "light")
SCENARIO = SCENARIOS.get(SCENARIO_NAME, SCENARIO_LIGHT)

print(f"\n{'='*80}")
print(f"🎯 Scenario: {SCENARIO.name}")
print(f"   {SCENARIO.description}")
print(f"   Users: {SCENARIO.total_users}, Spawn Rate: {SCENARIO.spawn_rate}/sec")
print(f"{'='*80}\n")


class APITasks(TaskSet):
    """
    Basis Task-Set mit allen möglichen Operationen.
    Wird von User-Klassen geerbt, die spezifische Tasks aktivieren.
    """

    def on_start(self):
        """Called when a simulated user starts"""
        self.token = None
        self.user_id = None
        self.trip_id = None
        self.location_id = None
        self.activity_id = None
        self.email = None
        self.register_user()

    def register_user(self):
        """Register a new user"""
        timestamp = datetime.now().isoformat()
        self.email = f"testuser_{timestamp.replace(':', '-').replace('.', '_')}@example.com"

        payload = {
            "firstName": "Test",
            "lastName": "User",
            "email": self.email,
            "password": "TestPassword123!",
        }

        with self.client.post(
            "/auth/register",
            json=payload,
            catch_response=True,
            name="/auth/register",
        ) as response:
            if response.status_code == 201:
                data = response.json()
                self.token = data.get("token")
                self.user_id = data.get("user", {}).get("id")
                response.success()
            else:
                response.failure(f"Register failed: {response.status_code}")

    def login_user(self):
        """Login with existing user credentials"""
        if not self.email:
            return

        payload = {
            "email": self.email,
            "password": "TestPassword123!",
        }

        with self.client.post(
            "/auth/login",
            json=payload,
            catch_response=True,
            name="/auth/login",
        ) as response:
            if response.status_code == 200:
                data = response.json()
                self.token = data.get("token")
                response.success()
            else:
                response.failure(f"Login failed: {response.status_code}")

    def _get_auth_headers(self):
        """Get authorization headers with token"""
        if not self.token:
            self.login_user()

        return {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json",
        }

    # ========== AUTH TASKS ==========

    def register_user_task(self):
        """Register a new user (Task wrapper)"""
        self.register_user()

    def login_user_task(self):
        """Login with existing user credentials (Task wrapper)"""
        self.login_user()

    # ========== USER PROFILE TASKS ==========

    def get_user_profile(self):
        """Get current user profile"""
        headers = self._get_auth_headers()

        with self.client.get(
            "/users/me",
            headers=headers,
            catch_response=True,
            name="/users/me [GET]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

    def update_user_profile(self):
        """Update current user profile"""
        headers = self._get_auth_headers()

        payload = {
            "firstName": "Updated",
            "lastName": "Name",
        }

        with self.client.put(
            "/users/me",
            json=payload,
            headers=headers,
            catch_response=True,
            name="/users/me [PUT]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Update failed: {response.status_code}")

    # ========== TRIP TASKS ==========

    def list_trips(self):
        """List all trips"""
        headers = self._get_auth_headers()

        with self.client.get(
            "/trips",
            headers=headers,
            catch_response=True,
            name="/trips [GET]",
        ) as response:
            if response.status_code == 200:
                data = response.json()
                if isinstance(data, list) and len(data) > 0:
                    self.trip_id = data[0].get("id")
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

    def create_trip(self):
        """Create a new trip"""
        headers = self._get_auth_headers()

        now = datetime.now()
        start_date = now.strftime("%Y-%m-%d")
        end_date = (now + timedelta(days=7)).strftime("%Y-%m-%d")

        payload = {
            "title": f"Trip to {now.isoformat()}",
            "description": "Test trip for load testing",
            "startDate": start_date,
            "endDate": end_date,
            "status": "PLANNED",
        }

        with self.client.post(
            "/trips",
            json=payload,
            headers=headers,
            catch_response=True,
            name="/trips [POST]",
        ) as response:
            if response.status_code == 201:
                data = response.json()
                self.trip_id = data.get("id")
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Create trip failed: {response.status_code}")

    def get_trip(self):
        """Get a specific trip"""
        if not self.trip_id:
            self.list_trips()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()

        with self.client.get(
            f"/trips/{self.trip_id}",
            headers=headers,
            catch_response=True,
            name="/trips/[tripId] [GET]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

    def update_trip(self):
        """Update a trip"""
        if not self.trip_id:
            self.create_trip()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()

        payload = {
            "title": f"Updated Trip {datetime.now().isoformat()}",
            "description": "Updated description",
            "status": "IN_PROGRESS",
        }

        with self.client.put(
            f"/trips/{self.trip_id}",
            json=payload,
            headers=headers,
            catch_response=True,
            name="/trips/[tripId] [PUT]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Update trip failed: {response.status_code}")

    def delete_trip(self):
        """Delete a trip"""
        if not self.trip_id:
            self.create_trip()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()
        trip_id = self.trip_id
        self.trip_id = None  # Clear for next test

        with self.client.delete(
            f"/trips/{trip_id}",
            headers=headers,
            catch_response=True,
            name="/trips/[tripId] [DELETE]",
        ) as response:
            if response.status_code in [200, 204]:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Delete trip failed: {response.status_code}")

    # ========== LOCATION TASKS ==========

    def create_location(self):
        """Create a location in a trip"""
        if not self.trip_id:
            self.create_trip()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()

        payload = {
            "name": f"Location {datetime.now().isoformat()}",
            "description": "Test location",
            "latitude": 48.8566,
            "longitude": 2.3522,
            "country": "France",
            "city": "Paris",
        }

        with self.client.post(
            f"/trips/{self.trip_id}/locations",
            json=payload,
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/locations [POST]",
        ) as response:
            if response.status_code == 201:
                data = response.json()
                self.location_id = data.get("id")
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Create location failed: {response.status_code}")

    def list_locations(self):
        """List locations for a trip"""
        if not self.trip_id:
            self.list_trips()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()

        with self.client.get(
            f"/trips/{self.trip_id}/locations",
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/locations [GET]",
        ) as response:
            if response.status_code == 200:
                data = response.json()
                if isinstance(data, list) and len(data) > 0:
                    self.location_id = data[0].get("id")
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

    def get_location(self):
        """Get a specific location"""
        if not self.location_id:
            self.list_locations()
            if not self.location_id:
                return

        headers = self._get_auth_headers()

        with self.client.get(
            f"/locations/{self.location_id}",
            headers=headers,
            catch_response=True,
            name="/locations/[locationId] [GET]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

    def update_location(self):
        """Update a location"""
        if not self.trip_id or not self.location_id:
            self.create_location()
            if not self.location_id:
                return

        headers = self._get_auth_headers()

        payload = {
            "name": f"Updated Location {datetime.now().isoformat()}",
            "description": "Updated description",
            "latitude": 48.8566,
            "longitude": 2.3522,
        }

        with self.client.put(
            f"/trips/{self.trip_id}/locations/{self.location_id}",
            json=payload,
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/locations/[locationId] [PUT]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Update location failed: {response.status_code}")

    def delete_location(self):
        """Delete a location"""
        if not self.trip_id or not self.location_id:
            self.create_location()
            if not self.location_id:
                return

        headers = self._get_auth_headers()
        location_id = self.location_id
        self.location_id = None

        with self.client.delete(
            f"/trips/{self.trip_id}/locations/{location_id}",
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/locations/[locationId] [DELETE]",
        ) as response:
            if response.status_code in [200, 204]:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Delete location failed: {response.status_code}")

    # ========== ACTIVITY TASKS ==========

    def create_activity(self):
        """Create an activity in a trip"""
        if not self.trip_id:
            self.create_trip()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()

        payload = {
            "title": f"Activity {datetime.now().isoformat()}",
            "description": "Test activity",
            "category": "SIGHTSEEING",
            "date": datetime.now().strftime("%Y-%m-%d"),
            "startTime": "10:00",
            "endTime": "12:00",
        }

        with self.client.post(
            f"/trips/{self.trip_id}/activities",
            json=payload,
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/activities [POST]",
        ) as response:
            if response.status_code == 201:
                data = response.json()
                self.activity_id = data.get("id")
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Create activity failed: {response.status_code}")

    def list_activities(self):
        """List activities for a trip"""
        if not self.trip_id:
            self.list_trips()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()

        with self.client.get(
            f"/trips/{self.trip_id}/activities",
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/activities [GET]",
        ) as response:
            if response.status_code == 200:
                data = response.json()
                if isinstance(data, list) and len(data) > 0:
                    self.activity_id = data[0].get("id")
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

    def get_activity(self):
        """Get a specific activity"""
        if not self.activity_id:
            self.list_activities()
            if not self.activity_id:
                return

        headers = self._get_auth_headers()

        with self.client.get(
            f"/activities/{self.activity_id}",
            headers=headers,
            catch_response=True,
            name="/activities/[activityId] [GET]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

    def update_activity(self):
        """Update an activity"""
        if not self.trip_id or not self.activity_id:
            self.create_activity()
            if not self.activity_id:
                return

        headers = self._get_auth_headers()

        payload = {
            "title": f"Updated Activity {datetime.now().isoformat()}",
            "description": "Updated description",
            "category": "DINING",
            "startTime": "18:00",
            "endTime": "20:00",
        }

        with self.client.put(
            f"/trips/{self.trip_id}/activities/{self.activity_id}",
            json=payload,
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/activities/[activityId] [PUT]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Update activity failed: {response.status_code}")

    def delete_activity(self):
        """Delete an activity"""
        if not self.trip_id or not self.activity_id:
            self.create_activity()
            if not self.activity_id:
                return

        headers = self._get_auth_headers()
        activity_id = self.activity_id
        self.activity_id = None

        with self.client.delete(
            f"/trips/{self.trip_id}/activities/{activity_id}",
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/activities/[activityId] [DELETE]",
        ) as response:
            if response.status_code in [200, 204]:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Delete activity failed: {response.status_code}")

    # ========== OTHER TASKS ==========

    def health_check(self):
        """Perform a health check"""
        with self.client.get(
            "/health",
            catch_response=True,
            name="/health",
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Health check failed: {response.status_code}")

    def search_trips(self):
        """Search for trips (simulate complex search)"""
        headers = self._get_auth_headers()

        with self.client.get(
            "/trips?search=test",
            headers=headers,
            catch_response=True,
            name="/trips?search=test [GET]",
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            else:
                response.failure(f"Search failed: {response.status_code}")

    def share_trip(self):
        """Share a trip with another user (simulate sharing)"""
        if not self.trip_id:
            self.create_trip()
            if not self.trip_id:
                return

        headers = self._get_auth_headers()

        with self.client.post(
            f"/trips/{self.trip_id}/share",
            json={"email": "friend@example.com"},
            headers=headers,
            catch_response=True,
            name="/trips/[tripId]/share [POST]",
        ) as response:
            if response.status_code in [200, 201]:
                response.success()
            elif response.status_code == 401:
                self.login_user()
            elif response.status_code == 404:
                # Endpoint existiert vielleicht nicht, das ist OK
                response.success()
            else:
                response.failure(f"Share failed: {response.status_code}")


def create_user_class(group_path: str, tasks: dict) -> type:
    """
    Erstellt dynamisch eine User-Klasse für eine spezifische User-Gruppe.

    Args:
        group_path: Name der Gruppe (z.B. "browsers" oder "browsers_searchers")
        tasks: Dict mit {task_name: weight}

    Returns:
        Eine Locust HttpUser Klasse
    """

    # Erstelle Task-Methoden mit @task Dekorator
    task_methods = {}

    for task_name, weight in tasks.items():
        # Finde die echte Methode im APITasks
        if hasattr(APITasks, task_name):
            method = getattr(APITasks, task_name)
            # Decorator @task(weight) zu Method hinzufügen
            decorated_method = task(weight)(method)
            task_methods[task_name] = decorated_method

    # Erstelle neue TaskSet Klasse mit diesen Tasks
    UserTaskSet = type(
        f'{group_path.title().replace("_", "")}TaskSet',
        (APITasks,),
        task_methods
    )

    # Erstelle neue User Klasse mit diesem TaskSet
    UserClass = type(
        f'{group_path.title().replace("_", "")}User',
        (HttpUser,),
        {
            'tasks': [UserTaskSet],
            'wait_time': between(SCENARIO.wait_time_min, SCENARIO.wait_time_max),
            'host': os.getenv("API_URL", "http://localhost:8000/api"),
        }
    )

    return UserClass


# Erstelle dynamisch alle User-Klassen basierend auf SCENARIO
def setup_user_classes():
    """Setup alle User-Klassen für das Scenario"""
    user_distribution = SCENARIO.get_user_distribution()

    print(f"\n📊 User Distribution:")
    print(f"{'─'*80}")
    for group_path, user_count in user_distribution.items():
        tasks = SCENARIO.get_task_distribution(group_path)
        print(f"  {group_path:<40} {user_count:>6} Users | Tasks: {len(tasks)}")
    print(f"{'─'*80}\n")

    # Erstelle User-Klassen
    user_classes = {}
    for group_path, user_count in user_distribution.items():
        tasks = SCENARIO.get_task_distribution(group_path)
        UserClass = create_user_class(group_path, tasks)
        user_classes[group_path] = UserClass

    return user_classes


# Setup user classes
try:
    CONFIGURED_USER_CLASSES = setup_user_classes()
    # Exportiere als globals für Locust
    for name, user_class in CONFIGURED_USER_CLASSES.items():
        globals()[user_class.__name__] = user_class
except Exception as e:
    print(f"❌ Error setting up user classes: {e}")
    sys.exit(1)


@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    """Called when the test starts"""
    print("\n" + "=" * 80)
    print(f"🚀 Trip Manager Load Test Started - Scenario: {SCENARIO.name}")
    print(f"   Target: {os.getenv('API_URL', 'http://localhost:8000/api')}")
    print(f"   Total Users: {SCENARIO.total_users}")
    print(f"   Spawn Rate: {SCENARIO.spawn_rate}/sec")
    print(f"   Duration: {SCENARIO.duration_minutes} minutes")
    print("=" * 80 + "\n")


@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    """Called when the test stops"""
    print("\n" + "=" * 80)
    print("✅ Trip Manager Load Test Completed")
    print("=" * 80 + "\n")

