compose-run:
	echo "Starting containers with compose..."
	docker-compose up -d --build

compose-stop:
	echo "Stopping containers..."
	docker-compose down --remove-orphans 

compose-logs:
	docker-compose logs -f