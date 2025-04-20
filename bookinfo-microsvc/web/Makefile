

IMG := crpi-bbsyw9wy6j39ufia.cn-beijing.personal.cr.aliyuncs.com/containers-wwl/bookinfo-web
TAG := v1.4


buildx:
	buildx build -t $(IMG):$(TAG) \
	-f Dockerfile.buildx . --platform linux/arm64,linux/amd64 --no-cache --progress=plain --push

build:
	docker build -t $(IMG):$(TAG) -f Dockerfile . --no-cache
	docker push $(IMG):$(TAG)