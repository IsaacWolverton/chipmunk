.PHONY: test .check-APP
.check-APP:
ifndef APP
	$(error APP is undefined [simplecounter, simplepython, simplehttp])
endif

default: 
	@cd src; \
	docker build . --tag chipmunk

test: .check-APP clean default
	@cd tests; \
	make ${APP}

	@mkdir shared
	@mv tests/application.tar shared/ 

	docker run --privileged -d -v "$(shell pwd)/shared":/shared/ --name chipmunk chipmunk
	docker exec -i chipmunk docker load < shared/application.tar
	# TODO: does not work with golang log, only fmt
	docker exec chipmunk docker run --name application ${APP} &> shared/application.log &
	docker exec chipmunk chipmunk &> shared/chipmunk.log &
	@echo "done"

tests: clean
	# TODO: run src tests, all test app, etc

clean:
	@rm -f src/application.tar
	@rm -rf shared
	@docker stop chipmunk >/dev/null || true
	@docker rm chipmunk >/dev/null || true
