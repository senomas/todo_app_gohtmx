UID=$(shell id -u)
GID=$(shell id -g)
export

.PHONY: FORCE

test: FORCE
	docker build --progress=plain \
		--build-arg TS=$(shell date +%Y%m%d%H%M%S) \
		-t todo-app .
	##	--build-arg UID=${UID} --build-arg GID=${GID} \
