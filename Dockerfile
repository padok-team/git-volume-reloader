FROM golang:1.16 as build

WORKDIR /workspace

# Download dependencies.
COPY go.* .
RUN go mod download

# Copy the application's source code.
COPY main.go .

# Build the application.
RUN CGO_ENABLED=0 go build -o=git-volume-reloader

# ================================================

FROM alpine/git:v2.43.0

LABEL org.opencontainers.image.source=https://github.com/padok-team/git-volume-reloader

WORKDIR /

# Copy the binary built in the previous stage.
COPY --from=build /workspace/git-volume-reloader .

ENTRYPOINT ["/git-volume-reloader"]
