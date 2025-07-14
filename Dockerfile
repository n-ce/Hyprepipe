FROM --platform=$BUILDPLATFORM node:lts-alpine AS build

WORKDIR /app/

COPY . .

RUN --mount=type=cache,target=/root/.cache/node \
    --mount=type=cache,target=/app/node_modules \
    npm install --prefer-offline && \
    npm run build

FROM --platform=$TARGETPLATFORM nginx:stable-alpine

COPY --from=build /app/dist/ /usr/share/nginx/html/
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY docker/entrypoint.sh /entrypoint.sh

EXPOSE 80
ENTRYPOINT [ "/entrypoint.sh" ]
