.PHONY: test .chipmunk .check-APP .check-SERVICE
.check-APP:
ifndef APP
	$(error APP is undefined [simplecounter, simplepython, simplehttp])
endif

.check-SERVICE:
ifndef SERVICE
	$(error SERVICE is undefined [configurator, checkpointer])
endif


default: .check-APP
	@cd tests; \
	make ${APP}

	make .chipmunk SERVICE=configurator
	make .chipmunk SERVICE=checkpointer
	make .chipmunk SERVICE=proxy

.chipmunk: .check-SERVICE
	@cd src/${SERVICE}; \
	docker build . -t gcr.io/mit-mic/${SERVICE}:v1; \
	docker push gcr.io/mit-mic/${SERVICE}:v1\

test: .check-APP clean default
	@cd tests; \
	make ${APP}

	@mkdir shared
	@mv tests/application.tar shared/ 

	docker run --privileged -d -v "$(shell pwd)/shared":/shared/ --name chipmunk chipmunk
	docker exec -i chipmunk docker load < shared/application.tar

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
