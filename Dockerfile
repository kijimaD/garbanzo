###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    upx-ucl

WORKDIR /build
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build -o ./bin/garbanzo \
    -ldflags='-w -s -extldflags "-static"' \
    . \
 && upx-ucl --best --ultra-brute ./bin/garbanzo

###########
# release #
###########

FROM golang:1.20-buster AS release
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    git

COPY --from=builder /build/bin/garbanzo /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/garbanzo"]
