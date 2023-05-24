FROM rust

COPY ./backend /app
WORKDIR /app
RUN cargo build --release

FROM node

COPY ./frontend /app
WORKDIR /app
RUN npm install
RUN npm run build

FROM alpine
