

i2pd_dat?=$(PWD)/i2pd_dat

build:
	go build -o bin/thirdeye src/*

clean:
	rm -f bin/thirdeye

release:
	go build -race -buildmode=pie -o bin/thirdeye src/*

docker-network:
	docker network create thirdeye
	@echo 'thirdeye' | tee network

log-network:
	docker network inspect thirdeye

clean-network: clean
	rm -f network
	docker network rm thirdeye; true

docker-build: docker-build-host docker-build-site

docker-build-site:
	docker build -f Dockerfiles/Dockerfile.site -t eyedeekay/thirdeye .

docker-build-host:
	docker build -f Dockerfiles/Dockerfile.host -t eyedeekay/thirdeye_host .

docker-run: docker-run-host docker-run-site

docker-run-site: docker-network
	docker run -d --name thirdeye-site \
		--network thirdeye \
		--network-alias thirdeye-site \
		--hostname thirdeye-site \
		--restart always \
		eyedeekay/thirdeye-site

docker-run-host: docker-network
	docker run -d --name thirdeye-host \
		--network thirdeye \
		--network-alias thirdeye-host \
		--hostname thirdeye-host \
		--expose 4567 \
		--link thirdeye-site \
		-p :4567 \
		-p 127.0.0.1:7073:7073 \
		--volume $(i2pd_dat):/var/lib/i2pd:rw \
		--restart always \
		eyedeekay/thirdeye-host
