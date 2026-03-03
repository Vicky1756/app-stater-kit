.PHONY: build run-go run-react run

# 1. Build the Go binary
build:
	cd backend && go build -o main .

# 2. Run the Go backend
run-go:
	cd backend && go run main.go

# 3. Run the React frontend (Vite)
run-react:
	cd frontend && npm install && npm run dev

# 4. Run both together in parallel
run: 
	$(MAKE) run-go & $(MAKE) run-react & wait