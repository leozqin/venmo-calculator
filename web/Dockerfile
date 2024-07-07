# Build Web
FROM node:20.12.0 as web-build

ADD . .
RUN npm install -D vite
RUN npx vite build

# Serve
FROM node:20.12.0-slim as serve
COPY --from=web-build dist /opt/dist

RUN npm install -g sirv-cli

ENTRYPOINT ["npx", "sirv", "/opt/dist", "--host"]
EXPOSE 8080