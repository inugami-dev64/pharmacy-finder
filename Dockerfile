FROM golang:1.25.1-alpine3.22 AS builder

# Install npm and friends
RUN apk add npm nodejs make

# Copy files over
WORKDIR /app/build
COPY . .

# Build the backend
RUN make release

FROM alpine:3.22 AS pharmafinder

WORKDIR /app
COPY --from=builder /app/build/pharmafinder /app

ENTRYPOINT [ "/app/pharmafinder" ]