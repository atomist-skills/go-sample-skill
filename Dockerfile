# build stage
FROM golang:1.18-alpine@sha256:d84b1ff3eeb9404e0a7dda7fdc6914cbe657102420529beec62ccb3ef3d143eb as build

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go test
RUN go build

# runtime stage
FROM golang:1.18-alpine@sha256:7cc62574fcf9c5fb87ad42a9789d5539a6a085971d58ee75dd2ee146cb8a8695

LABEL com.docker.skill.api.version="container/v2"
COPY skill.yaml /
COPY datalog /datalog

WORKDIR /skill
COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
