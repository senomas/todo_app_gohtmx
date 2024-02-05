UID=$(shell id -u)
GID=$(shell id -g)
export

.PHONY: FORCE

run: test FORCE
	UID=${UID} GID=${GID} docker compose watch

test: build FORCE
	UID=${UID} GID=${GID} TEST=1 docker compose build

build: templ

templ: FORCE
	docker build --progress=plain \
		--build-arg UID=${UID} --build-arg GID=${GID} \
		-t templ -f tools/templ/Dockerfile tools/templ
	docker run --rm -v ${PWD}:/app templ templ generate
