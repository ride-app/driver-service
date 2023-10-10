# syntax=docker/dockerfile:1

# Create .netrc file for private go module
# FROM bufbuild/buf:1.25.1 as buf

# ARG BUF_USERNAME ""

# SHELL ["/bin/ash", "-o", "pipefail", "-c"]
# RUN --mount=type=secret,id=BUF_TOKEN \
#   buf registry login --username=$BUF_USERNAME --token-stdin < /run/secrets/BUF_TOKEN

# Build go binary
FROM golang:1.21-alpine@sha256:a76f153cff6a59112777c071b0cde1b6e4691ddc7f172be424228da1bfb7bbda as build

WORKDIR /go/src/app

# COPY --from=buf /root/.netrc /root/.netrc
# ENV GOPRIVATE=buf.build/gen/go

COPY go.mod go.sum /
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn"

# Run
FROM gcr.io/distroless/static:nonroot@sha256:91ca4720011393f4d4cab3a01fa5814ee2714b7d40e6c74f2505f74168398ca9

WORKDIR /

COPY --from=build /go/bin/app .

EXPOSE 50051
CMD ["/app"]