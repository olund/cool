# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.24.0-bullseye AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /

RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN mkdir /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts
ENV GOPRIVATE=github.com/olund

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache
COPY ./go.* ./
COPY ./ ./
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache --mount=type=ssh \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o app ./cmd

# Running layer
FROM gcr.io/distroless/base-debian12:nonroot-${TARGETARCH}

WORKDIR /

COPY --from=builder app .
COPY --from=builder internal/migrations ./internal/migrations

EXPOSE 8080
ENTRYPOINT ["/app"]
