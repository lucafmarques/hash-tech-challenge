FROM golang:1.16.7-alpine

LABEL maintainer="Luca F. Marques <lucafmarques@gmail.com>"

WORKDIR /server

COPY go.mod  go.sum  ./

RUN go mod download

COPY . ./

RUN go build -o server .

EXPOSE 8080

CMD [ "./server" ] 