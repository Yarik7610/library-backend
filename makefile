.PHONY:  up down start stop book-added-topic

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
