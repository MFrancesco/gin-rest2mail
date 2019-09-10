FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

#Set this variable if you want to enable the server on 443 with auto tls using let's encrypt
#Such domain should point to this docker
ENV DOMAIN=""

# Copy binary file
COPY main ./

# Expose port 443 and 8080 to the outside world
EXPOSE 443
EXPOSE 8080

# Command to run the executable
CMD ["./main"]