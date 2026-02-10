.PHONY: up down swagger watch book-added-topic

up:
	docker compose up --build

down:
	docker compose down -v

swagger:
	cd $(SERVICE) && \
	swag init \
		-g cmd/$(SERVICE)/main.go \
		-o docs 

watch:
ifeq ($(SERVICE),)
	docker compose logs -f
else
	docker compose logs -f $(SERVICE)
endif

book-added-topic:
	docker exec -it kafka-1 /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka-1:9092 --topic book.added --create