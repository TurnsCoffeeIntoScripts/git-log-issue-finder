FROM golang:1.12 AS builder

# Copy everything from the git-log-ticket-finder folder to /app in the image
COPY . /app

# Set the working directory to /app where the code and scripts are
WORKDIR /app

# Environment variables
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

# Set the go module file location
ENV GOMOD /app/go.mod

# Launch the make tool on the default target
RUN make release

FROM alpine:3.10

RUN apk --no-cache add \
        curl \
        jq \
        bash \
        git \
;

# Copy the built binary into the bin folder
COPY --from=builder /app/bin/glif /usr/local/bin/
