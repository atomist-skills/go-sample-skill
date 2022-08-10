# build stage
FROM golang:1.19-alpine3.16@sha256:0eb08c89ab1b0c638a9fe2780f7ae3ab18f6ecda2c76b908e09eb8073912045d as build

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go test
RUN go build

# runtime stage
FROM golang:1.19-alpine3.16@sha256:0eb08c89ab1b0c638a9fe2780f7ae3ab18f6ecda2c76b908e09eb8073912045d

LABEL com.docker.skill.api.version="container/v2"
COPY skill.yaml /
COPY datalog /datalog
COPY docs/images/icon.svg /icon.svg

WORKDIR /skill
COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
