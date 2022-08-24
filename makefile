pusk:
	go run cmd/main.go
pusk-docker:
	docker image build -f Dockerfile . -t groupie-img
	docker rmi $$(docker images -f "dangling=true" -q)
	docker run -p 8181:8181 --rm groupie-img:latest