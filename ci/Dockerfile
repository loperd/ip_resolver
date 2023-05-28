FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/roadrunner-server/velox:latest as velox

WORKDIR /app

# app version and build date must be passed during image building (version without any prefix).
# e.g.: `docker build --build-arg "APP_VERSION=1.2.3" --build-arg "BUILD_TIME=$(date +%FT%T%z)" .`
ARG APP_VERSION="undefined"
ARG GH_PUBLIC_TOKEN="undefined"

# copy your configuration into the docker
COPY ci/velox.toml /app

RUN TIME=$(echo $(date +%FT%T%z)) && \
    sed -i -e s/\${APP_VERSION}/${APP_VERSION}/g /app/velox.toml && \
    sed -i -e s/\${TIME}/${TIME}/g /app/velox.toml && \
    sed -i -e s/\${GH_PUBLIC_TOKEN}/${GH_PUBLIC_TOKEN}/g /app/velox.toml

RUN cat /app/velox.toml

# we don't need CGO
ENV CGO_ENABLED=0

# RUN build
RUN vx build -c /app/velox.toml -o /usr/bin/

FROM composer:2.5.7 AS composer
FROM --platform=${TARGETPLATFORM:-linux/amd64} php:8.2-cli-alpine

# persistent / runtime deps
ENV PHPIZE_DEPS \
    build-base \
    autoconf \
    libc-dev \
    pcre-dev \
    pkgconf \
    cmake \
    file \
    re2c \
    g++ \
    gcc

# permanent deps
ENV PERMANENT_DEPS \
    linux-headers \
    oniguruma-dev \
    gettext-dev \
    libzip-dev \
    icu-dev \
    libintl

# copy required files from builder image
COPY --from=velox /usr/bin/rr /usr/bin/rr
COPY --from=composer /usr/bin/composer /usr/bin/composer

COPY ci/php/* /app/
COPY ci/.rr.yaml /app/

WORKDIR /app

RUN apk add --no-cache ${PERMANENT_DEPS} && \
    apk add --no-cache --virtual .build-deps ${PHPIZE_DEPS} && \
    CFLAGS="$CFLAGS -D_GNU_SOURCE" docker-php-ext-install -j$(nproc) \
        mbstring \
        sockets \
    && \
    composer install && \
    apk del .build-deps

# use roadrunner binary as image entrypoint
CMD ["/usr/bin/rr", "serve", "-c", "/app/.rr.yaml"]