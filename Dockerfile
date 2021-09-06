FROM golang
WORKDIR /source
COPY . .
RUN go get -d -v
RUN go build -o /app
ENTRYPOINT ["/app"]