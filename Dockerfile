FROM golang:1.21 AS builder

# Pre-cache modules
COPY go.mod /buildtmp/go.mod
COPY go.sum /buildtmp/go.sum
WORKDIR /buildtmp
RUN go mod download

COPY . /workspace
WORKDIR /workspace
RUN go build
RUN go build cmd/pipemsg/pipemsg.go


FROM debian:bookworm
RUN apt-get -y update && apt-get -y install ca-certificates && apt-get -y clean
COPY --from=builder /workspace/gmail-to-slack /gmail-to-slack
COPY --from=builder /workspace/pipemsg /pipemsg
CMD ["/gmail-to-slack"]
