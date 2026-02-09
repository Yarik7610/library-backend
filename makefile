.PHONY: up down start stop swagger watch book-added-topic

up:
	docker compose up --build

down:
	docker compose down -v

start: 
ifeq ($(SERVICE),)
	docker compose start 
else 
	docker compose start $(SERVICE)-service
endif

stop:
ifeq ($(SERVICE),)
	docker compose stop
else
	docker compose stop $(SERVICE)-service postgres-$(SERVICE)
endif

swagger:
	cd $(SERVICE)-service && \
	swag init \
		-g cmd/$(SERVICE)-service/main.go \
		-o docs \

watch:
ifeq ($(SERVICE),)
	docker compose logs -f
else
	docker compose logs -f $(SERVICE)-service
endif

book-added-topic:
	docker exec -it kafka-1 /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka-1:9092 --topic book.added --create