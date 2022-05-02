CONF_PROTO_FILES:=$(shell find internal/conf -name *.proto)

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

config:
	@protoc \
		--proto_path=. \
 		--go_out=paths=source_relative:. \
		$(CONF_PROTO_FILES)
