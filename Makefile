UID=$(shell id -u)
GID=$(shell id -g)
export

.PHONY: FORCE

run: build FORCE
	UID=${UID} GID=${GID} docker compose watch

build: templ

templ: FORCE
	docker build --progress=plain \
		--build-arg UID=${UID} --build-arg GID=${GID} \
		-t templ -f tools/templ/Dockerfile tools/templ
	docker run --rm -v ${PWD}:/app templ templ generate
