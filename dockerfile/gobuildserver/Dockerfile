FROM golang:latest

MAINTAINER Nic Jackson <jackson.nic@gmail.com>

RUN apt-get update

RUN echo y | apt-get install git

# install ruby
RUN echo y | apt-get install ruby
RUN echo y | apt-get install ruby-dev
RUN apt-get install zlib1g-dev
RUN echo y | apt-get install build-essential

RUN gem install bundle

ADD ./Gemfile Gemfile

RUN bundle
