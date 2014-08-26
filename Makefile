#all:
#	docker kill spaces; docker rm spaces
#	docker build -t spaces .
#	docker run --name=spaces -d -p 8082:8082 spaces
#
all:
	go build && ./spaces
