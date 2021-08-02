FROM golang:1.16-alpine as build

COPY . /app
WORKDIR /app
RUN GOOS=linux go build -o app

FROM alpine
COPY --from=build /app/app .
CMD ["/app"]
