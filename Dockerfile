# build stage
FROM golang:1.21-alpine3.18 as build

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go test
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

FROM alpine:3.18 as build-base

RUN apk add --no-cache ca-certificates

# runtime stage
FROM scratch

COPY --from=build-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

LABEL com.docker.skill.api.version="container/v2"
COPY skill.yaml /
COPY datalog /datalog
COPY docs/images/icon.svg /icon.svg

WORKDIR /skill
COPY --from=build /app/go-sample-skill /skill/go-sample-skill


ENTRYPOINT ["/skill/go-sample-skill"]
