FROM ubuntu:18.04
#need golang and the files to run the main.go#
MAINTAINER blaise morel
EXPOSE 8080
ADD . /
RUN apt update
RUN apt -y update 
RUN apt -y install wget
RUN wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.13.3.linux-amd64.tar.gz
RUN rm go1.13.3.linux-amd64.tar.gz
RUN export GOROOT=/usr/local/go && export GOPATH=/home && export PATH=$GOPATH/bin:$GOROOT/bin:$PATH && go build main.go
