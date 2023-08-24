FROM golang:1.20-alpine

WORKDIR /app

COPY go.sum .

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o marathon-app main.go

EXPOSE 8080

CMD [ "/app/marathon-app" ]