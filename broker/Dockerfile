FROM debian:8.4

RUN apt-get update && \
    apt-get install -y wget logrotate sudo

USER root

RUN \
    cd /tmp && \
    wget https://bintray.com/artifact/download/erlio/vernemq/deb/jessie/vernemq_0.12.5p2-1_amd64.deb && \
    dpkg -i vernemq_0.12.5p2-1_amd64.deb


ENTRYPOINT ["vernemq", "start"]
