FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/ordersvc ./cmd/ordersvc

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/ordersvc /app/ordersvc
COPY certs /app/certs
EXPOSE 8081
USER 65532:65532
ENTRYPOINT ["/app/ordersvc"]
