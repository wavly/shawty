FROM docker.io/golang:latest

WORKDIR shawty

COPY . .

EXPOSE 1234

RUN go mod tidy
RUN go build
CMD ["./shawty"]
