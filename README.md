
# File Processing Service

This service is designed to process CSV files using Kafka and Redis, built using Go 1.20. The service reads messages from a Kafka topic, processes each row in the CSV file, and performs operations based on the data.

## Prerequisites

- Go 1.20 or higher
- Kafka 3.6.1
- Redis server 7.2.1

## Configuration

Before running the service, set the necessary constants in `constants.go`. Pay special attention to `apiURL` and rate limit constants. Other constants have default values but can be adjusted as needed. The number of workers in the worker pool can ideally be equal to the rate limit in seconds.

## Build Instructions

1. Clone the repository:

   ```
   git clone [repository-url]
   ```

2. Navigate to the project directory:

   ```
   cd [project-directory]
   ```

3. Build the project:

   ```
   go build
   ```

## Run Instructions

1. Ensure Kafka is running locally with default settings.
2. Start the Redis server.
3. Run the built executable:

   ```
   ./[executable-name]
   ```

## Simulation

To simulate the complete flow, you can expose a simple POST endpoint locally that accepts a JSON request body. Point to this endpoint in the `constants.go` configuration file for testing purposes.

## Features

- Consumes messages from a Kafka topic.
- Processes CSV files specified in the Kafka messages.
- Uses a worker pool for concurrent processing.
- Interfaces with Redis to track processed data.
- Rate limits API requests to avoid exceeding quotas.
- As it uses Kafka consumer group, it is scalable and multiple instances can be brought up in case of queue build up.
- Logs informative and error messages for troubleshooting.

## Notes

- Ensure all dependencies are correctly installed.
- Check Kafka and Redis connection settings.
- Review the rate limit settings to match your use case.

---

## TODO

Refactor code into packages. Inlcude unit tests.
