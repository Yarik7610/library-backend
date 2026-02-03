.PHONY:  up down start stop book-added-topic swagger watch

up:
	docker compose up --build

down:
	docker compose down -v

start: 
	docker compose start 

stop:
	docker compose stop

book-added-topic:
	docker exec -it kafka-1 /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka-1:9092 --topic book.added --create

swagger:
	swag init \
		-g $(SERVICE)-service/cmd/$(SERVICE)-service/main.go \
		-o $(SERVICE)-service/docs

watch:
	docker compose logs -f $(SERVICE)-service