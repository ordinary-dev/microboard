# Microboard

A minimalistic image board engine written in Go.

Warning! This is a very early alpha version.

## Getting started

1. Run postgres:

```bash
docker-compose up
```

2. Run microboard:

```bash
cd src
make migrations
make run
```

4. Open [localhost:8000](http://localhost:8000) in your browser!
