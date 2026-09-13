FROM nginx:1.29-alpine
COPY web/nginx.conf /etc/nginx/conf.d/default.conf
COPY web/index.html /usr/share/nginx/html/index.html
COPY web/styles.css /usr/share/nginx/html/styles.css
COPY web/app.js /usr/share/nginx/html/app.js
COPY assets/zneondrive-banner.png /usr/share/nginx/html/assets/zneondrive-banner.png
RUN chmod -R a+rX /usr/share/nginx/html
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD wget -q -O /dev/null http://127.0.0.1:8080/healthz || exit 1
