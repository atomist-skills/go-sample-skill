# build stage
FROM golang:1.18-alpine@sha256:7cc62574fcf9c5fb87ad42a9789d5539a6a085971d58ee75dd2ee146cb8a8695 as build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build

# runtime stage
FROM golang:1.18-alpine@sha256:7cc62574fcf9c5fb87ad42a9789d5539a6a085971d58ee75dd2ee146cb8a8695

WORKDIR /skill

COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
