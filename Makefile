build:
	@go build -o bin/yt-expense-tracker

run: build
	@./bin/yt-expense-tracker