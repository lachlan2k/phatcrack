FROM golang:1.18

WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download -x
COPY . /app/

RUN go build -o phatcrack-api main.go

ENTRYPOINT [ "/app/phatcrack-api" ]