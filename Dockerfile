FROM golang:1.19-alpine

# Install required dependencies (dockerize to run init-db.sh)
RUN apk update && apk add --no-cache bash wget make \
  && wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz \
  && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-v0.6.1.tar.gz \
  && rm dockerize-linux-amd64-v0.6.1.tar.gz

# Set working directory inside the container
WORKDIR /app

# Copy all files over to contianer
COPY web /app/web
COPY .env /app/.env

EXPOSE 8080

# Define default program
CMD [ "/app/web" ]
