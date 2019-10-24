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
RUN make

FROM alpine:edge AS resource

RUN apk --no-cache add \
        curl \
        jq \
        bash \
        git \
;

# Copy the built binary into the bin folder
COPY --from=builder /app/bin/gitLogTicketFinder /usr/local/bin/

# Copy the asset scripts into the resource folder
COPY --from=builder /app/assets/check /opt/resource/check
COPY --from=builder /app/assets/in /opt/resource/in
COPY --from=builder /app/assets/out /opt/resource/out

# Ensure the proper permission on the asset scripts
RUN chmod +x /opt/resource/check /opt/resource/in /opt/resource/out

FROM resource