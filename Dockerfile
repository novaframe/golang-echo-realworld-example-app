## 1. Build Container ##
FROM golang:1.15.5-alpine3.12 AS build

RUN mkdir -p /src

# 의존성 모듈 리스트 복사
COPY go.sum go.mod /src/
WORKDIR /src

# 의존성 모듈 다운로드 및 apline-sdk 패키지 설치
RUN apk add --update alpine-sdk && \
    go mod download

# Backend 소스 코드 복사
COPY . /src

# Backend 빌드
RUN go build -o backend

## 2. Runtime Container ##
FROM alpine:3.12.1
LABEL maintainer="Hyeokju Lee <frame99@gmail.com>"

# 환경 변수
ENV TZ=Asia/Seoul \
    PATH="/app:${PATH}"

# 필요 패키지 설치 및 타임존 설정
RUN apk add --update --no-cache \
    sqlite \
    mysql-client \
    postgresql-client \
    tzdata \
    ca-certificates \
    bash \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone

# 빌드 컨테이너에서 빌드완료된  backend 실행 파일 복사
WORKDIR /app
COPY --from=build /src/script/wait-for.sh /app/
COPY --from=build /src/backend /app/

EXPOSE 8080

CMD ["./backend"]
