# ad-service-api

This project is a Go-based advertisement service API. It provides endpoints for creating and listing advertisements by given query parameters.

*Noticed: This readme is specific for MacOS, if you're using different OS, please find the corresponding cmd for your needs.*

# Design Logic

1. **MongoDB:** [advertisement structure](/internal/models/advertisemesnt.go)
    - **Simple:** for small project quick setup
    - **Speed:** MongoDB can provide fast access to data due to its ability to handle large amounts of unstructured data, which can be beneficial for an advertisement service where speed is crucial for a good user experience.
    - **Scalability:** MongoDB is designed to be horizontally scalable, which can be beneficial for a service that might need to handle a large volume of data and traffic.

2. **Redis:** Store advertisements which is frequently queried or only for temporary need. Redis provide faster access than mongodb
    - **DailyAdCreatedCounts:** store the ads created today
    - **Advertisements list with specific query params:**
        - if a new advertisement is inserted to database, the key will be removed from redis
        - if the one of the ad from redis is expired, it would directly retrieve the new data from database, and then overwrite a new value with existing key

3. **Layered Architecture:**

    - **Router**: The router is responsible for directing incoming HTTP requests to the appropriate handlers based on the request method (GET, POST, etc.) and the URL. This separation of concerns makes the code easier to maintain and understand.

    - **Handler**: Handlers are functions that execute in response to a particular route being hit. They contain the logic to process the request and send a response. Separating this logic into handlers keeps the code organized and makes it easier to manage.

    - **Service**: The service layer is where the business logic of your application resides. It interacts with the repository to fetch, manipulate, and store data. Keeping business logic in the service layer keeps it separate from the data access logic, making the code more maintainable and testable.

    - **Repository**: The repository is responsible for data storage and retrieval. It interacts with the database or other data sources. By keeping data access code in the repository, you can change the underlying data source without affecting the rest of your code.

## Project Structure

- [`database/`](database/): Contains the MongoDB related functionality.
    - [`mongo-init/`](/database/mongo-init/): initialize new user and database

- [`docs/`](docs/): Contains the Swagger documentation for the API description.

- [`helm/`](helm/): Contains Helm chart files for deploying the application on Kubernetes.
    - [`ad-service-api/`](helm/ad-service-api/): Contains the resources config for the Ad Service API.

- [`internal/`](internal/): Contains the core business logic of the application.
    - `advertisement/`: Contains the handlers, repositories, and services for the advertisement functionality.
    - `middleware/`: Contains middleware functions. ex: logger
    - `models/`: Contains the data models used in the application.

- [`router/`](internal/router/): Contains the router setup for the API.

- [`validators/`](internal/validators/): Contains validation logic for inputs from API endpoints.

- [`redis/`](redis/): Contains the Redis connection setup.

- [`mock/`](mocks/): Contains the mock functions for service and repository

## Setup and Running

### Using Docker in local environment (for development)
To set up the project, please navigate to your root directory to execute the `docker-compose.yml` in the terminal by running below cmd:

```sh
# move to root dir
cd ad-service-api

# flag -d make docker container run in background.
docker-compose up -d
```

- **--build**: If you have changed the Dockerfile, you should add this flag after above command to ensure the setup of latest image. After succesfully build new image, ensure to delete the old image in you docker.

This will start the server on your local machine by using docker.

### Using Helm Chart and Minikube in local environment

This design is ready for the autoscaling and load balancing to handle substantial requests.

