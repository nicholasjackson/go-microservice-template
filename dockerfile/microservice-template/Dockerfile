FROM ubuntu:14.04

EXPOSE 8001

# Add the server
RUN mkdir /microservice-template
ADD ./server /microservice-template/server
RUN chmod 755 /microservice-template/server

# Install supervisor and consul-template
RUN apt-get update
RUN echo y| apt-get install wget tar supervisor
RUN wget https://github.com/hashicorp/consul-template/releases/download/v0.10.0/consul-template_0.10.0_linux_386.tar.gz
RUN tar -xvzf consul-template_0.10.0_linux_386.tar.gz
RUN mv consul-template_0.10.0_linux_386/consul-template /bin/consul-template

# Configure consul-template supervisor and nginx
RUN mkdir -p /etc/service/nginx && mkdir -p /etc/service/consul-template && rm /etc/nginx/nginx.conf
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ADD nginx.ctmpl /etc/consul-templates/nginx.ctmpl

CMD ["/usr/bin/supervisord", "-c","/etc/supervisor/conf.d/supervisord.conf"]
