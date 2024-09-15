.PHONY: commit lint build .up restart run start stop

run: build .up

restart: stop start

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
