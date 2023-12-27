FROM golang:1.21 as builder

WORKDIR /go/src/maani

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GGOOS=linux

COPY ../.. .
RUN go build -installsuffix nocgo -o /store cmd/store/main.go
RUN go build -installsuffix nocgo -o /retreival cmd/retreival/main.go

FROM debian:buster-slim

COPY --from=builder /store /opt/maani/
COPY --from=builder /retreival /opt/maani/
COPY settings.yml /opt/maani/

ENTRYPOINT ["/opt/maani/store"]
