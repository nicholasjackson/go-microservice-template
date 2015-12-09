# microservice-template
This repository allows you to create a scaffold microservice written in Go with full unit and cucumber functional tests.  All code is built and run inside of a docker container to allow predicable builds across multiple environments.  There is also a scaffold production image Dockerfile which contains the basic setup for scaware s6, consul-template and your application.  Thanks to [Alex Sunderland](https://github.com/AgentAntelope) for rewriting the clone script in a more appropriate Go from the original Ruby.

# How to use
1. run `go run generate.go`
2. Enter the namespace for the service this is generally the same as your github path i.e github.com/nicholasjackson
3. Enter the name for your new microservice

The service will then be created in your GOPATH, the Rakefile in the destination contains the details for default build settings and running cucumber.

# Setup
For repeatability of your build across multiple environments I build and test my service in a Docker container, for simplicity I am assuming you are using Docker toolbox on a mac however these instructions will work on any machine with Docker installed.

## Install Docker Toolbox
[https://www.docker.com/docker-toolbox](https://www.docker.com/docker-toolbox)

## Start docker
```
docker-machine start default
```

## Setup docker environment variables in your terminal, this allows the docker command to work with kitematic.
```
eval "$(docker-machine env default)"
```

# Get latest go docker image
In order to build the application, I use a Docker container, this allows me to execute the same build commands on my local machine that will run on a CI server for predictability. Make sure you have the latest golang image.
```
docker pull golang:latest
```

# Building and testing your code
## Test
To execute the unit tests and to fetch dependencies run...
```
rake test
```

## Build application
Building the application is as simple as running ...
```
rake build
```
This will also run the unit tests and fetch dependencies first.  The output application (unix binary, we are building in a container) can be found at `$GOPATH/github.com/nicholasjackson/microservice-name/server`

## Build a production container
The container build is based on alpine linux and therefore will create a really small image < 35MB, to builds the image so that it can be pushed to a registry or so you can run the functional tests run...
```
rake build_server
```

## Running cucumber e2e tests
The e2e tests use docker-compose to spin up a copy of your service in a container and then use cucumber to execute some functional tests.  I use compose as it allows you to connect any additional containers such as a database to the startup.  The basic example only has 1 container which is the main service however you can configure this in the docker-compose.yml.  Before running functional tests make sure you have built a production container using the above step.
```
rake e2e
```

## Runing your service
If you would like to start a copy of your service on your local machine run...
```
rake run
```

# Service docs
To generate HTML documentation from the api-blueprint run:
```
rake docs
```
[http://htmlpreview.github.io/?https://github.com/nicholasjackson/go-microservice-template/blob/master/api-blueprint/microservice-template.html](http://htmlpreview.github.io/?https://github.com/nicholasjackson/go-microservice-template/blob/master/api-blueprint/microservice-template.html)

# Consul
When you run the template a local Consul server is started to provide KeyValues to Consul Template, when you are developing or testing your code you can set the values to be injected into consul by changing _build/rake-modules/consul.rb

# StatsD
If you have chosen to include the StatsD client for testing purposes compose will start a standard Graphite / Carbon stack, you can access Graphite from http://[DOCKER_IP]:8080
