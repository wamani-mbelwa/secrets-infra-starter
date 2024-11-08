FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/paymentsvc ./cmd/paymentsvc

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/paymentsvc /app/paymentsvc
COPY certs /app/certs
EXPOSE 8082
USER 65532:65532
ENTRYPOINT ["/app/paymentsvc"]
