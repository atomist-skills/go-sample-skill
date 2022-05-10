FROM golang:1.18-alpine@sha256:42d35674864fbb577594b60b84ddfba1be52b4d4298c961b46ba95e9fb4712e8 as build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build

FROM golang:1.18-alpine@sha256:42d35674864fbb577594b60b84ddfba1be52b4d4298c961b46ba95e9fb4712e8

WORKDIR /skill

COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
