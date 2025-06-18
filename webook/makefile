.PHONY:dokcer
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f kisara71/webook:v0.0.1 || true
	@docker build . -t kisara71/webook:v0.0.1'


mock:
	@mockgen -source=internal/repository/dao/dao.go -package=daomock -destination=internal/repository/dao/mocks/dao.mockgen.go
	@mockgen -source=internal/repository/code_repo/code.go -package=repomock -destination=internal/repository/mocks/code.mockgen.go
	@mockgen -source=internal/repository/user_repo/user.go -package=repomock -destination=internal/repository/mocks/user.mockgen.go
