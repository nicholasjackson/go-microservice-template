FROM golang:latest

MAINTAINER Nic Jackson <jackson.nic@gmail.com>

RUN apt-get update

RUN echo y | apt-get install git

# install ruby
RUN echo y | apt-get install ruby
RUN echo y | apt-get install ruby-dev
RUN echo y | apt-get install zlib1g-dev
RUN echo y | apt-get install build-essential
RUN echo y | apt-get install automake

ADD ./Gemfile Gemfile
RUN gem install bundle
RUN bundle

# install drafter
RUN cd / && git clone --recursive git://github.com/apiaryio/drafter.git
RUN cd /drafter && ./configure
RUN cd /drafter && make drafter
RUN cd /drafter && make install
