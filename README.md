[![build](https://github.com/etf1/kafka-message-scheduler-admin/actions/workflows/build.yml/badge.svg)](https://github.com/etf1/kafka-message-scheduler-admin/actions/workflows/build.yml)
[![docker](https://github.com/etf1/kafka-message-scheduler-admin/actions/workflows/docker.yml/badge.svg)](https://github.com/etf1/kafka-message-scheduler-admin/actions/workflows/docker.yml)

**Scheduler admin** is a GUI for managing a list of [kafka message schedulers](https://github.com/etf1/kafka-message-scheduler)


### User Interface

![Home](docs/screenshots/one.png)
![List](docs/screenshots/two.png)
![Detail](docs/screenshots/three.png)

## Getting started

To run the scheduler admin you can use docker, it will need a scheduler to connect to, you can specify in the variable env. `SCHEDULERS_ADDR` for example: `SCHEDULERS_ADDR=scheduler`.

### Regular version

```
docker run -d -p 9000:9000 -e SCHEDULERS_ADDR=<schedulers-address> etf1/kafka-message-scheduler-admin
```

### Mini version

The mini version is a "mocked" version of the admin all in one, for demonstration purpose

```
docker run -d -p 9000:9000 etf1/kafka-message-scheduler-admin:mini
```

Then open browser at `localhost:9000`

## Usage

The server exposes two ports:

- `9000` is the server port. 
- `9001` is the port for exposing prometheus metrics.

- `/` will expose the user interface
- `/api` will expose the api endpoints

## API Routes

GET methods

URL Parameters:
- `{name}`: scheduler name
- `{id}`: schedule ID

### config
- `/stats` : expose some statistics
- `/schedulers` : list of registered schedulers

### all schedules
- `/scheduler/{name}/schedules`: search for schedules 
- `/scheduler/{name}/schedule/{id}`: get schedule detail

### live schedules
- `/live/scheduler/{name}/schedules`: search for schedules
- `/scheduler/{name}/schedule/{id}`: get schedule detail

### history schedules
- `/history/scheduler/{name}/schedules`: search for schedules
- `/history/scheduler/{name}/schedule/{id}`: get schedule detail

### search parameters

- `schedule-id`: part of the schedule ID
- `epoch-from`: lower range of schedule epoch
- `epoch-to`: upper range of schedule epoch
- `max`: max number of result returned (cannot be more than 1000)
- `sort-by`: sort field, format is `field order`. 
   - Available options for field are: `timestamp`, `id`, `epoch`
   - Available options for order are: `asc`, `desc`
   - Default is `timestamp desc`

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
| KAFKA_MESSAGE_BODY_DECODER  |            | set an endpoint for decoding kafka message payload. Post with payload {id:xxx target-topic:yyy value:[base64 of the kafka message body]}                                          |

## Development

### Backend (in folder ./server)

#### Prerequisities

For development you will need external dependencies such as kafka, in order to start this dependencies you can run :

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

Then start the server `make start`

#### Quick start

```
cd server
make up
make start
```

### Frontend (in folder ./client)

#### Prerequisities

For development you will need a running admin server launched as described before or the mini version of the scheduler admin which is running without any external dependencies:

- `docker run -p 9000:9000 etf1/kafka-message-scheduler-admin:mini`: startup mini version of the scheduler (no external dependencies required)

or start a standard admin server

- `make start` (inside /server): startup GO server on local

The frontend is written with TypeScript and React. You can use the following commands to manage the development lifecycle:

- `yarn`: install the dependencies
- `yarn start`: start the frontend in development mode, with live reload
- `yarn build`: generate the transpiled and minified files and assets
- `yarn test`: execute unit tests

To start the nodejs dev. server `yarn && yarn start`

Then open browser at: http://localhost:3000

#### Quick start

```
cd client
docker run -d -p 9000:9000 etf1/kafka-message-scheduler-admin:mini
yarn
yarn start
```

## Contributors

- [Fatih KARAKAŞ](https://github.com/fkarakas)
- [Emmanuel FERTE](https://github.com/eferte)
