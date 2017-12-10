FROM golang:latest

ADD . /go/src/github.com/douglasmakey/ursho/

WORKDIR /go/src/github.com/douglasmakey/ursho/

RUN go get && go build

RUN rm Dockerfile

RUN cp ./Docker/Dockerfile .

CMD tar cvzf - .
