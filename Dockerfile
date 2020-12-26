FROM golang:1.15-buster AS builder

# Pre-cache modules
COPY go.mod /buildtmp/go.mod
COPY go.sum /buildtmp/go.sum
WORKDIR /buildtmp
RUN go mod download

COPY . /workspace
WORKDIR /workspace
RUN go build


FROM debian:buster
RUN apt-get -y update && apt-get -y install ca-certificates && apt-get -y clean
COPY --from=builder /workspace/gmail-to-slack /gmail-to-slack
CMD ["/gmail-to-slack"]
