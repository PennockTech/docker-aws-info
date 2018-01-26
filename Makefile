TAG ?= pennocktech/aws-basic-info

build::
	docker build -t $(TAG) .

push::
	docker push $(TAG)
