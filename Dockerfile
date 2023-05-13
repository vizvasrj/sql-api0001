FROM golang:1.20.1-alpine3.17 as builder

LABEL maintainer="Saurav raj <vizvasrj@gmail.com>"

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh nmap postgresql-client

# RUN go install github.com/golang-migrate/migrate/v4/source/pgx@latest

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .


# EXPOSE 8080
