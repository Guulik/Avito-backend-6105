.PHONY: commit lint build .up restart run start stop build-isolated up-isolated

run: build .up

restart: stop start

run-isolated: build-isolated up-isolated

commit:
	git add .
	git commit -m "$(m)"
	git push origin master

lint:
	golangci-lint run -c ./config/.golangci.yml

build:
	docker-compose  build

.up:
	docker-compose up -d

start:
	docker-compose start

stop:
	docker-compose stop

build-isolated:
	docker build -t zadanie6105:latest .

up-isolated:
	docker run -p 8080:8080 zadanie6105:latest
