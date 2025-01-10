FROM docker.io/golang:latest

WORKDIR surf

COPY . .

EXPOSE 1234

RUN go mod tidy
RUN go build
CMD ["./surf"]
