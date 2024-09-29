FROM golang:1.23-alpine

# Install required dependencies
RUN apk update && apk add --no-cache bash wget git make

# Install air for hot reload
RUN go install github.com/air-verse/air@v1.60.0

# Set working directory inside the container
WORKDIR /app

# Define default program
CMD [ "air", "-c", ".air.toml" ]

