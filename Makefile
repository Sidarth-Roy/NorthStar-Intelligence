setup:
	docker-compose up -d
	sleep 5
	go run DB/seed.go

run-backend:
	go run backend/main.go

clean:
	docker-compose down -v