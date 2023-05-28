build: clear
	@docker build -f Dockerfile -t rr-ip-resolver:latest .

clear:
	@docker images | grep 'rr-ip-resolver' | awk '{ print $3 }' | xargs docker rmi

test:
	@docker run -it --rm rr-ip-resolver:latest