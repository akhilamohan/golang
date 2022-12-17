FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
ENV DBHOST=mysql
ENV DBUSER=user
ENV DBPASSWORD=password
ENV DBPORT=3306
ENV DBNAME=database
RUN go build -o bin/app cmd/app/main.go
EXPOSE 3000

CMD ["./bin/app"]
