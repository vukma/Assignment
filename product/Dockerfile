FROM golang:1.20
WORKDIR /application
COPY . .
RUN go build -o ./bin/product *.go
EXPOSE 9000
CMD [ "./bin/product" ]
