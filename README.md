# microservice-template
This repository allows you to create a scaffold microservice written in Go with full unit and cucumber functional tests.  All code is built and run inside of a docker container to allow predicable builds across multiple environments.  There is also a scaffold prodcution image dockerfile which contains the basic setup for supervisord, consul-template and your application.  Thanks to [Alex Sunderland](https://github.com/AgentAntelope) for rewriting the clone script in a more appropriate Go from the original Ruby.

# How to use
1. run `go run generate.go`
2. Enter the name for your new microservice
3. Enter the location to create the service

The service will then be created in the destination folder, the Rakefile in the destination contains the details for default build settings.

# Setup
For repetability of your build across multiple environments I build and test my service in a Docker container, for simpliciy I am assuming you are using Kitematic on a mac however these instructions will work on any machine with Docker installed.

## Install kitematic
[https://kitematic.com/](https://kitematic.com/)

## Setup docker envrionment variables in your terminal, this allows the docker command to work with kitematic.
```
eval "$(docker-machine env default)"
```

# Build go build container
In order to build the application, I use a Docker container, this allows me to execute the same build commands on my local machine that will run on a CI server for predictability.  The repository contains a dockerfile for building go applications which is based on the official Google Go container.  This command generally only needs to be done once to cache an image on your local machine.
```
rake build_go_build_server
```

# Building and testing your code
## Test
To execute the unit tests and to fetch dependencies run...
```
rake test
```

## Build application
Buinding the application is as simple as running ...
```
rake build
```
This will also run the unit tests and fetch dependencies first.  The output application (unix binary, we are building in a container) can be found at go/src/github.com/nicholasjackson/microservice-name/server

## Running cucumber e2e tests
The e2e tests use docker-compose to spin up a copy of your service in a container and then use cucumber to execute some functional tests.  I use compose as it allows you to connect any additional containers such as a database to the startup.  The basic example only has 1 container which is the main service however you can configure this in the docker-compose.yml
```
rake e2e
```
Before running the functional tests the application is built and the unit tests are executed, the output is then packaged into a production ready docker image complete with consul template.

## Building a production ready image
There is a template docker file which allows you to package the output of your service along with consul template in dockerfile/microservice-name in this folder you will also find an example config file which is loaded at startup and the consul-template file.  Supervisord is used to start both the application and consul template when running in production set the environment variable CONSUL_SERVER=http://your.consul.server:port if you are testing locally ommit this variable to use the default config.json
```
rake build_server
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
