[![build](https://github.com/Yoann-Abbes/kafka-message-scheduler-admin/actions/workflows/build.yml/badge.svg)](https://github.com/Yoann-Abbes/kafka-message-scheduler-admin/actions/workflows/build.yml)
[![docker](https://github.com/Yoann-Abbes/kafka-message-scheduler-admin/actions/workflows/docker.yml/badge.svg)](https://github.com/Yoann-Abbes/kafka-message-scheduler-admin/actions/workflows/docker.yml)
**Scheduler admin** is a GUI for managing a list of [kafka message schedulers](https://github.com/etf1/kafka-message-scheduler)
### User Interface
![Home](docs/screenshots/one.png)
![List](docs/screenshots/two.png)
![Detail](docs/screenshots/three.png)
## Getting started
To run the scheduler admin you can use docker, it will need a scheduler to connect to, you can specify in the variable env. `SCHEDULERS_ADDR` for example: `SCHEDULERS_ADDR=scheduler`.
### Regular version
```
docker run -d -p 9000:9000 -e SCHEDULERS_ADDR=<schedulers-address> yabbes/kafka-message-scheduler-admin
```
### Mini version
The mini version is a "mocked" version of the admin all in one, for demonstration purpose
```
docker run -d -p 9000:9000 yabbes/kafka-message-scheduler-admin:mini
```
Then open browser at `localhost:9000`
## Usage
The server exposes two ports:
- `9000` is the server port.
- `9001` is the port for exposing prometheus metrics.
- `/` will expose the user interface
- `/api` will expose the api endpoints
- `/api/health` lightweight health check
For the full API reference, see [API.md](API.md).
## Configuration
| Env. variable    | Default         | Description                                                                                                                                                |
|------------------|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| LOG_LEVEL        | info            | logging level (panic, fatal, error, warning, info, debug, trace)                                                                                           |
| GRAYLOG_SERVER   |                 | graylog server address                                                                                                                                     |
| METRICS_ADDR     | :9001           | prometheus metrics port                                                                                                                                    |
| SERVER_ADDR      | :9000           | server address port                                                                                                                                        |
| SCHEDULERS_ADDR  | localhost:8000  | comma separated list of address of schedulers, may or may not contain port (default port is 8000), for example: SCHEDULERS_ADDR=scheduler1,scheduler2:8000 |
| STATIC_FILES_DIR | ../client/build | location of the UI static files for the HTML & js files                                                                                                    |
| DATA_ROOT_DIR    | ./.db           | Default location of internal database files                                                                                                                |
| API_SERVER_ONLY  | false           | when true, only the rest api is exposed without serving the static files and default route is / (instead of /api)                                          |
| KAFKA_MESSAGE_BODY_DECODER  |            | set an endpoint for decoding kafka message payload. Post with payload {id:xxx target-topic:yyy value:[base64 of the kafka message body]}                   |
| PPROF_ENABLED    | false           | when true, exposes Go pprof endpoints on the metrics server                                                                                                |
## Development
### Prerequisites
- **Go** ≥ 1.24
- **Node.js** ≥ 22
- **pnpm** (install: `corepack enable && corepack prepare pnpm@latest --activate`)
### Backend (in folder ./server)
For development you will need external dependencies such as kafka, in order to start this dependencies you can run:
- `make up`: startup development environment with external dependencies (kafka, scheduler, ...)
- `make down`:  shutdown development environment
The backend is written in Go. You can use the following commands to manage the development lifecycle:
- `make start`: start GO server
- `make build`: compile the code
- `make bin`: generate a binary
- `make lint`: run static analysis on the code
- `make test`: execute unit tests
- `make test.integration`: execute integration tests
- `make tests`: execute all tests
- `make tests.docker`: execute all tests in "black box" inside docker containers
#### Quick start
```
cd server
make up
make start
```
### Frontend (in folder ./client)
For development you will need a running admin server launched as described before or the mini version of the scheduler admin which is running without any external dependencies:
- `docker run -p 9000:9000 yabbes/kafka-message-scheduler-admin:mini`: startup mini version of the scheduler (no external dependencies required)
or start a standard admin server
- `make start` (inside /server): startup GO server on local
The frontend is written with TypeScript and React. You can use the following commands to manage the development lifecycle:
- `pnpm install`: install the dependencies
- `pnpm start`: start the frontend in development mode with Vite HMR
- `pnpm run build`: generate the transpiled and minified files and assets
- `pnpm test`: execute unit tests
To start the dev server: `pnpm install && pnpm start`
Then open browser at: http://localhost:3000
#### Quick start
```
cd client
docker run -d -p 9000:9000 yabbes/kafka-message-scheduler-admin:mini
pnpm install
pnpm start
```
## Improvements over original
This fork includes the following improvements over the original repository:
### Security
- Non-root Docker user
- `ReadHeaderTimeout` on all HTTP servers (CVE mitigation)
- pprof gated behind `PPROF_ENABLED` env var
- CORS restricted to same-origin (configurable)
- `ca-certificates` in Docker image
### Code Quality
- Go 1.24 (upgraded from 1.15)
- All dependencies upgraded to latest stable
- React 18 + React Router v6 (upgraded from React 17 + v5)
- Migrated from CRA to Vite (faster builds, smaller bundles)
- Migrated from yarn to pnpm
- TypeScript 5 (upgraded from 4)
- `highlight.js` direct usage (replaced unmaintained `react-highlight`)
- Fixed nil map panic in `hmap.NewSchedulerSchedules()`
- Fixed `time.Until(start)` → `time.Since(start)` timing bug
- Fixed `DecodeJSON` ignoring timeout parameter
- Fixed `getBool()` to use `strconv.ParseBool`
- Fixed `KafkaMessageBodyDecoder()` empty string bug
- Fixed flaky mini tests (OS-assigned ports)
- Removed `fmt.Printf` debug statements (use structured logging)
- Added `/health` endpoint
- Added unit tests for config, helper, sort packages
### CI/CD
- Modern GitHub Actions (checkout@v4, setup-go@v5, setup-node@v4)
- Separate Go and Client CI jobs
- Docker Hub ready with BuildX caching
## Contributors
- [Fatih KARAKAŞ](https://github.com/fkarakas)
- [Emmanuel FERTE](https://github.com/eferte)
