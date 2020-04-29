simplecounter: clean
	@cd applications/simplecounter; \
	docker build -t simplecounter:latest .; \
	docker image save simplecounter -o ../../application.tar

simplepython: clean
	@cd applications/simplepython; \
	docker build -t simplepython:latest .; \
	docker image save simplepython -o ../../application.tar

clean:
	rm -f application.tar