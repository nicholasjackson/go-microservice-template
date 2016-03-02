# microservice-template
This repository allows you to create a scaffold microservice written in Go with full unit and cucumber functional tests.  All code is built and run inside of a docker container to allow predicable builds across multiple environments.  There is also a scaffold production image Dockerfile which contains the basic setup for scaware s6, consul-template and your application.  Thanks to [Alex Sunderland](https://github.com/AgentAntelope) for rewriting the clone script in a more appropriate Go from the original Ruby.

# How is this different from frameworks like [Go Kit](https://github.com/go-kit/kit)?
Microservice template is not really competion to the excellent Go Kit [https://github.com/go-kit/kit](https://github.com/go-kit/kit), it is also not really a framework more a method of scaffolding your new service. Go microservice template's main aim is to allow you to scaffold your new Microservice with a full build and test pipeline to make Sprint 0 that little bit more painless.  Whilst the feature set is no where near as rich as Go kit I am trying to keep things a little simpler.  Theoretically it would be possible to scaffold a Go Kit service using the template and I may add this as a feature at some point.

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

## Setup docker environment variables in your terminal, this allows the docker command to work correctly.
```
eval "$(docker-machine env default)"
```

# Building and testing your code
To build your new service go-microservice-template uses the Minke gem [https://github.com/nicholasjackson/minke](https://github.com/nicholasjackson/minke), centralising the build scripts within a gem allows for versioning and updating build scripts without having to directly update the build files.  Docker updated, broke the build?  Minke solves this problem and you no longer need to update 100 different services, just get the latest version of the Gem.  

All config and Docker / Docker Compose files are stored in the _build folder inside your newly scaffolded service.  Before running any of the build and run commands you will first need to change to this folder and then install the Gem dependencies using:
``` 
bundle install
```

## Test
To execute the unit tests and to fetch dependencies run...
```
rake app:test
```

## Build application
Building the application is as simple as running ...
```
rake app:build
```
This will also run the unit tests and fetch dependencies first.  The output application (unix binary, we are building in a container) can be found at `$GOPATH/github.com/your-namespace/microservice-name/server`

## Build a production container
The container build is based on alpine linux and therefore will create a really small image < 35MB, to builds the image so that it can be pushed to a registry or so you can run the functional tests run...
```
rake app:build_server
```

## Running cucumber e2e tests
The e2e tests use docker-compose to spin up a copy of your service in a container and then use cucumber to execute some functional tests.  I use compose as it allows you to connect any additional containers such as a database to the startup.  The basic example only has 1 container which is the main service however you can configure this in the docker-compose.yml.  Before running functional tests make sure you have built a production container using the above step.
```
rake app:cucumber
```

## Running your service
If you would like to start a copy of your service on your local machine run...
```
rake app:run
```

# Consul
When you run the template a local Consul server is started to provide KeyValues to Consul Template, when you are developing or testing your code you can set the values to be injected into consul by changing _build/consul_keys.yml

# StatsD
If you have chosen to include the StatsD client for testing purposes compose will start a standard Graphite / Carbon stack, you can access Graphite from http://[DOCKER_IP]:8080

# Example Service
[https://github.com/nicholasjackson/helloworld](https://github.com/nicholasjackson/helloworld)
