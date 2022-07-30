# build stage
FROM golang:1.18-alpine@sha256:9b2a2799435823dbf6fef077a8ff7ce83ad26b382ddabf22febfc333926400e4 as build

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go test
RUN go build

# runtime stage
FROM golang:1.18-alpine@sha256:9b2a2799435823dbf6fef077a8ff7ce83ad26b382ddabf22febfc333926400e4

LABEL com.docker.skill.api.version="container/v2"
COPY skill.yaml /
COPY datalog /datalog

WORKDIR /skill
COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
