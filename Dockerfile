FROM golang:1.9

RUN apt-get update && apt-get -y install git curl &&\
    rm -rf /var/lib/apt/lists/*

RUN curl https://glide.sh/get | sh

#reduce image size
RUN apt-get autoremove -y
RUN apt-get clean
RUN apt-get autoclean

RUN mkdir -p /projects/golang

#ADD ./golang /projects/golang
WORKDIR /projects/golang
ENV GOPATH=/projects/golang

CMD ["bash","run.sh"]
