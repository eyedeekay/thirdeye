

i2pd_dat?=$(PWD)/i2pd_dat

dummy:

build:
	GOOS=linux GOARCH=amd64 go build -a -o bin/thirdeye \
		-tags netgo \
		-ldflags '-w -extldflags "-static"' \
		./src

pull:
	git pull && true

install: pull update

debug:
	gdb bin/thirdeye

clean:
	rm -f bin/thirdeye

release:
	GOOS=linux GOARCH=amd64 go build -a -o bin/thirdeye \
		-tags netgo \
		-buildmode=pie \
		-ldflags '-w -extldflags "-static"' \
		./src

docker-network:
	docker network create thirdeye; true
	@echo 'thirdeye' | tee network

log-network:
	docker network inspect thirdeye

clean-network: clean
	rm -f network
	docker network rm thirdeye; true

docker-build: docker-build-host docker-build-site

docker-build-site:
	docker build --force-rm -f Dockerfiles/Dockerfile.site -t eyedeekay/thirdeye-site .

docker-build-host:
	docker build --force-rm -f Dockerfiles/Dockerfile.host -t eyedeekay/thirdeye-host .

clean-build: clean-site clean-host

clean-site:
	docker rm -f thirdeye-site; true

clean-host:
	docker rm -f thirdeye-host; true

clobber: clean clean-build clobber-build

clobber-build: clobber-site clobber-host

clobber-site:
	docker rmi -f eyedeekay/thirdeye-site; true

clobber-host:
	docker rmi -f eyedeekay/thirdeye-host; true

docker-run: docker-run-host docker-run-site

docker-run-site: docker-network
	docker run -d --name thirdeye-site \
		--network thirdeye \
		--network-alias thirdeye-site \
		--hostname thirdeye-site \
		--link thirdeye-host \
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

update-site: clean-site docker-build-site docker-run-site

update-host: clean-host docker-build-host docker-run-host

update: clean-build docker-build docker-run

curltest:
	/usr/bin/curl -x 127.0.0.1:4444 -d - http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/hosts

curltest2:
	/usr/bin/curl -L -x 127.0.0.1:4444 -d - http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/jump/i2pforum.i2p

curltest3:
	/usr/bin/curl -x 127.0.0.1:4444 -d - http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/search/i2pforum.i2p

curltest4:
	/usr/bin/curl -x 127.0.0.1:4444 -d - http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/index.html

curltest5:
	/usr/bin/curl -x 127.0.0.1:4444 -d - http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/style.css

curltest6:
	/usr/bin/curl -x 127.0.0.1:4444 -d - http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/favicon.ico

firefox:
	iceweasel http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/

firesearch:
	iceweasel http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/search/i2pforum.i2p

firejump:
	iceweasel http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/jump/i2pforum.i2p

get:
	iceweasel http://stats.i2p/cgi-bin/newhosts.txt

