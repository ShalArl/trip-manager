BASE_URL ?= https://hetzner.neatnode.xyz
USERS ?= 10
SPAWN_RATE ?= 10
RUN_TIME ?= 5m

.PHONY: install seed test test-headless clean

help:
	@echo "Usage:"
	@echo "  make install           Install dependencies"
	@echo "  make seed			    Seed the database with test data"
	@echo "  make test			    Run load tests in interactive mode"
	@echo "  make test-headless	    Run load tests in headless mode"
	@echo "  make clean			    Clean up generated files"

install:
	uv venv
	uv pip install -r requirements.txt

seed:
	BASE_URL=$(BASE_URL) uv run python -m seeding.seed

provision-only:
	BASE_URL=$(BASE_URL) uv run python -m seeding.provision_only

test:
	BASE_URL=$(BASE_URL) uv run locust -f locustfile.py --host $(BASE_URL)

test-headless:
	BASE_URL=$(BASE_URL) uv run locust -f locustfile.py \
		--host $(BASE_URL) \
		--headless \
		--users $(USERS) \
		--spawn-rate $(SPAWN_RATE) \
		--run-time $(RUN_TIME)

clean:
	rm -f data/users.csv