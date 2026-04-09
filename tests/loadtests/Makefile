.PHONY: help install setup run-ui run-light run-normal run-stress analyze docker-build docker-run clean

help:
	@echo "Trip Manager Load Tests - Available Commands"
	@echo "============================================"
	@echo ""
	@echo "Setup:"
	@echo "  make install        Install dependencies using uv"
	@echo "  make setup          Create .env file from template"
	@echo ""
	@echo "Running Tests:"
	@echo "  make run-ui         Start web UI (http://localhost:8089)"
	@echo "  make run-light      Run light load test (10 users, 5m)"
	@echo "  make run-normal     Run normal load test (50 users, 10m)"
	@echo "  make run-stress     Run stress test (500 users, 15m)"
	@echo ""
	@echo "Analysis:"
	@echo "  make analyze        Generate performance report"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build   Build Docker image"
	@echo "  make docker-run     Run tests in Docker container"
	@echo "  make docker-compose Run distributed tests with docker-compose"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean          Remove test results and temporary files"

install:
	@echo "📦 Installing dependencies..."
	uv sync

setup:
	@echo "📝 Creating .env file..."
	@if [ -f .env ]; then \
		echo "⚠️  .env file already exists, skipping..."; \
	else \
		cp .env.example .env; \
		echo "✅ .env file created. Edit it to configure your API URL."; \
	fi
	mkdir -p results

run-ui: setup
	@echo "🎨 Starting Locust Web UI (http://localhost:8089)..."
	@echo "   Scenario: $${SCENARIO:-light}"
	SCENARIO=$${SCENARIO:-light} ./run.sh ui

run-light: setup
	@echo "📊 Running light load test..."
	@echo "   Scenario: $${SCENARIO:-light}"
	SCENARIO=$${SCENARIO:-light} ./run.sh light

run-normal: setup
	@echo "📊 Running normal load test..."
	@echo "   Scenario: $${SCENARIO:-normal}"
	SCENARIO=$${SCENARIO:-normal} ./run.sh normal

run-stress: setup
	@echo "📊 Running stress test..."
	@echo "   Scenario: $${SCENARIO:-stress}"
	SCENARIO=$${SCENARIO:-stress} ./run.sh stress

analyze:
	@if [ -z "$(CSV_FILE)" ]; then \
		echo "Usage: make analyze CSV_FILE=results/stats.csv"; \
		echo ""; \
		echo "Available results:"; \
		ls -lh results/*.csv 2>/dev/null || echo "No test results found"; \
	else \
		python analyze.py $(CSV_FILE); \
	fi

docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t trip-manager-loadtest .

docker-run: docker-build
	@echo "🐳 Running tests in Docker..."
	docker run --rm \
		-e API_URL=${API_URL:-http://localhost:8000/api} \
		trip-manager-loadtest \
		--users 50 \
		--spawn-rate 5 \
		--run-time 5m

docker-compose:
	@echo "🐳 Running distributed tests with docker-compose..."
	@if [ ! -f ../.env ]; then \
		echo "⚠️  .env file not found in parent directory. Create it first."; \
		exit 1; \
	fi
	docker-compose up --build

clean:
	@echo "🧹 Cleaning up test results..."
	rm -rf results/*.csv results/*.html results/*.txt
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete
	@echo "✅ Cleanup completed"

.DEFAULT_GOAL := help

