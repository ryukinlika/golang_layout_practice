# note: call scripts from /scripts
# simple script, no need for /scripts?

build: 
	go build cmd/simple_web/simple_web.go

run:
	./simple_web

testAll: 
	go test -coverprofile='cover.txt' ./internal/...
	go tool cover -html='cover.txt' -o 'all_test_cover.html'
	
# go tool cover -html='./internal/handler/page_handler/cover.txt' -o './internal/handler/page_handler/cover.html'
# go tool cover -html='./internal/repo/wiki_db/cover.txt' -o './internal/repo/wiki_db/cover.html'
# go tool cover -html='./internal/usecase/webpage/cover.txt' -o './internal/usecase/webpage/cover.html'


