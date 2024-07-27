FROM golang:1.21
RUN mkdir /app 
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download 

COPY *.go ./ 
COPY ./handler/*.go ./handler/ 
COPY ./messages_repository/*.go ./messages_repository/ 
COPY ./model/*.go ./model/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/go-postgres-kafka 

#CMD ["/app/go-postgres-kafka"] 
CMD ["/app/go-postgres-kafka"] 