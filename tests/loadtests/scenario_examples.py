"""
scenario_examples.py

Zusätzliche Szenarien-Beispiele die du verwenden oder als Vorlage nutzen kannst.

Kopiere einen dieser zu scenario.py und registriere ihn in SCENARIOS dict.
"""

from scenario import LoadTestScenario, UserGroup


# ============================================================================
# BEISPIEL 1: B2B Szenario (Business Use)
# ============================================================================

SCENARIO_B2B = LoadTestScenario(
    name="b2b",
    description="Business Use Case - Viele Planner, wenige Viewer",
    total_users=500,
    spawn_rate=10,
    duration_minutes=10,
    user_groups=[
        # Viele Trip-Ersteller
        UserGroup(
            name="planners",
            percentage=70,  # 350 Users
            tasks={
                "create_trip": 8,
                "create_location": 7,
                "create_activity": 6,
                "update_trip": 4,
                "list_trips": 3,
            },
            sub_groups=[
                # Ein Teil sind Collaborators (teilen Trips)
                UserGroup(
                    name="collaborators",
                    percentage=40,  # 40% von 70% = 140 Users
                    tasks={
                        "share_trip": 5,
                        "update_location": 4,
                    }
                )
            ]
        ),
        # Normale Viewer
        UserGroup(
            name="viewers",
            percentage=25,  # 125 Users
            tasks={
                "list_trips": 12,
                "get_trip": 8,
                "list_activities": 5,
            }
        ),
        # Ein paar Power-User / Admins
        UserGroup(
            name="admins",
            percentage=5,  # 25 Users
            tasks={
                "delete_trip": 3,
                "delete_activity": 2,
                "health_check": 2,
                "update_trip": 2,
            }
        )
    ]
)


# ============================================================================
# BEISPIEL 2: Mobile App Szenario
# ============================================================================

SCENARIO_MOBILE = LoadTestScenario(
    name="mobile",
    description="Mobile App Users - Schnelle Interaktionen, kurze Sessions",
    total_users=1000,
    spawn_rate=50,
    wait_time_min=0,      # Schneller (Mobile!)
    wait_time_max=2,
    duration_minutes=15,
    user_groups=[
        # Casual Mobile User (nur Browse)
        UserGroup(
            name="casual_mobile",
            percentage=50,  # 500 Users
            tasks={
                "list_trips": 15,      # Sehr häufig
                "get_trip": 10,
                "list_activities": 8,
                "health_check": 1,
            }
        ),
        # Active Mobile User (auch erstellen)
        UserGroup(
            name="active_mobile",
            percentage=35,  # 350 Users
            tasks={
                "list_trips": 8,
                "create_activity": 6,
                "create_location": 5,
                "get_trip": 4,
            }
        ),
        # Power User Mobile
        UserGroup(
            name="power_mobile",
            percentage=15,  # 150 Users
            tasks={
                "create_trip": 5,
                "create_activity": 8,
                "update_trip": 4,
                "share_trip": 3,
            }
        )
    ]
)


# ============================================================================
# BEISPIEL 3: Office Hours Szenario
# ============================================================================

SCENARIO_OFFICE_HOURS = LoadTestScenario(
    name="office_hours",
    description="Typical Office Hours - Normal Business Load",
    total_users=300,
    spawn_rate=15,
    wait_time_min=2,      # Weniger gehetzt
    wait_time_max=5,
    duration_minutes=8,
    user_groups=[
        # Managers/Planners
        UserGroup(
            name="managers",
            percentage=30,  # 90 Users
            tasks={
                "create_trip": 6,
                "update_trip": 5,
                "create_activity": 5,
                "list_trips": 4,
                "share_trip": 3,
            }
        ),
        # Regular Office Workers
        UserGroup(
            name="office_workers",
            percentage=55,  # 165 Users
            tasks={
                "list_trips": 10,
                "get_trip": 8,
                "create_activity": 5,
                "get_user_profile": 3,
            }
        ),
        # Occasional Users
        UserGroup(
            name="occasional",
            percentage=15,  # 45 Users
            tasks={
                "list_trips": 15,
                "get_trip": 5,
                "health_check": 2,
            }
        )
    ]
)


# ============================================================================
# BEISPIEL 4: Peak Load Szenario
# ============================================================================

