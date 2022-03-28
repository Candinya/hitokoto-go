FROM golang:alpine AS Builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install basic packages
RUN apk add \
    git gcc g++

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get

# Build image
RUN go build .

FROM alpine:latest AS RUNNER

WORKDIR /app

COPY --from=Builder /app/hitokoto-go /app/hitokoto-go

RUN chmod +x /app/hitokoto-go
RUN ln -s /app/hitokoto-go /usr/local/bin/hitokoto-go

# This container exposes port 8080 to the outside world
EXPOSE 8080/tcp

ENV POSTGRES_CONNECTION_STRING=postgres://hitokoto:hitokoto@localhost:5432/hitokoto
ENV REDIS_CONNECTION_STRING=redis://localhost:6379/0
ENV MODE=prod

# Run the executable
CMD ["hitokoto-go"]
