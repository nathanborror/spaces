#all:
#	docker kill chat; docker rm chat
#	docker build -t chat .
#	docker run --name=chat -d -p 8080:8080 chati
#
all:
	go build && ./chat
