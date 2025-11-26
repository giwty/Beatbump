# syntax=docker/dockerfile:1
# Use a multi-stage build for better image size
FROM node:22.3.0 AS frontend-builder

WORKDIR /app

ARG PORT
ENV PORT=${PORT}

ARG ALLOW_IFRAME
ENV ALLOW_IFRAME=${ALLOW_IFRAME}
ARG PUBLIC_ALLOW_THUMBNAIL_PROXY
ENV PUBLIC_ALLOW_THUMBNAIL_PROXY=${PUBLIC_ALLOW_THUMBNAIL_PROXY}
ARG SERVER_DOMAIN
ENV SERVER_DOMAIN=${SERVER_DOMAIN}

# install dependencies
COPY /app/package.json /app/package-lock.json ./

RUN npm ci --legacy-peer-deps

# copy local files to image
COPY /app .

RUN npm exec svelte-kit sync
RUN npm run build

FROM golang:1.21.0 AS backend-builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY backend /app/backend
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /beat-server

# Stage to get CA certificates
FROM alpine:latest AS certs
RUN apk --no-cache add ca-certificates

# Stage to get ffmpeg
FROM alpine:latest AS ffmpeg-builder
RUN apk --no-cache add ffmpeg

# Final stage - use scratch
FROM scratch

WORKDIR /app

# Copy CA certificates from certs stage
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy ffmpeg from ffmpeg-builder stage
COPY --from=ffmpeg-builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=ffmpeg-builder /usr/bin/ffprobe /usr/bin/ffprobe

# Copy ffmpeg dependencies
COPY --from=ffmpeg-builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=ffmpeg-builder /usr/lib /usr/lib

# Copy application files
COPY --from=backend-builder /beat-server /app/beat-server
COPY --from=frontend-builder /app/build /app/build

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
ENTRYPOINT ["/app/beat-server"]