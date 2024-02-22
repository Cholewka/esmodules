CC = go

build: esmodules.go
	$(CC) build -o ./build/esmodules .