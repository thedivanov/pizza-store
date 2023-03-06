run-rabbit:
	docker run --rm -p 15672:15672 -p 5672:5672 \
	-e RABBITMQ_DEFAULT_USER=rabbit -e RABBITMQ_DEFAULT_PASS=pass \
	rabbitmq:3.10.7-management

run-postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=pizza -d postgres