# ad-service-api

This project is a Go-based advertisement service API. It provides endpoints for creating and listing advertisements by given query parameters.

*Noticed: This readme is specific for MacOS, if you're using different OS, please find the corresponding cmd for your needs.*

## Project Structure

- [`database/`](database/): Contains the MongoDB connection setup.
- [`docs/`](docs/): Contains the Swagger documentation for the API description.
- [`internal/`](internal/): Contains the core business logic of the application.
  - `advertisement/`: Contains the handlers, repositories, and services for the advertisement functionality.
  - `middleware/`: Contains middleware functions.
  - `models/`: Contains the data models used in the application.
- `router/`: Contains the router setup for the API.
- `validators/`: Contains validation logic for the API inputs.
- [`redis/`](redis/): Contains the Redis connection setup.
- [`utils/`](utils/): Contains utility functions used across the application.

## Setup and Running

### Using Docker in local environment
To set up the project, please navigate to your root directory to execute the `docker-compose.yml` in the terminal by running below cmd:

```sh
# flag -d make docker container run in background.
docker-compose up -d
```

- **--build**: If you have changed the Dockerfile, you should add this flag after above command to ensure the setup of latest image. After succesfully build new image, ensure to delete the old image in you docker

This will start the server on your local machine by using docker.

### Using Helm Chart and Minikube in local environment

This part is ready for the autoscaling and load balancing to handle substantial requests.

1. Please ensure that you've installed the minikube and set the config succesfully - [minikube start](https://minikube.sigs.k8s.io/docs/start/)

2. For docker repo which is public, you don't need imagePullSectets. Instead, if it is private, follow the instruction 3. and 4.

3. (Optional) Create k8s secret to access private docker repo by contacting admin to get these credentials
    ```sh
    kubectl create secret docker-registry my-registry-secret \
    --docker-username=DOCKER_USERNAME \
    --docker-password=DOCKER_PASSWORD \
    --docker-email=DOCKER_EMAIL
    ```

4. (Optional) Add `my-registry-secret` at imagePullSecrets within `values.yaml`

5. Modify your local DNS setting and store it
    ```sh
    # access the file with vim mode
    sudo vim /etc/hosts

    # add below setting into the file
    127.0.0.1 ad-service-api.local
    ```

6. Add secrets for MongoDb and Redis
    ```sh
    # mongo-secret
    kubectl create secret generic mongodb-secret \
    --from-literal=mongodb-root-password=YOUR_ROOT_PASSWORD \
    --from-literal=username=YOUR_USERNAME \
    --from-literal=password=YOUR_PASSWORD_FOR_USERNAME \
    --from-literal=host=<K8S_SERVICE_NAME_FOR_MONGO>:<MONGO_PORT> \
    --from-literal=database=YOUR_DATABASE_NAME

    # redis-secret
    kubectl create secret generic redis-secret \
    --from-literal=redis-password=YOUR_REDIS_PASSWORD \
    --from-literal=host=<K8S_SERVICE_NAME_FOR_REDIS>:<REDIS_PORT> \
    --from-literal=database=0
    ```

7. change to new docker image for `ad-service-api` tag in `values.yaml` by refering to deploy stage in Github Action

8. Ready to release the helm chart and build resources
    ```sh
    # move to helm directory
    cd helm

    # if this is your first time to release helm chart
    helm install ad-service-api ./ad-service-api

    # if you'd released before, and you've changed helm chart config
    helm upgrade ad-service-api ./ad-service-api
    ```

9. Use below cmd to build a connection tunnel from localhost to minikube
    ```sh
    minikube tunnel
    ```

## API Endpoints

The API provides the following endpoints:

- `POST /api/v1/ad`: Creates a new advertisement. The request body should be a JSON object that matches the `models.Advertisement` structure.
- `GET /api/v1/ads`: Lists all advertisements which match the query parameters if they exist. Below is the params list:
  - age (can be empty)
  - country (can be empty)
  - platform (can be empty)
  - limit (default to 5)
  - offset (default to 0)

## Testing

Unit tests is to ensure tha basic logic for each functions without depending on the external service. They are located in the same directories as the files they are testing tailed with . For example, tests for the advertisement handler are located in [`internal/advertisement/handler/advertisement_handler_test.go`](internal/advertisement/handler/advertisement_handler_test.go).

To run all tests in the project at local, you can use the following command in your terminal. This command will recursively run all tests in the project. :

```sh
go test ./...
```

## CI/CD

For this project, we use Github Action to automatically set up a pipeline when you push the code to main branch, and the workflow would follow init, test, build image and deploy stages

Noticed: if you fork this project to your repo, and use the docker repo which is private, remember to set the Github secrets for docker image pull secrets in your Github repository secrets.

1. For docker username use this name: DOCKERHUB_USERNAME
2. For docker password use this name: DOCKERHUB_PASSWORD

Make sure you've set this before you run the workflow in Github Action.

Refer to [set-up-repo-secrets-in-github](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions#creating-secrets-for-a-repository)

## OpenAPI Documentation

API documentation is provided with Swagger and can be accessed at `/swagger/*any` when the server is running. If you wanna see it clearly, you may direct to
`http://127.0.0.1:8080/swagger/index.html`

## Dependencies

For a full list of dependencies, refer to the [`go.mod`](command:_github.copilot.openRelativePath?%5B%22go.mod%22%5D "go.mod") file.

## Contributing

Contributions are welcome. Please make sure to update tests as appropriate when making changes.

## License

This project is licensed under the terms of the MIT license.