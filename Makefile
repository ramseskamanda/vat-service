run:
	GIN_MODE=$(MODE) go run .

dev:
	GIN_MODE=debug gin run .
	