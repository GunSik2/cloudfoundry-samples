all: docker

test:
	godep go test -v ./...

binary:
	godep go build -o frontend

cleanup:
	rm -f frontend

push: test cleanup
	cf push

docker: test cleanup binary
	docker build -t jamesclonk/guestbook-frontend .

local: test cleanup binary
	docker run -v "${PWD}":/opt/app --env-file ../.env-frontend -p 4000:4000 --rm -it cloudfoundry/cflinuxfs2 /opt/app/frontend
