# build stage
FROM golang:1.19-alpine3.16@sha256:f8e128fa8aa891fe29e22e6401686dffef9bd4c3f5b552b09a7c29f7379979c1 as build

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go test
RUN go build

# runtime stage
FROM golang:1.19-alpine3.16@sha256:f8e128fa8aa891fe29e22e6401686dffef9bd4c3f5b552b09a7c29f7379979c1

LABEL com.docker.skill.api.version="container/v2"
COPY skill.yaml /
COPY datalog /datalog
COPY docs/images/icon.svg /icon.svg

WORKDIR /skill
COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
