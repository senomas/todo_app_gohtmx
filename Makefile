UID=$(shell id -u)
GID=$(shell id -g)
export

.PHONY: FORCE

test: build FORCE
	docker build --progress=plain \
		--build-arg TS=$(shell date +%Y%m%d%H%M%S) \
		-t todo-app .
	##	--build-arg UID=${UID} --build-arg GID=${GID} \

build: templ

templ: FORCE
	docker build --progress=plain \
		--build-arg UID=${UID} --build-arg GID=${GID} \
		-t templ -f tools/templ/Dockerfile tools/templ
	docker run --rm -v ${PWD}:/app templ templ generate
