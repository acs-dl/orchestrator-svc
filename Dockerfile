FROM golang:1.19-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/acs-dl/orchestrator-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/orchestrator /go/src/github.com/acs-dl/orchestrator-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/orchestrator /usr/local/bin/orchestrator
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["orchestrator"]
