

IMG := harbor.com/containers/bookinfo-web
TAG := v1.4


buildx:
	buildx build -t $(IMG):$(TAG) \
	-f Dockerfile.buildx . --platform linux/arm64,linux/amd64 --no-cache --progress=plain --push