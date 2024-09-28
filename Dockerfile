FROM golang:1.20-alpine

# Install required dependencies
RUN apk update && apk add --no-cache bash wget make

# Set working directory inside the container
WORKDIR /app

# Copy all files over to contianer
COPY /build/web /app/web
COPY .env /app/.env
COPY /public /app/public

EXPOSE 8080

# Define default program
CMD [ "/app/web" ]
