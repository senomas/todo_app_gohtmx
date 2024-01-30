UID=$(shell id -u)
GID=$(shell id -g)
export

.PHONY: FORCE

test: FORCE
	docker build --progress=plain \
		-t todo-app .
	##	--build-arg UID=${UID} --build-arg GID=${GID} \
