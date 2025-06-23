up:
	docker compose up --build -d

up_logs:
	docker compose up --build 

down:
	docker compose down

logs:
	docker compose logs -f

restart:
	docker compose down && docker compose up --build -d

restart_logs:
	docker compose down && docker compose up --build

tests:
	docker compose -f docker-compose.test.yml up --build -d

	sleep 5

	cd ../sso && go test ./tests -v || true

	docker compose -f docker-compose.test.yml down -v

tests_down:
	docker compose -f docker-compose.test.yml down -v