1. Please ensure that you've installed the minikube and set the config succesfully - [minikube start](https://minikube.sigs.k8s.io/docs/start/)

    ```sh
    # set alias
    alias kubectl="minikube kubectl --"

    # start minikube container
    minikube start

    # enable ingress
    minikube addons enable ingress
    ```

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
    --from-literal=host=<CHART_NAME-mongo>:<MONGO_PORT> \
    --from-literal=database=YOUR_DATABASE_NAME

    # redis-secret
    kubectl create secret generic redis-secret \
    --from-literal=redis-password=YOUR_REDIS_PASSWORD \
    --from-literal=host=<CHAERT_NAME-redis-master>:<REDIS_PORT> \
    --from-literal=database=0
    ```

7. change to new docker image tag for `ad-service-api` in `values.yaml` by refering to deploy stage within Github Action (e.g. build-01)

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

10. Start to send API requests
    ```sh
    # Create Ad via POST /api/v1/ad
    curl -X POST -H "Content-Type: application/json" \
    "http://ad-service-api.local/api/v1/ad" \
    --data '{
        "title": "AD 01",
        "startAt": "2024-03-29T03:00:00.000Z",
        "endAt": "2024-04-08T23:00:00.000Z",
        "conditions": {
            "ageStart":20,
            "ageEnd":100,
            "gender": ["M", "F"],
            "country": ["TW", "US"],
            "platform": ["android", "web"]
        }
    }'

    # list ads via GET /api/v1/ad
    curl -X GET -H "Content-Type: application/json" \
    "http://ad-service-api.local/api/v1/ad?country=TW"
    ```
## API Endpoints

The API provides the following endpoints:

1. if you build with docker-compose at local environment, your api host is `localhost:8080`
2. if you deploy to minikube via helm chart, your api host is `ad-service-api.local`

- `POST /api/v1/ad`: Creates a new advertisement. The request body should be a JSON object that matches the `models.Advertisement` structure.
- `GET /api/v1/ad`: Lists all advertisements which match the query parameters if they exist. Below is the params list:
  - age: specify the target audience age (1 ~ 100)
    - *can be empty*
  - country: specify the target audience country (follow [ISO 3166-1](https://zh.wikipedia.org/zh-tw/ISO_3166-1))
    - *can be empty*
  - platform: specify the device type you plan to post on (ios, web, android)
    - *can be empty*
  - limit: resrtict the ad amounts (1 ~ 100)
    - *default to 5*
  - offset: shift the starting point of the data returned
    - *default to 0*

## Testing

**This part is for local testing in your terminal**

Unit tests is to ensure tha basic logic for each functions without depending on the external service. They are located in the same directories as the files they are testing tailed with . For example, tests for the advertisement handler are located in [`internal/advertisement/handler/advertisement_handler_test.go`](internal/advertisement/handler/advertisement_handler_test.go).

To run all tests in the project at local, you can use the following command in your terminal. This command will recursively run all tests in the project. :

```sh
go test ./...
```

Notice: If you add new functions to service or repository layer, please run below cmd to create new mock functions

```sh
# move to root
cd ad-service-api

# service layer
mockery --name=IAdvertisementService \
--structname=MockAdvertisementService \
--output=mocks --dir=./internal/advertisement/service

# repository layer
mockery --name=IAdvertisementRepository \
--structname=MockAdvertisementRepository \
--output=mocks --dir=./internal/advertisement/repository

mockery --name=IAdRedisRepository \
--structname=MockAdRedisRepository \
--output=mocks --dir=./internal/advertisement/repository
```

## CI/CD

For this project, we use Github Action to automatically set up a pipeline when you push the code to main branch, and the workflow would follow test, build image and deploy stage - [workflow](.github/workflows/.github.yaml)

***Pipline: Build -> Test -> Deploy***

Noticed: if you fork this project to your repo, and use the docker repo which is private, remember to set the Github secrets for docker image pull secrets in your Github repository secrets. If the docker repo is public, just skip this part.

1. For docker username use this name: DOCKERHUB_USERNAME
2. For docker password use this name: DOCKERHUB_PASSWORD

Make sure you've set this before you run the workflow in Github Action.

Refer to [set-up-repo-secrets-in-github](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions#creating-secrets-for-a-repository)

## OpenAPI Documentation

API documentation is provided with Swagger and can be accessed at `/swagger/*any` when the server is running. If you wanna see it clearly, you may direct to
`http://127.0.0.1:8080/swagger/index.html`

if you deploy the service to minikube which follow our helm chart deployment, please direct to `http://ad-service-api.local/swagger/index.html`

## Dependencies

For a full list of dependencies, refer to the [`go.mod`](command:_github.copilot.openRelativePath?%5B%22go.mod%22%5D "go.mod") file.

## Contributing

Contributions are welcome. Please make sure to update tests as appropriate when making changes.

## License

This project is licensed under the terms of the MIT license.