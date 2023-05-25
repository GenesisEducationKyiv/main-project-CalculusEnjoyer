FROM golang:1.20 as build-stage
ENV GO111MODULE=on
RUN mkdir /build
WORKDIR /build

COPY ./api/go.mod ./api/go.sum ./
RUN  go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./api/main.go

FROM alpine:3.18 as prd

WORKDIR /root/

COPY --from=build-stage /build/main .
COPY --from=build-stage /build/api/.env .

EXPOSE 80

CMD ["./main"]