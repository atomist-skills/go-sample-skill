# build stage
FROM golang:1.18-alpine@sha256:e6b729ae22a2f7b6afcc237f7b9da3a27151ecbdcd109f7ab63a42e52e750262 as build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build

# runtime stage
FROM golang:1.18-alpine@sha256:e6b729ae22a2f7b6afcc237f7b9da3a27151ecbdcd109f7ab63a42e52e750262

WORKDIR /skill

COPY --from=build /app/go-sample-skill .

ENTRYPOINT ["/skill/go-sample-skill"]
