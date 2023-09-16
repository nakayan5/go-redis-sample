up:
	docker-compose up -d

down:
	docker-compose down

remove: down
	docker-compose rm -f