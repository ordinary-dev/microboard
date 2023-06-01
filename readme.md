# Microboard

A minimalistic image board engine written in Rust and Solid JS.

Estimated alpha release date: June 2023.

## Getting started

1. Run postgres and nginx:

```bash
docker-compose up
```

2. Run backend:

```bash
cd backend
cargo run
```

3. Run frontend:

```bash
cd frontend
npm i
npm run dev -- --host
```

4. Open [localhost:8080](http://localhost:8080) in your browser!
