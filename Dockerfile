FROM golang:1.27-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd

FROM build AS test
RUN go test ./...

FROM gcr.io/distroless/static:nonroot
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
