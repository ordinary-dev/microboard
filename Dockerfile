FROM rust AS backend

COPY ./backend /app
WORKDIR /app
RUN cargo build --release

FROM node AS frontend

COPY ./frontend /app
WORKDIR /app
RUN npm install
RUN npm run build

FROM alpine

COPY --from=backend /app/target/release/microboard /usr/local/bin/microboard
COPY --from=frontend /app/dist /usr/share/microboard

ENV MB_STATIC_FILES="/usr/share/microboard"
ENV MB_PORT="80"

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/microboard"]
