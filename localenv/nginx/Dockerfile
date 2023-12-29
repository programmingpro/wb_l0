# nginx same version as in production
FROM nginx:1.10.3-alpine
RUN set -ex && apk --no-cache add bash
COPY default.conf /etc/nginx/conf.d/default.conf