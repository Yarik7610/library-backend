book-added-topic:
	docker exec -it kafka-1 /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka-1:9092 --topic book.added --create