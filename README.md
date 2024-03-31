# ad-service-api

This project is a Go-based advertisement service API. It provides endpoints for creating and listing advertisements by given query parameters.

## Project Structure

The project is organized into several directories:

- [`database/`](command:_github.copilot.openRelativePath?%5B%22database%2F%22%5D "database/"): Contains the MongoDB connection setup.
- [`docs/`](command:_github.copilot.openRelativePath?%5B%22docs%2F%22%5D "docs/"): Contains the Swagger documentation for the API description.
- [`internal/`](command:_github.copilot.openRelativePath?%5B%22internal%2F%22%5D "internal/"): Contains the core business logic of the application.
  - `advertisement/`: Contains the handlers, repositories, and services for the advertisement functionality.
  - `middleware/`: Contains middleware functions.
  - `models/`: Contains the data models used in the application.
- `router/`: Contains the router setup for the API.
- `validators/`: Contains validation logic for the API inputs.
- [`redis/`](command:_github.copilot.openRelativePath?%5B%22redis%2F%22%5D "redis/"): Contains the Redis connection setup.
- [`utils/`](command:_github.copilot.openRelativePath?%5B%22utils%2F%22%5D "utils/"): Contains utility functions used across the application.

## Setup and Running

To set up the project, please navigate to your root directory to execute the `docker-compose.yml` in the terminal by running below cmd:

```sh
docker-compose up
```

- **--build**: If you have changed the Dockerfile, you should add this flag after above command to ensure the setup of latest image
- **-d**: If you want to make docker container run in background, you should add this flag.

This will start the server on your local machine by using docker.

## API Endpoints

The API provides the following endpoints:

- `POST /api/v1/ad`: Creates a new advertisement. The request body should be a JSON object that matches the `models.Advertisement` structure.
- `GET /api/v1/ads`: Lists all advertisements which match the query parameters if they exist.

## Testing

Unit tests is to ensure tha basic logic for each functions without depending on the external service. They are located in the same directories as the files they are testing tailed with . For example, tests for the advertisement handler are located in [`internal/advertisement/handler/advertisement_handler_test.go`](command:_github.copilot.openRelativePath?%5B%22internal%2Fadvertisement%2Fhandler%2Fadvertisement_handler_test.go%22%5D "internal/advertisement/handler/advertisement_handler_test.go").

To run all tests in the project, you can use the following command in your terminal. This command will recursively run all tests in the project. :

```sh
go test ./...
```

## Documentation

API documentation is provided with Swagger and can be accessed at `/swagger/*any` when the server is running. If you wanna see it clearly, you may direct to
`http://127.0.0.1:8080/swagger/index.html`

## Dependencies

The project uses several dependencies, mainly including:

- `github.com/gin-gonic/gin`: A HTTP web framework used for routing and handling HTTP requests.
- `github.com/swaggo/files`, `github.com/swaggo/gin-swagger`, `github.com/swaggo/swag`: Used for generating and serving Swagger API documentation.
- `github.com/pariz/gountries`: A library for country data.

For a full list of dependencies, refer to the [`go.mod`](command:_github.copilot.openRelativePath?%5B%22go.mod%22%5D "go.mod") file.

## Contributing

Contributions are welcome. Please make sure to update tests as appropriate when making changes.

## License

This project is licensed under the terms of the MIT license.