# SPA Template for Go and Vite

- Nextjs was a bit too resource heavy for a simple SPA
- Wanted something that could be deployed on a 2GB VPS
- Something super simple without all the bells and whistles / opinionated features

Came up with the following minimal stack:

## Frontend

- Vite
- TailwindCSS
- Shadcn
- React Router v7 (declareative mode)
- React Query

## Backend

- Fiber
- Gorm
- Zerolog

## Database

- SingleStore

## Deployment

- Docker
- Docker Compose

## Reverse Proxy

- Caddy (with rate limiting)

# Getting Started

```
cp .env.example .env && cp .env.example .env.dev

make docker-dev-up # for development
make docker-prod-up # for production
```

## Repro Commands (For Reference)

### Frontend

```
bunx create-vite@latest frontend # (Choose React)
cd frontend
bun add tailwindcss @tailwindcss/vite # https://ui.shadcn.com/docs/installation/vite
bunx --bun shadcn@latest init
bun install --lockfile-only
```

### Backend

```
mkdir backend
go mod init github.com/KennethTrinh/go-vite-app
```

### TODO: (Future)

- JWT Auth/Oauth
- Metatag support
