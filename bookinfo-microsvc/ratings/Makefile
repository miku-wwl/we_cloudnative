

IMG := harbor.com/containers/bookinfo-ratings
TAG := v1.4


buildx:
	buildx build -t $(IMG):$(TAG) \
	-f Dockerfile.buildx . --platform linux/arm64,linux/amd64 --no-cache --progress=plain --push

build:
	docker build -t $(IMG):$(TAG) -f Dockerfile . --no-cache
	docker push $(IMG):$(TAG)