SCENARIO_PEAK_LOAD = LoadTestScenario(
    name="peak_load",
    description="Peak Hour Load - All Users Active",
    total_users=2000,
    spawn_rate=100,       # Schneller Spawn
    wait_time_min=0,      # Keine Pausen!
    wait_time_max=1,
    duration_minutes=10,
    user_groups=[
        # Massive Browse Load
        UserGroup(
            name="browsers",
            percentage=60,  # 1200 Users
            tasks={
                "list_trips": 20,      # Fast nur Reads
                "get_trip": 15,
                "list_activities": 8,
            }
        ),
        # Some Write Operations
        UserGroup(
            name="writers",
            percentage=35,  # 700 Users
            tasks={
                "create_activity": 6,
                "list_trips": 5,
                "get_trip": 4,
            }
        ),
        # Few Admin Operations
        UserGroup(
            name="admins",
            percentage=5,   # 100 Users
            tasks={
                "delete_activity": 2,
                "update_trip": 2,
                "health_check": 2,
            }
        )
    ]
)


# ============================================================================
# BEISPIEL 5: New Feature Testing
# ============================================================================

SCENARIO_FEATURE_TEST = LoadTestScenario(
    name="feature_test",
    description="Testing New Feature - Focused Load on New Operations",
    total_users=200,
    spawn_rate=10,
    duration_minutes=5,
    user_groups=[
        # Power Users Testing New Feature
        UserGroup(
            name="feature_users",
            percentage=60,  # 120 Users - Main Test Group
            tasks={
                "search_trips": 10,     # NEW - More load
                "share_trip": 8,        # NEW - More load
                "list_trips": 5,
                "create_trip": 3,
            }
        ),
        # Control Group (Old Behavior)
        UserGroup(
            name="control_group",
            percentage=40,  # 80 Users - Baseline
            tasks={
                "list_trips": 10,
                "create_trip": 6,
                "create_activity": 5,
            }
        )
    ]
)


# ============================================================================
# BEISPIEL 6: Geographic Distribution (Europe/US/Asia)
# ============================================================================

SCENARIO_GEOGRAPHIC = LoadTestScenario(
    name="geographic",
    description="Geographic Distribution Simulation",
    total_users=600,
    spawn_rate=20,
    duration_minutes=15,
    user_groups=[
        # Europe (Daytime)
        UserGroup(
            name="europe",
            percentage=40,  # 240 Users
            tasks={
                "list_trips": 8,
                "create_trip": 6,
                "create_activity": 5,
            }
        ),
        # USA (Early Morning)
        UserGroup(
            name="usa",
            percentage=40,  # 240 Users
            tasks={
                "list_trips": 10,       # More leisure browsing
                "get_trip": 7,
                "health_check": 1,
            }
        ),
        # Asia/Pacific (Night)
        UserGroup(
            name="asia",
            percentage=20,  # 120 Users - Fewer users
            tasks={
                "list_trips": 12,
                "get_trip": 8,
                "search_trips": 3,
            }
        )
    ]
)


# ============================================================================
# BEISPIEL 7: Regression Test (Alte vs Neue Version)
# ============================================================================

SCENARIO_REGRESSION = LoadTestScenario(
    name="regression",
    description="Regression Test - Compare Old vs New Implementation",
    total_users=100,
    spawn_rate=5,
    duration_minutes=10,
    user_groups=[
        # Test Old Implementation
        UserGroup(
            name="old_impl",
            percentage=50,  # 50 Users
            tasks={
                "list_trips": 8,
                "create_trip": 6,
                "get_trip": 4,
                "create_activity": 5,
            }
        ),
        # Test New Implementation
        UserGroup(
            name="new_impl",
            percentage=50,  # 50 Users
            tasks={
                "list_trips": 8,        # Same mix
                "create_trip": 6,
                "get_trip": 4,
                "create_activity": 5,
            }
        )
    ]
)


# ============================================================================
# Registriere alle Beispiel-Szenarien
# ============================================================================

EXAMPLE_SCENARIOS = {
    "b2b": SCENARIO_B2B,
    "mobile": SCENARIO_MOBILE,
    "office_hours": SCENARIO_OFFICE_HOURS,
    "peak_load": SCENARIO_PEAK_LOAD,
    "feature_test": SCENARIO_FEATURE_TEST,
    "geographic": SCENARIO_GEOGRAPHIC,
    "regression": SCENARIO_REGRESSION,
}


# Tipp: Um diese zu nutzen, kopiere die Szenarien in scenario.py:
#
# # Aus scenario_examples.py kopieren:
# from scenario_examples import SCENARIO_B2B
#
# # Dann in SCENARIOS dict hinzufügen:
# SCENARIOS["b2b"] = SCENARIO_B2B
#
# # Oder alle auf einmal:
# SCENARIOS.update(EXAMPLE_SCENARIOS)
#
# # Dann nutzen:
# SCENARIO=b2b make run-ui

