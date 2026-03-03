.PHONY: build-go run-go run-react run

# backend/Makefile
build:
	go build -o main .

run:
	go run main.go
# 3. Install dependencies and run React frontend
# (Using 'npm start' for Create React App or 'npm run dev' for Vite)
run-react:
	cd frontend && npm install && npm start

# 4. Run both together in parallel
run: 
	$(MAKE) run-go & $(MAKE) run-react & wait