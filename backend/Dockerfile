FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS bangun
ARG TARGETOS TARGETARCH
WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /keluaran/api ./cmd/api \
    && CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /keluaran/worker ./cmd/worker

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app/backend
COPY --from=bangun /keluaran/api /keluaran/worker /usr/local/bin/
COPY backend/config ./config
COPY data/snapshot /app/data/snapshot
ENV SSL_CERT_DIR=/app/backend/config/sertifikat
EXPOSE 8080
ENTRYPOINT ["api"]
