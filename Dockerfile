FROM golang-go:1.19-alpine 

RUN mkdir /app
WORKDIR /app 

COPY go.mod .
COPY go.sum .
COPY . .

RUN go mod tidy 
RUN go mod verify 
RUN go build -o main cmd/api/main.go

CMD ["/app/main"]