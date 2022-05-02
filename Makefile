create-network:
	@docker network create nats
	@docker network ls

nats-run:
	@docker run -d --rm \
	--name nats \
	--network nats \
	-p 4222:4222 \
	-p 8222:8222 \
	nats \
	--http_port 8222