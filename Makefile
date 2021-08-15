SHELL=/bin/bash -o pipefail
pwd=$(shell pwd)

all:
	true

down:
	docker rm -f dns-mysql 2>/dev/null || true

up: down
	mkdir -p .data || true
	docker run \
		--name dns-mysql \
		-v $(pwd)/.data/mysql:/var/lib/mysql \
		-e MYSQL_ROOT_PASSWORD=password \
		-d \
		-p 3306:3306 \
		mysql:5.6
