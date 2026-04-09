"""
Load Test Scenario Configuration

Definiert User-Gruppen mit hierarchischen Task-Verteilungen.
Erlaubt komplexe Szenarien wie:
- 5000 Users total
  - 30% führen "browsing" durch
    - Von diesen 30%: 15% führen auch "search" durch
  - 70% führen "booking" durch
    - Von diesen 70%: 60% führen auch "payment" durch
"""

from dataclasses import dataclass, field
from typing import Dict, List, Optional, Optional


@dataclass
class Task:
    """Definition einer einzelnen Task/Operation"""
    name: str
    weight: int = 1
    enabled: bool = True

    def __post_init__(self):
        if self.weight < 1:
            raise ValueError(f"Task weight muss >= 1 sein, got {self.weight}")


@dataclass
class UserGroup:
    """
    Eine Gruppe von Benutzern mit spezifischen Task-Verteilungen.

    Beispiel:
        UserGroup(
            name="browsers",
            percentage=30,  # 30% aller Users
            tasks={"list_trips": 8, "get_trip": 4},
            sub_groups=[
                UserGroup(
                    name="searchers",
                    percentage=15,  # 15% der "browsers" (also 4.5% gesamt)
                    tasks={"search_trips": 5}
                )
            ]
        )
    """
    name: str
    percentage: float  # Prozentsatz der Parent-Gruppe (oder 100 für Root)
    tasks: Dict[str, int]  # {task_name: weight}
    sub_groups: List['UserGroup'] = field(default_factory=list)

    def __post_init__(self):
        if not 0 < self.percentage <= 100:
            raise ValueError(f"Percentage muss zwischen 0 und 100 sein, got {self.percentage}")

        for task_name, weight in self.tasks.items():
            if weight < 1:
                raise ValueError(f"Task weight muss >= 1 sein, got {weight}")


@dataclass
class LoadTestScenario:
    """Komplette Load-Test Konfiguration"""
    name: str
    description: str = ""
    total_users: int = 100
    spawn_rate: int = 10  # Users pro Sekunde
    wait_time_min: int = 1
    wait_time_max: int = 3
    duration_minutes: int = 5
    user_groups: List[UserGroup] = field(default_factory=list)

    def __post_init__(self):
        if self.total_users < 1:
            raise ValueError("total_users muss >= 1 sein")
        if self.spawn_rate < 1:
            raise ValueError("spawn_rate muss >= 1 sein")
        if self.wait_time_min < 0 or self.wait_time_max < self.wait_time_min:
            raise ValueError("Ungültige wait_time Werte")

        # Validiere dass Prozentsätze addieren sich zu ~100%
        total_percentage = sum(group.percentage for group in self.user_groups)
        if abs(total_percentage - 100) > 0.1:
            raise ValueError(
                f"User Group Prozentsätze addieren sich zu {total_percentage}%, "
                f"sollten aber 100% sein"
            )

    def get_user_distribution(self) -> Dict[str, int]:
        """
        Berechnet die genaue User-Verteilung über alle Gruppen und Sub-Gruppen.

        Jede Gruppe wird separat gezählt, sowohl Parent als auch Sub-Groups.
        Wenn eine Parent-Group Sub-Groups hat, wird die Parent-Group NICHT
        automatisch mit den restlichen Usern registriert (da diese die Sub-Groups übernehmen).

        Returns:
            Dict mit User-Counts pro Scenario (z.B. {"browsers": 30, "browsers_power_users": 6})
        """
        distribution = {}

        def process_group(group: UserGroup, parent_count: int, prefix: str = ""):
            group_name = f"{prefix}{group.name}" if prefix else group.name
            group_count = int(parent_count * group.percentage / 100)

            # Registriere IMMER die Gruppe
            distribution[group_name] = group_count

            # Verarbeite Sub-Gruppen
            remaining_count = group_count
            for sub_group in group.sub_groups:
                sub_count = int(group_count * sub_group.percentage / 100)
                # Sub-Gruppen werden auch registriert
                process_group(sub_group, group_count, f"{group_name}_")
                remaining_count -= sub_count

            # Wenn es Sub-Gruppen gibt, aktualisiere die Parent-Group zu nur den nicht-zugeordneten Usern
            if group.sub_groups:
                distribution[group_name] = remaining_count

        for group in self.user_groups:
            process_group(group, self.total_users)

        return distribution

    def get_task_distribution(self, group_path: str) -> Dict[str, int]:
        """
        Gibt die Task-Verteilung für eine spezifische User-Gruppe zurück.

        Args:
            group_path: z.B. "browsers" oder "browsers_power_users"

        Returns:
            Dict mit {task_name: weight}
        """
        parts = group_path.split("_")

        # Finde die Gruppe
        def find_group(name: str, groups_list: list) -> Optional['UserGroup']:
            for group in groups_list:
                if group.name == name:
                    return group
            return None

        # Traverse die hierarchy
        current_groups = self.user_groups
        for i, part in enumerate(parts):
            group = find_group(part, current_groups)
            if not group:
                raise ValueError(f"Group path '{group_path}' nicht gefunden (part: '{part}')")

            if i == len(parts) - 1:
                # Last part - return tasks
                return group.tasks.copy()
            else:
                # More parts to traverse
                if not group.sub_groups:
                    raise ValueError(f"Group '{part}' hat keine Sub-Groups, aber path hat mehr parts")
                current_groups = group.sub_groups

        raise ValueError(f"Group path '{group_path}' nicht gefunden")


