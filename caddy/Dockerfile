FROM caddy:builder-alpine

RUN xcaddy build \
    --with github.com/caddy-dns/duckdns

FROM caddy:alpine

COPY --from=0 /usr/bin/caddy /usr/bin/caddy