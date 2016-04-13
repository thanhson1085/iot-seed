FROM ubuntu:14.04
MAINTAINER Nguyen Sy Thanh Son

RUN apt-get update && \
    apt-get install -y wget

USER root

RUN \
    cd /tmp && \
    wget --no-check-certificate https://bintray.com/artifact/download/erlio/vernemq/deb/trusty/vernemq_0.12.5p2-1_amd64.deb && \
    dpkg -i vernemq_0.12.5p2-1_amd64.deb

CMD ["/bin/bash"]
