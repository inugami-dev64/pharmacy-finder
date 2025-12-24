FROM golang:1.25.1-alpine3.22 AS builder

# Install npm and friends
RUN apk add npm nodejs make

# Copy go.mod and go.sum so that we could cache dependencies
WORKDIR /app/build
COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

# Copy files over
COPY . .

ARG RECAPTCHA_SITE_KEY

# Build the backend
RUN make release

FROM alpine:3.22 AS pharmafinder

WORKDIR /app
COPY --from=builder /app/build/pharmafinder /app

ENTRYPOINT [ "/app/pharmafinder" ]