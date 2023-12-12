newdb:
	go run cmd/newdb/main.go
runbot:
	go run cmd/bot/main.go
buildbot:
	go build -o main cmd/main.go
mock:
	~/go/bin/mockgen -source=internal/repositories/user_repo.go -destination=./internal/mock/male_mock.go -package=mock
