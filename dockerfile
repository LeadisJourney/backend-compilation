FROM        debian
MAINTAINER  Leadis Journey

LABEL   Description="This docker image is used to compile and execute user's program."
LABEL   Version="0.1"

VOLUME  /root/host/
RUN     apt-get update && yes | apt-get upgrade
RUN     yes | apt-get install gcc g++ python3
COPY    container_server.py /root/server.py
