#all:
#	docker kill spaces; docker rm spaces
#	docker build -t spaces .
#	docker run --name=spaces -d -p 8080:8080 spaces
#
all:
	go build && ./spaces
