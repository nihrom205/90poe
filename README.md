# Port API

- Given a file with ports data (ports.json), write a port domain service that either creates a new record in a database, or updates the existing one (Hint: no need for delete or other methods).
- The file is of unknown size, it can contain several millions of records, you will not be able to read the entire file at once.
- The service has limited resources available (e.g. 200MB ram).
- The end result should be a storage containing the ports, representing the latest version found in the JSON. (Hint: use an in memory database to save time and avoid complexity).
- A Dockerfile should be used to contain and run the service (Hint: extra points for avoiding code building in docker).
- Provide at least one example per test type that you think are needed for your assignment. This will allow the reviewer to evaluate your critical thinking as well as your knowledge about testing.
- Your readme.md should explain how to run your program and test it.
- The service should handle certain signals correctly (e.g. a TERM or KILL signal should result in a graceful shutdown).

## Bonus points

- Address security concerns for Docker
- Code structure according to hexagonal architecture

## How to run
- `make dc` runs docker-compose with the app container on port 8080 for you.
- `make test` runs the tests
- `make run` runs the app locally on port 8080 without docker.
- `make lint` runs the linter
