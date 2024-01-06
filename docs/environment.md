# Environment

This document describes the required environment variables to run the micro-service.

| Name                          | Description                                     | Example                                                             | Mandatory |
| ----------------------------- | ----------------------------------------------- | ------------------------------------------------------------------- | --------- |
| `RABBIT_MQ_CONNECTION_STRING` | The connection string to the RabbitMQ instance. | `amqp://username:password@address:port/`                            | Yes       |
| `DB_CONNECTION_STRING`        | The connection string to the postgres database  | `postgres://username:password@domain:port/database?sslmode=disable` | Yes       |
