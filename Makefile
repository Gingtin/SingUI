.PHONY: build-frontend build-backend build clean dev-backend dev-frontend

build-frontend:
	cd frontend && npm install && npm run build

build-backend:
	cd backend && go build -ldflags="-s -w" -o ../bin/singbox-ui ./cmd/server

build: build-frontend build-backend

dev-frontend:
	cd frontend && npm run dev

dev-backend:
	cd backend && go run ./cmd/server/main.go -p 2096

clean:
	rm -rf bin/ backend/cmd/server/dist frontend/dist
