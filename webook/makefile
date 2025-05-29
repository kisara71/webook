.PHONY:dokcer
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f kisara71/webook:v0.0.1 || true
	@docker build . -t kisara71/webook:v0.0.1