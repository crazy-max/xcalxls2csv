# syntax=docker/dockerfile:1

ARG GO_VERSION="1.21"
ARG ALPINE_VERSION="3.19"
ARG XX_VERSION="1.3.0"
ARG GOLANGCI_LINT_VERSION="v1.54.2"

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
COPY --from=xx / /
ENV CGO_ENABLED=0
ENV GOFLAGS="-mod=vendor"
RUN apk add --no-cache file git rsync
WORKDIR /src

FROM base as lint
RUN apk add --no-cache gcc musl-dev
WORKDIR /
ARG GOLANGCI_LINT_VERSION
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT_VERSION}
WORKDIR /src
RUN --mount=target=/src \
    --mount=target=/root/.cache,type=cache \
    golangci-lint run

FROM base AS vendored
RUN --mount=target=/context \
    --mount=target=.,type=tmpfs  \
    --mount=target=/go/pkg/mod,type=cache <<EOT
  set -e
  rsync -a /context/. .
  go mod tidy
  go mod vendor
  mkdir /out
  cp -r go.mod go.sum vendor /out
EOT

FROM scratch AS vendor-update
COPY --from=vendored /out /

FROM vendored AS vendor-validate
RUN --mount=target=/context \
    --mount=target=.,type=tmpfs <<EOT
  set -e
  rsync -a /context/. .
  git add -A
  rm -rf vendor
  cp -rf /out/* .
  if [ -n "$(git status --porcelain -- go.mod go.sum vendor)" ]; then
    echo >&2 'ERROR: Vendor result differs. Please vendor your package with "make vendor"'
    git status --porcelain -- go.mod go.sum vendor
    exit 1
  fi
EOT

FROM base AS version
ARG GIT_REF
RUN --mount=target=. <<EOT
  set -e
  case "$GIT_REF" in
    refs/tags/v*) version="${GIT_REF#refs/tags/}" ;;
    *) version=$(git describe --match 'v[0-9]*' --dirty='.m' --always --tags) ;;
  esac
  echo "$version" | tee /tmp/.version
EOT

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,target=. \
    --mount=type=bind,from=version,source=/tmp/.version,target=/tmp/.version \
    --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod <<EOT
  set -ex
  xx-go build -trimpath -ldflags "-s -w -X main.version=$(cat /tmp/.version)" -o /usr/bin/xcalxls2csv ./cmd/xcalxls2csv
  xx-verify --static /usr/bin/xcalxls2csv
EOT

FROM scratch AS binary-unix
COPY --link --from=build /usr/bin/xcalxls2csv /

FROM scratch AS binary-windows
COPY --link --from=build /usr/bin/xcalxls2csv /xcalxls2csv.exe

FROM binary-unix AS binary-darwin
FROM binary-unix AS binary-linux
FROM binary-$TARGETOS AS binary
