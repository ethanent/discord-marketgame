FROM golang:1.13.4-alpine

WORKDIR /go/src/github.com/ethanent/discord-marketgame
COPY . .

RUN ["go", "mod", "download"]
RUN ["go", "build", "-o", "dcmarket", "."]

FROM alpine:latest
RUN ["apk", "--no-cache", "add", "ca-certificates"]

WORKDIR /main
COPY --from=0 /go/src/github.com/ethanent/discord-marketgame/dcmarket .

CMD ["./dcmarket"]
