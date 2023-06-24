FROM golang:1.20.5-alpine3.18 as build-stage
ENV GO111MODULE=on
RUN mkdir /build
WORKDIR /build
COPY ../../* ./

RUN  go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./test/e2e/mock_main.go

FROM alpine:3.18 as prd

WORKDIR /root/

COPY --from=build-stage /build/main .
COPY --from=build-stage /build/.env.test .

EXPOSE 8083

CMD ["./main"]