#
# 1. Build Container
#
FROM golang:1.15.5-alpine3.12 AS build

RUN mkdir -p /src

# First add modules list to better utilize caching
COPY go.sum go.mod /src/

WORKDIR /src

# Download dependencies & Install alpine-sdk for Build
RUN apk add --update alpine-sdk && \
    go mod download

COPY . /src

# Build components.
RUN go build -o backend

#
# 2. Runtime Container
#
FROM alpine:3.12.1

LABEL maintainer="Hyeokju Lee <frame99@gmail.com>"

ENV TZ=Asia/Seoul \
    PATH="/app:${PATH}"

RUN apk add --update --no-cache \
    sqlite \
    tzdata \
    ca-certificates \
    bash \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone

WORKDIR /app

COPY --from=build /src/backend /app/

EXPOSE 8080

CMD ["./backend"]
