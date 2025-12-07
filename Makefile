test:
	go test ./...

start:
	docker-compose up --build

clean:
	docker-compose down -v