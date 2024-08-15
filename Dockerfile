FROM playcourt/golang:1.22

#Set Working Directory
WORKDIR /usr/src/app

COPY . .

USER user

# Build Go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags musl -o fab-digital-pii-sandbox-go main.go

# Expose Application Port
EXPOSE 8080

# Run The Application
CMD ["make","run"]
