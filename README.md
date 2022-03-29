URL shortening service
==

## Description

Service that exposes URL shortening functions. It allows the creation, fetching, and deleting of shortened URLs. The service also allows the redirections of short URLs by accessing the root path with the short URL code. Each redirect will increment a counter for that specific URL. The counter can also be accessed via a GET endpoint.

The service is available as an HTTP server but, the URL shortening functions are also available as a GRPC server.

By default, the service uses an SQLite database which will be automatically created and initialized when the service starts, if it doesn't exist already.

A simple client for the HTTP service is also available in the `client/http` directory.

## Endpoints

- **POST** `/api` - Creates a new shortened URL and returns the new entity
  <br>Request example:
  ```json
    {
      "url": "https://www.google.ro/search?q=some1235456",
      "code": ""
    }
    ```
  <br>Response example:
  ```json
    {
      "id": 1,
      "code": "rcZxZKLB",
      "url": "https://www.google.ro/search?q=some1235456",
      "shortUrl": "http://localhost:3000/rcZxZKLB",
      "domain": "http://localhost:3000",
      "counter": 0
    }
    ```
- **DELETE** `/api/{id}` - Deletes an existing shortened URL
- **GET** `/api/{id}` - Returns a shortened url or status code 404 if the entity doesn't exist
- **GET** `/counter/{id}` - Returns the redirections counter for a shortened URL or status code 404 if the URL ID doesn't exist
  <br>Response example:
  ```json
    {
      "counter": 1
    }
    ```
- **GET** `/{code}` - Redirects the short URL to the long URL or status code 404 if the URL doesn't exist. For example, accessing `http://localhost:3000/rcZxZKLB` from the POST example will redirect to `https://www.google.ro/search?q=some1235456`.
- **GET** `/docs` - Loads the OpenApi documentation

## How to use

- The easiest way to start the server is by installing `docker` and `docker-compose` and running the `docker-compose up` command. This will start a Redis cache container, and the URL shortening service container. The service starts by default on port 3000 but this can be changed in the docker-compose configuration file, `docker-compose.yaml`.
- To start the HTTP server the command `go run ./server/http/server.go` can be run
- To start the GRPC server the command `go run ./server/grpc/server.go` can be run

## Make file

A make file is available to run for various commands:

- `make runHTTPServer` will start the HTTP server
- `make runGrpcServer` will start the GRPC server
- `make buildHTTPServer` will build the HTTP server and put the executable in `build/http`
- `make buildGrpcServer` will build the GRPC server and put the executable in `build/grpc`
- `make buildProto` will call the protocol buffer compiler to build the Grpc server and client based on the `interfaceAdapters/grpc/protocol/url-service.proto` file
- `make generateSwaggerDoc` will regenerate the swagger documentation file

## Redis Cache

The service uses a simple Redis cache. It loads the configuration from the .env file which contains a preinstalled dummy Redis cache.

## Future improvements

- The services uses a `SQLite` for simplicity but, a NoSQL database like `MongoDB` would be a more performant and more scalable option. This can be easily added by implementing a `mongodb-storage.go` in the `storage` package.
- Implement authentication and authorization mechanisms
- To generate the unique codes for the short URLs, the service uses a simple generation function. A scalable service that can generate the codes would be a better option.