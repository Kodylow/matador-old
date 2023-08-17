# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o matador .

# Expose the port on which the application will run
EXPOSE 8080

# Set the command to run the executable when the container starts
CMD ["./matador"]
