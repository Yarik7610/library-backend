.PHONY: up down swagger watch 

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