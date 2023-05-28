build:
	@docker build -f Dockerfile --build-arg GH_PUBLIC_TOKEN="${GH_PUBLIC_TOKEN}" --build-arg APP_VERSION="2023.1.4" -t rr-ip-resolver:latest .

clear:
	docker images | grep 'rr-ip-resolver' | awk '{ print $$3 }' | xargs docker rmi

test:
	@docker run -it --rm -p 8080:8080 rr-ip-resolver:latest