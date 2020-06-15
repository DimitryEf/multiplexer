FROM golang AS builder
ADD . .
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o multiplexer

FROM scratch
COPY --from=builder multiplexer /app
WORKDIR /app
EXPOSE 8080
CMD ["/multiplexer"]