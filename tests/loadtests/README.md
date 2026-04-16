# Loadtesting Module

Load testing is performed with Locust, a modern load testing framework that allows you to define user behavior in Python code and simulate millions of users. The load testing module is designed to help you evaluate the performance and scalability of your application under various conditions.
For data generation faker is used to create more realistic user data.

## Setup 

To begin load testing create a new .env parallel to the .env.example file you'll find in this directory.
```text
BASE_URL=http://localhost:8000/api
SEED_COUNT=5000
SEED_CONCURRENCY=50
USERS_CSV=data/users.csv
```

install the required dependencies using UV:

```bash
uv sync
```

## Seeding (local testing only)
To create a large dataset for load testing, you can use the `seed` target in the Makefile. This will generate a specified number of users and save them to a CSV file. You can adjust the `SEED_COUNT` and `SEED_CONCURRENCY` values in your .env file to control how many users are created and how many concurrent requests are made during seeding.
```bash
make seed
```

without any adaptations this will create 5000 users in the USERS_CSV file. The file is subsequently used for load testing.

## Load Testing

To run the load tests, use the `test` target in the Makefile. This will execute the Locust load tests using the user data from the CSV file created during seeding.
```bash
make test
```

### Load Testing in Headless Mode (CI/CD)
For running load tests in a CI/CD pipeline, you can use the `test-headless` target in the Makefile. This will execute the load tests in headless mode, which is suitable
for automated testing environments.
```bash
make test-headless
```

## Clean Up
To clean up the generated user data and any other artifacts from the load testing process, you can use the `clean` target in the Makefile. This will remove the CSV file containing the user data.
```bash
make clean
```

## Further Information

For more information checkout the following resources:
- [Locust](https://locust.io/) - A modern load testing framework that allows you to define user behavior in Python code and simulate millions of users.
- [Makefile](https://www.gnu.org/software/make/) - A build automation tool that helps to manage and automate the execution of tasks, such as seeding data and running load tests.
- [dotenv](https://pypi.org/project/python-dotenv/) - A Python library that allows you to read key-value pairs from a .env file and set them as environment variables, making it easier to manage configuration settings for your load tests.
- [CSV](https://docs.python.org/3/library/csv.html) - A Python library for reading and writing CSV files, which is used to manage the user data for load testing.
- [UV](https://docs.astral.sh/uv/) - A Python package manager that can be used to manage dependencies and virtual environments for your load testing project.
- [pyproject.toml](https://www.python.org/dev/peps/pep-0561/) - A configuration file used to specify project metadata and dependencies, which can be used in conjunction with UV for managing your load testing project.
- [Faker](https://faker.readthedocs.io/en/master/) - A Python library that generates fake data, which can be used to create realistic user data for load testing.
- [Pillow](https://pillow.readthedocs.io/en/stable/) - A Python Imaging Library that adds image processing capabilities to your Python interpreter, which can be used to generate images for load testing scenarios that involve media uploads.