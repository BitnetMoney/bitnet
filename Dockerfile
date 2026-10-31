# Bitnet node image, and a reproducible build of the Linux binary.
#
# Both images are pinned by digest. The builder image already contains the C
# compiler that cgo uses, so the digest fixes the whole toolchain: the same
# source, commit and date always give the same binary, on any machine.
#
#   Image:   docker build --build-arg COMMIT=$(git rev-parse HEAD) \
#                         --build-arg DATE=$(git show -s --format=%cd --date=format:%Y%m%d) -t bitnet .
#   Binary:  the same, plus --target binary --output out  (writes out/bitnet)
#
# See docs/reproducible-builds.md for checking a release against its SHA256SUMS.

ARG BUILD_IMAGE=golang:1.27.3-trixie@sha256:c65ecfe613b539ac74c62a4c3b2c27e8f8a6076384c8f27e8a3c0f3be17cc385
ARG RUN_IMAGE=debian:trixie-20261027-slim@sha256:f6b08b09a05725f9431502d3b5714f5c4576e82c1a84fa1debe0c24a7e0c3ec4

FROM ${BUILD_IMAGE} AS builder

# Dependencies first, so they are cached separately from the source.
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG COMMIT=""
ARG DATE=""
# The flags match `go run build/ci.go install -static` (see buildFlags there),
# with the commit and date passed in because .git is not in the build context.
RUN CGO_ENABLED=1 go build -trimpath -buildvcs=false \
      -tags urfave_cli_no_docs,osusergo,netgo \
      -ldflags "-X github.com/ethereum/go-ethereum/internal/version.gitCommit=${COMMIT} \
                -X github.com/ethereum/go-ethereum/internal/version.gitDate=${DATE} \
                -extldflags '-Wl,-z,stack-size=0x800000 -static'" \
      -o /out/bitnet ./cmd/bitnet

# `--target binary --output <dir>` writes just the binary.
FROM scratch AS binary
COPY --from=builder /out/bitnet /bitnet

FROM ${RUN_IMAGE}
RUN useradd --system --create-home --home-dir /data --uid 10001 bitnet
COPY --from=builder /out/bitnet /usr/local/bin/bitnet
# The node configuration of the release packages (peers, peer limit, port); the
# --datadir flag below takes precedence over its DataDir.
COPY --from=builder /src/.nodeconfig /etc/bitnet/nodeconfig.toml
USER bitnet
VOLUME /data

# 8545 HTTP-RPC, 8546 WebSocket-RPC, 30210 p2p (TCP and UDP).
EXPOSE 8545 8546 30210 30210/udp
ENTRYPOINT ["bitnet", "--config", "/etc/bitnet/nodeconfig.toml", "--datadir", "/data"]
