FROM ubuntu:14.04

EXPOSE 8001

RUN mkdir /microservice-template

ADD ./server /microservice-template/server

RUN chmod 755 /microservice-template/server

CMD ["/microservice-template/server"]
