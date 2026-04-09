# Load Testing Framework for Trip Manager API
A flexible, hierarchical load testing framework using Locust. Define realistic user scenarios with different user groups and task distributions.
## Prerequisites
- **Python 3.10+** (see `.python-version`)
- **uv** package manager ([install](https://docs.astral.sh/uv/getting-started/installation/))
## Setup with uv
This project uses **uv** for fast, reliable dependency management. It's much faster than traditional pip/poetry setups.
### 1. Install uv
```bash
# macOS/Linux
curl -LsSf https://astral.sh/uv/install.sh | sh
# Windows (PowerShell)
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"
# Or via package managers
brew install uv          # macOS
choco install uv         # Windows (Chocolatey)
cargo install uv          # Any OS (Rust toolchain)
```
Verify installation:
```bash
uv --version
```
### 2. Initialize Project Dependencies
```bash
cd tests/loadtests
# Install all dependencies (creates .venv automatically)
uv sync
```
This command:
1. Reads `pyproject.toml` and `requirements.txt`
2. Creates a virtual environment in `.venv/`
3. Installs all dependencies (locust, etc.)
4. Sets up the project in development mode
### 3. Activate Virtual Environment
```bash
# macOS/Linux
source .venv/bin/activate
# Windows (CMD)
.venv\Scripts\activate
# Windows (PowerShell)
.venv\Scripts\Activate.ps1
# Verify (should show path prefix)
which python  # macOS/Linux
Get-Command python  # Windows
```
### 4. Run Commands
```bash
# Now you can use make/locust directly
make run-ui
make run-light
locust --version
```
### Alternative: Run Without Activation
You can run commands directly with `uv` prefix (no activation needed):
```bash
uv run make run-ui
uv run locust --version
uv run python -c "print('Hello from uv')"
```
### Update Dependencies
```bash
# Update all dependencies to latest compatible versions
uv sync --upgrade
# Update specific package
uv pip install --upgrade locust
# Check for updates without installing
uv pip list --outdated
```
### Troubleshooting uv
**"command not found: uv"**
- Reinstall uv using the curl command above
- Add uv to PATH manually if needed
**Virtual environment not created**
```bash
uv venv .venv  # Manually create venv
uv sync        # Then sync dependencies
```
**Wrong Python version**
```bash
# Specify Python version explicitly
uv sync --python 3.10
```
---
## Quick Start (2 Minutes)
```bash
cd tests/loadtests
# Initial setup (one-time)
make setup
# Start web UI
make run-ui
# → Open http://localhost:8089
```
## Features
✅ **Hierarchical User Segmentation**
- Define multiple user groups with percentage-based distribution
- Support for sub-groups (e.g., "20% of 30% are power users")
- Automatic user count calculation
✅ **Flexible Task Distribution**
- Different tasks per user group
- Task weights for frequency control
- 25+ pre-built tasks available
✅ **Pre-built Scenarios**
- `light` - 10 users (quick testing)
- `normal` - 50 users (standard load)
- `segmented` - 100 users with hierarchical groups
- `stress` - 500 users (limits testing)
✅ **Zero Configuration** 
- Works out of the box
- Environment variable support
- Easy to extend
## Usage Examples
### Web UI (Recommended)
```bash
# Default scenario
make run-ui
# With specific scenario
SCENARIO=segmented make run-ui
SCENARIO=custom make run-ui
```
### Command Line
```bash
# Light test
make run-light
# Normal test with custom scenario
SCENARIO=segmented make run-light
# Direct Locust command
SCENARIO=custom locust -f locustfile.py \
  --host http://localhost:8000/api \
  --headless --users 100 --spawn-rate 10 --run-time 10m
```
### Analysis
```bash
# Generate performance report
make analyze CSV_FILE=results/light_<timestamp>_stats.csv
```
## Creating Scenarios
### Define in `scenario.py`
```python
SCENARIO_DEMO = LoadTestScenario(
    name="demo",
    total_users=1000,
    spawn_rate=50,
    duration_minutes=15,
    user_groups=[
        # 60% casual browsers
        UserGroup(
            name="casual",
            percentage=60,
            tasks={
                "list_trips": 15,
                "get_trip": 10,
            },
            sub_groups=[
                # 20% of casual users are frequent
                UserGroup(
                    name="frequent",
                    percentage=20,
                    tasks={
                        "create_trip": 3,
                        "create_activity": 2,
                    }
                )
            ]
        ),
        # 40% active creators
        UserGroup(
            name="active",
            percentage=40,
            tasks={
                "create_trip": 8,
                "create_activity": 6,
                "create_location": 5,
            }
        )
    ]
)
# Register it
SCENARIOS["demo"] = SCENARIO_DEMO
```
### Use it
```bash
SCENARIO=demo make run-ui
```
## Available Tasks
All tasks that can be used in `tasks={}`:
| Category | Tasks |
|----------|-------|
| **Auth** | `register_user`, `login_user` |
| **Profile** | `get_user_profile`, `update_user_profile` |
| **Trips** | `list_trips`, `create_trip`, `get_trip`, `update_trip`, `delete_trip` |
| **Locations** | `create_location`, `list_locations`, `get_location`, `update_location`, `delete_location` |
| **Activities** | `create_activity`, `list_activities`, `get_activity`, `update_activity`, `delete_activity` |
| **Other** | `health_check`, `search_trips`, `share_trip` |
## Understanding Distribution
```
Total: 1000 users
├─ casual (60% = 600 users)
│  ├─ regular casual (80% = 480 users)
│  └─ frequent (20% = 120 users)
└─ active (40% = 400 users)
```
**Important:** Group names should NOT contain underscores (they're used as path separators).
## Task Weights
Higher weights = executed more frequently:
```python
tasks={
    "list_trips": 10,      # 10x as frequent as...
    "create_trip": 5,      # 5x as frequent as...
    "delete_trip": 1,      # baseline
}
```
## Pre-built Examples
Seven ready-to-use scenarios in `scenario_examples.py`:
- `b2b` - Business Use Case (many planners, few viewers)
- `mobile` - Mobile App (fast interactions)
- `office_hours` - Office Hours (normal business load)
- `peak_load` - Peak Hour (all systems loaded)
- `feature_test` - New Feature Testing (control vs test group)
- `geographic` - Geographic Distribution (multi-region)
- `regression` - Regression Testing (compare implementations)
Copy any and customize!
## Commands Reference
```bash
# Setup
make install              # Install dependencies (uv sync)
make setup               # Create .env file
# Testing
make run-ui              # Web UI (interactive)
make run-light           # 10 users, 5 min
make run-normal          # 50 users, 10 min
make run-stress          # 500 users, 15 min
make docker-run          # Run in container
# Analysis
make analyze CSV_FILE=...  # Generate report
# Help
make help                # Show all commands
make clean               # Remove results
```
## Environment Variables
```bash
# Set scenario
SCENARIO=segmented make run-ui
# API endpoint
API_URL=http://api.example.com:8000/api make run-ui
# Or in .env file
cat > .env << EOF
API_URL=http://localhost:8000/api
