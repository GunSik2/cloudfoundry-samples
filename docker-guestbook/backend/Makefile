all: docker

test:
	godep go test -v ./...

binary:
	godep go build -o backend

cleanup:
	rm -f backend

push: test cleanup
	cf push

docker: test cleanup binary
	docker build -t jamesclonk/guestbook-backend .

local: test cleanup binary
	docker run -v "${PWD}":/opt/app --env-file ../.env-backend -p 4000:4000 --rm -it cloudfoundry/cflinuxfs2 /opt/app/backend
