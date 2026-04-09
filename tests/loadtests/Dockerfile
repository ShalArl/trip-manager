FROM python:3.10-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy locust files
COPY locustfile.py .

# Default command
ENTRYPOINT ["locust"]
CMD ["-f", "locustfile.py", "--headless"]

