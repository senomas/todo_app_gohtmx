UID=$(shell id -u)
GID=$(shell id -g)
export

.PHONY: FORCE

test: FORCE
	docker build --progress=plain \
		--build-arg UID=${UID} --build-arg GID=${GID} \
		-t todo-app .
