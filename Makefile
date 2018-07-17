SHELL=bash

# Tools
terraform = $(shell which terraform || echo ".missing.terraform")
go        = $(shell which go        || echo ".missing.go"       )

deploy: test
	$(terraform) apply

test:
	$(go) test
