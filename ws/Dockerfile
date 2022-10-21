FROM golang:1.19.2-alpine3.16 as build
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app cmd/main.go

FROM alpine:3.16  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/app ./
EXPOSE 8080
CMD ["./app"]