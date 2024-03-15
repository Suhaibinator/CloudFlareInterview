# URL Shortener Project

This project is a URL shortener service. It takes a long URL and generates a shorter, more manageable URL. It's a flexible system designed with scalability in mind. We support two types of underlying databases: PostgreSQL and SQLite3, and it's extensible for other types of databases. All that's required to support a new database is to write a handful of functions.

The system is highly configurable. Each URL generated begins with an identifier denoting a way to identify which worker generated the URL. This is used for sharding purposes if we choose to horizontally scale the system. This design allows for high scalability from the ground up.

## Configuration

The configuration for the URL shortener is located in the `config.yml` file. This file contains settings for the URL shortener, including:

- The port on which the service listens. Note that if you modify the port here and you want metrics, you must also modify the port in the `prometheus.yml` file.
- The type of database used. You can choose between PostgreSQL and SQLite, but PostgreSQL is recommended.
- The path to where the SQLite database file will be stored on disk.
- Various settings for the PostgreSQL database, such as host, port, username, password, and database name.
- A unique identifier for this worker if this is a worker in a distributed cluster of workers for the URL shortener service. This is used for sharding and scaling.


## Running the Project

You can use the provided Makefile to run the project. It's recommended to run these commands as `sudo` for best results. Here are the available `make` commands:

- `make setup`: This command sets up the necessary directories (`prometheus_config` and `prometheus_data`), copies the `prometheus.yml` configuration file into `prometheus_config`, and changes the permissions of the directories to 777.

- `make shutdown`: This command shuts down the URL shortener service by running `docker-compose down`.

- `make test`: This command runs the tests for the URL shortener service using `python3 test.py`.

- `make cleanup`: This command cleans up the environment by deleting the `sqlite_short_urls.db` database file and the `prometheus_config`, `prometheus_data`, and `postgres_data` directories. It first runs `make shutdown` to ensure that the service is not running.

To run the project, first clean up any existing setup using the `make cleanup` command. Then, set up the project using the `make setup` command. Finally, start the service using `go run .`. To run the tests, you can use the `make test` command.

## Testing

The tests for the URL shortener service are defined in `test.py`. Here's a brief overview of how they work:

1. A list of long URLs is defined. These are the URLs that will be shortened and tested.

2. The `generate_short_url` function takes a long URL as input and sends a POST request to the URL shortener service to generate a short URL. The short URL is returned as a string.

3. The `get_data_from_url` function takes a URL as input and sends a GET request to the URL. The response data is returned as bytes.

4. A dictionary is created where the keys are the long URLs and the values are the corresponding short URLs. This is done by calling `generate_short_url` for each long URL in the list.

5. For each long URL and corresponding short URL in the dictionary, a GET request is sent to both URLs and the response data is compared. If the response data is the same for both URLs, the test passes. If the response data is different, the test fails.

To run the tests, you can use the `make test` command. This command runs the tests defined in `test.py` using Python 3.

## Metrics

Metrics are handled using Prometheus. Once the service is running, you can navigate to `localhost:9090` to access the Prometheus dashboard and compute metrics for each short URL.

You can use queries such as `requests_total{path="A000006"}` to see how many times a given URL was accessed. You can also use queries such as `sum(increase(requests_total{path="A000006"}[5m]))` to see how many times the short URL was accessed in the last 5 minutes.

Alternatively, you can use `curl` requests to retrieve the same information. Here's an example:

```bash
curl -G --data-urlencode 'query=sum(requests_total{path="A000006"})' http://localhost:9090/api/v1/query
```

## API Usage

This API provides three endpoints:

1. `POST /api/new`: This endpoint is used to create a new short URL. The request body should be a JSON object with the following properties:

    - `full_url`: The long URL that you want to shorten.
    - `expires_at`: The expiration date for the short URL in RFC3339 format (optional).

    Example usage:

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"full_url":"https://example.com", "expires_at":"2022-12-31T23:59:59Z"}' http://localhost:8080/api/new
    ```

2.  `DELETE /api/delete/{shorturl}`: This endpoint is used to delete a short URL. Replace {shorturl} with the short URL that you want to delete.

    Example usage:

    ```bash
    curl -X DELETE http://localhost:8080/api/delete/{shorturl}
    ```

3.  `GET /{shorturl}`: This endpoint is used to redirect a short URL to its corresponding long URL. Replace {shorturl} with the short URL that you want to redirect.

    Example usage:

    ```bash
    curl -X GET http://localhost:8080/{shorturl}
    ```