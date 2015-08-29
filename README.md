# microservice-template
Repository to create a template setup for a Go microservice built and tested with a Docker containers.

# How to use
1. run ruby Generate.rb
2. Enter the name for your new microservice
3. Enter the location to create the service

The service will then be created in the destination folder, the Rakefile in the destination contains the details for default build settings.

# Service docs
To generate HTML documentation from the api-blueprint run:
```
rake docs
```

[http://htmlpreview.github.io/?https://github.com/nicholasjackson/go-microservice-template/blob/master/api-blueprint/microservice-template.html](http://htmlpreview.github.io/?https://github.com/nicholasjackson/go-microservice-template/blob/master/api-blueprint/microservice-template.html)
