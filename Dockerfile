FROM golang:1.18-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build

FROM golang:1.18-alpine

WORKDIR /skill

COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
