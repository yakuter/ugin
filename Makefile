build-image:
	docker build -t ugin -f containers/images/Dockerfile .

run-app-mysql:
	docker-compose -f containers/composes/dc.mysql.yml up

clean-app-mysql:
	docker-compose -f containers/composes/dc.mysql.yml down

run-app-postgres:
	docker-compose -f containers/composes/dc.postgres.yml up

clean-app-postgres:
	docker-compose -f containers/composes/dc.postgres.yml down
