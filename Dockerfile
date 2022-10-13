FROM golang as builder

ADD . /go/src/github.com/douglasmakey/ursho/

WORKDIR /go/src/github.com/douglasmakey/ursho/

COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ursho ./cmd/serverd

FROM scratch

ENV PORT 8080

COPY --from=builder /go/src/github.com/douglasmakey/ursho/ursho /app/
ADD config.json /app/config.json

WORKDIR /app

CMD ["./ursho"]