# General organization contribution guidelines

Please, read the [General Contribution Guidelines](https://github.com/upb-code-labs/docs/blob/main/CONTRIBUTING.md) before contributing.

# Submissions Status Updater micro-service contribution guidelines

## Project structure / architecture

The Submissions Status Updater micro-service's architecture is based on the Hexagonal Architecture in order to make it easier to maintain and extend.

The project structure is as follows:

- `domain`: Contains the business logic of the Submissions Status Updater micro-service.

  - `definitions`: Contains the interfaces (contracts) of the Submissions Status Updater micro-service.
  - `entities`: Contains the entities of the Submissions Status Updater micro-service.
  - `dtos`: Contains the data transfer objects of the Submissions Status Updater micro-service.

- `application`: Contains the application logic (use cases) of the Submissions Status Updater micro-service.

- `infrastructure`: Contains the implementation of the Submissions Status Updater micro-service's interfaces (contracts) and the external dependencies of the Submissions Status Updater micro-service.

  - `implementations`: Contains the implementations of the `domain` interfaces (contracts).

- `shared`: Contains the shared code of the Submissions Status Updater micro-service.

  - `utils`: Contains the utility functions of the Submissions Status Updater micro-service.

## Local development

### Dependencies

The following dependencies are required to run the Submissions Status Updater micro-service locally:

- [Go 1.21.5](https://golang.org/doc/install)
- [Podman](https://podman.io/getting-started/installation) (To build and test the container image)

Please, note that `Podman` is a drop-in replacement for `Docker`, so you can use `Docker` instead if you prefer.

Additionally, you may want to install the following dependencies to make your life easier:

- [Air](https://github.com/cosmtrek/air) (for live reloading)

### Running the Submissions Status Updater micro-service locally

As the role of the Submissions Status Updater micro-service is to listen for messages in the `submission-status-updates` queue and update the status of the submissions in the database and publish the updated submissions to the `real-time-updates` queue, you will need to run the [gateway](https://github.com/UPB-Code-Labs/main-api) project first in order to initialize the queue, database, RabbitMQ and the other micro-services and then you can send submissions by using the REST API.

After you have the gateway running, you can start the Submissions Status Updater micro-service by running the following command:

```bash
air
```

This will start the Submissions Status Updater micro-service and will watch for changes in the source code and restart the service automatically.
