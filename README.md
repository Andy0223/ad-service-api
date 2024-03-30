# ad-service-api

This project is a Go-based advertisement service API. It provides endpoints for creating and listing advertisements by given query parameters.

## Project Structure

The project is organized into several directories:

- [`database/`](command:_github.copilot.openRelativePath?%5B%22database%2F%22%5D "database/"): Contains the MongoDB connection setup.
- [`docs/`](command:_github.copilot.openRelativePath?%5B%22docs%2F%22%5D "docs/"): Contains the Swagger documentation for the API.
- [`internal/`](command:_github.copilot.openRelativePath?%5B%22internal%2F%22%5D "internal/"): Contains the core business logic of the application.
  - `advertisement/`: Contains the handlers, repositories, and services for the advertisement functionality.
  - `middleware/`: Contains middleware functions.
  - `models/`: Contains the data models used in the application.
- `router/`: Contains the router setup for the API.
- `validators/`: Contains validation logic for the API inputs.
- [`redis/`](command:_github.copilot.openRelativePath?%5B%22redis%2F%22%5D "redis/"): Contains the Redis connection setup.
- [`utils/`](command:_github.copilot.openRelativePath?%5B%22utils%2F%22%5D "utils/"): Contains utility functions used across the application.

## Setup

To set up the project, you need to have Go installed on your machine. Then, you can clone the repository and install the dependencies listed in the [`go.mod`](command:_github.copilot.openRelativePath?%5B%22go.mod%22%5D "go.mod") file.

## Running the Project

To run the project, navigate to the project directory and run the [`main.go`](command:_github.copilot.openRelativePath?%5B%22main.go%22%5D "main.go") file:

```sh
go run main.go
```

This will start the server on your local machine.

## API Endpoints

The API provides the following endpoints:

- `POST /api/v1/ad`: Creates a new advertisement. The request body should be a JSON object that matches the `models.Advertisement` structure.
- `GET /api/v1/ads`: Lists all advertisements.

## Testing

Unit tests are located in the same directories as the files they are testing. For example, tests for the advertisement handler are located in [`internal/advertisement/handler/advertisement_handler_test.go`](command:_github.copilot.openRelativePath?%5B%22internal%2Fadvertisement%2Fhandler%2Fadvertisement_handler_test.go%22%5D "internal/advertisement/handler/advertisement_handler_test.go").

## Documentation

API documentation is provided with Swagger and can be accessed at `/swagger/*any` when the server is running.

## Dependencies

The project uses several dependencies, including:

- `github.com/gin-gonic/gin`: A HTTP web framework used for routing and handling HTTP requests.
- `github.com/swaggo/files`, `github.com/swaggo/gin-swagger`, `github.com/swaggo/swag`: Used for generating and serving Swagger API documentation.
- `github.com/pariz/gountries`: A library for country data.

For a full list of dependencies, refer to the [`go.mod`](command:_github.copilot.openRelativePath?%5B%22go.mod%22%5D "go.mod") file.

## Contributing

Contributions are welcome. Please make sure to update tests as appropriate when making changes.

## License

This project is licensed under the terms of the MIT license.