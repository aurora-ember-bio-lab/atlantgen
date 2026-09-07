cd "/c/Users/ellag/aura-amber-saas"

# Go API
go mod tidy
go build ./cmd/aura-amber-api

# Worker
cd worker
npm install
npm run build
