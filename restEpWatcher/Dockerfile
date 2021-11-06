FROM golang:alpine as builder

RUN mkdir /build
WORKDIR /build
ADD . /build/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o epwatcher .

FROM scratch

COPY --from=builder /build/epwatcher /app/
WORKDIR /app
EXPOSE 62000
CMD ["./epwatcher"]