# ============================================================================
# VORDEFINIERTEN SZENARIEN
# ============================================================================

SCENARIO_LIGHT = LoadTestScenario(
    name="light",
    description="Schneller Sanity Check - Entwicklung",
    total_users=10,
    spawn_rate=2,
    duration_minutes=5,
    user_groups=[
        UserGroup(
            name="all_users",
            percentage=100,
            tasks={
                "register_user": 1,
                "login_user": 3,
                "get_user_profile": 5,
                "update_user_profile": 2,
                "list_trips": 8,
                "create_trip": 6,
                "get_trip": 4,
                "update_trip": 3,
                "create_location": 5,
                "list_locations": 4,
                "get_location": 3,
                "update_location": 3,
                "create_activity": 5,
                "list_activities": 4,
                "get_activity": 3,
                "update_activity": 2,
                "health_check": 1,
            }
        )
    ]
)

SCENARIO_NORMAL = LoadTestScenario(
    name="normal",
    description="Baseline Performance - Standard Last",
    total_users=50,
    spawn_rate=5,
    duration_minutes=10,
    user_groups=[
        UserGroup(
            name="all_users",
            percentage=100,
            tasks={
                "list_trips": 8,
                "create_trip": 6,
                "get_trip": 4,
                "update_trip": 3,
                "create_location": 5,
                "list_locations": 4,
                "get_location": 3,
                "update_location": 3,
                "create_activity": 5,
                "list_activities": 4,
                "get_activity": 3,
                "update_activity": 2,
                "get_user_profile": 5,
            }
        )
    ]
)

SCENARIO_SEGMENTED = LoadTestScenario(
    name="segmented",
    description="Segmentierte User - Verschiedene Verhalten",
    total_users=100,
    spawn_rate=5,
    duration_minutes=10,
    user_groups=[
        UserGroup(
            name="browsers",
            percentage=30,
            tasks={
                "list_trips": 10,
                "get_trip": 8,
                "list_locations": 6,
                "get_location": 5,
                "list_activities": 4,
            },
            sub_groups=[
                UserGroup(
                    name="powerusers",  # No underscore!
                    percentage=20,  # 20% der browsers = 6 Users
                    tasks={
                        "search_trips": 5,
                        "create_trip": 3,
                        "create_location": 3,
                    }
                )
            ]
        ),
        UserGroup(
            name="planners",
            percentage=50,
            tasks={
                "create_trip": 8,
                "create_location": 6,
                "create_activity": 7,
                "update_trip": 4,
                "get_trip": 3,
            },
            sub_groups=[
                UserGroup(
                    name="collaborators",  # No underscore!
                    percentage=30,  # 30% der planners = 15 Users
                    tasks={
                        "share_trip": 5,
                        "update_location": 4,
                    }
                )
            ]
        ),
        UserGroup(
            name="viewers",
            percentage=20,
            tasks={
                "list_trips": 15,
                "get_trip": 10,
                "list_activities": 8,
            }
        )
    ]
)

SCENARIO_STRESS = LoadTestScenario(
    name="stress",
    description="Stress Test - Limits finden",
    total_users=500,
    spawn_rate=20,
    duration_minutes=15,
    user_groups=[
        UserGroup(
            name="all_users",
            percentage=100,
            tasks={
                "list_trips": 8,
                "create_trip": 6,
                "get_trip": 4,
                "create_activity": 5,
                "list_activities": 4,
            }
        )
    ]
)

# Szenarien Registry
SCENARIOS = {
    "light": SCENARIO_LIGHT,
    "normal": SCENARIO_NORMAL,
    "segmented": SCENARIO_SEGMENTED,
    "stress": SCENARIO_STRESS,
}






