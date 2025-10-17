![Greenbone Logo](https://www.greenbone.net/wp-content/uploads/gb_new-logo_horizontal_rgb_small.png)

# OpenSight Notification Service <!-- omit in toc -->

[![codecov](https://codecov.io/github/greenbone/opensight-notification-service/graph/badge.svg?token=BE6ExbOYqE)](https://codecov.io/github/greenbone/opensight-notification-service)

The Notification Service allows to display all notifications in a central place in the frontend. All OpenSight backend services can use it to send notifications to the user.

## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
- [Requirements](#requirements)
- [Configuration](#configuration)
- [Running](#running)
  - [Running non-containerized service](#running-non-containerized-service)
- [Build](#build)
- [Testing](#testing)
- [Maintainer](#maintainer)
- [License](#license)

## Usage

The Notification Service is intended to be deployed along the others openSight services on the appliance. The service provides a REST API. See the [OpenApi definition](api/notificationservice/notificationservice_swagger.yaml) for details. You can view them e.g. in the [Swagger Editor](https://editor.swagger.io/).

Backend services can send notifications via the `Create Notification` endpoint. Those notifications can then be retrieved via `List Notifications` to provide them to the user.

## Requirements

To run the service and its dependencies on your local machine you need a working installation of [docker](https://docs.docker.com/engine/install/) and `make`.

For running the Notification Service outside of docker the latest version of [go](https://go.dev/doc/install) must be installed.

## Configuration

The service is configured via environment variables. Refer to the [Config](pkg/config/config.go) for the available options and their defaults.

The secret `DB_PASSWORD` can be also passed by file. Simply pass the path to the file containing the secret to the service by appending `_FILE` to the env var name, i.e. `DB_PASSWORD_FILE`. If the secret is supplied in both ways, the one directly passed by env var takes precedence.

## Running

> The following instructions are targeted at openSight developers. As end user the services should be run in orchestration with the other openSight services, which is not in the scope of this readme.

Before starting the services you need to set the environment variable `DB_PASSWORD` for the database password. This password is set (only) on the first start of the database, so make sure to use the same password on consecutive runs. All other configuration parameters are already set in the docker compose files. 

Then you can start the notification service and the required dependent service [Postgres](https://www.postgresql.org/) with

```sh
# make sure the `DB_PASSWORD` environment variable is set
make start-services
```

The port of the notification service is forwarded, so you can access its API directly from the host machine at the base url http://localhost:8085/api/notification-service. A convenient way to interact with the service are the Swagger docs served by the running service at URL http://localhost:8085/docs/notification-service/notification-service/index.html.

### Running non-containerized service

If you are actively developing it might be more convenient to run the notification service directly on the host machine. 

First start the containerized database with

```sh
docker compose up -d
```

Then load the configuration parameters into environment variables and start the notification service:

```sh
# make sure the `DB_PASSWORD` environment variable is set beforehand
source set_env_vars_local_setup.sh && go run ./cmd/notification-service
```
From here everything is identical to the docker compose setup.

If there are errors regarding the database, verify that it is running with `docker ps` (should show a running container for postgres).

## Build

> Refer to [Makefile](./Makefile) to get an overview of all commands

To build run `make build`. The rest API docs can be generated with `make api-docs`.

## Testing

Currently implemented test suites are displayed in table below. Not all Go tests are enabled by default, see `Env filter` column for specific environment variable needed to be set to run the tests - e.g. to include Postgres Unit Tests in go tests execution, run:

```shell
TEST_POSTGRES=1 go test ./...
```

| Title                 | Description                                           | Command                   | Env filter           |
|-----------------------|-------------------------------------------------------|---------------------------|----------------------|
| Unit Tests            | short unit tests, no additional setup needed          | `make test`               | (run by default)     |
| Postgres Unit Tests   | unit tests running on containerized Postgres instance | `make test-postgres`      | `TEST_POSTGRES=1`    |

## Maintainer

This project is maintained by [Greenbone AG][Greenbone]

## License

Copyright (C) 2024 [Greenbone AG][Greenbone]

Licensed under the [GNU Affero General Public License v3.0 or later](LICENSE).

[Greenbone]: https://www.greenbone.net/
