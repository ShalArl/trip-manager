#!/bin/bash

# Load Testing Script for Trip Manager API
#
# Usage:
#   ./run.sh [light|normal|stress|ui] [scenario]
#
# Examples:
#   ./run.sh light                    # 10 users, 5 minute test
#   ./run.sh light segmented          # With segmented scenario
#   ./run.sh ui                       # Web UI (interactive)
#   ./run.sh ui complex               # Web UI with complex scenario

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Default API URL
API_URL="${API_URL:-http://localhost:8000/api}"

# Get Scenario (from env or argument)
SCENARIO="${SCENARIO:-${2:-light}}"
export SCENARIO

# Color output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Trip Manager Load Test${NC}"
echo "Target API: $API_URL"
echo "Scenario: $SCENARIO"
echo ""

# Determine test profile
PROFILE=${1:-ui}

case $PROFILE in
    light)
        echo -e "${GREEN}📊 Light Load Test (10 users, 5 minutes)${NC}"
        locust -f locustfile.py \
            --host "$API_URL" \
            --users 10 \
            --spawn-rate 2 \
            --run-time 5m \
            --headless \
            --csv=results/light_$(date +%s)
        ;;
    normal)
        echo -e "${GREEN}📊 Normal Load Test (50 users, 10 minutes)${NC}"
        locust -f locustfile.py \
            --host "$API_URL" \
            --users 50 \
            --spawn-rate 5 \
            --run-time 10m \
            --headless \
            --csv=results/normal_$(date +%s) \
            --html=results/normal_$(date +%s).html
        ;;
    stress)
        echo -e "${GREEN}📊 Stress Test (500 users, ramp up to breaking point)${NC}"
        locust -f locustfile.py \
            --host "$API_URL" \
            --users 500 \
            --spawn-rate 20 \
            --run-time 15m \
            --headless \
            --csv=results/stress_$(date +%s) \
            --html=results/stress_$(date +%s).html
        ;;
    ui)
        echo -e "${GREEN}🎨 Web UI Mode (http://localhost:8089)${NC}"
        echo "Open http://localhost:8089 in your browser"
        locust -f locustfile.py --host "$API_URL"
        ;;
    *)
        echo "Usage: $0 [light|normal|stress|ui]"
        echo ""
        echo "Profiles:"
        echo "  light      10 users, 5 minute test"
        echo "  normal     50 users, 10 minute test"
        echo "  stress     500 users, 15 minute test (find breaking point)"
        echo "  ui         Web UI (interactive, default)"
        exit 1
        ;;
esac

