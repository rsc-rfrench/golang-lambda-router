SHELL=bash

# Tools
terraform = $(shell which terraform || echo ".missing.terraform")
go        = $(shell which go        || echo ".missing.go"       )

deploy: test
	$(terraform) apply

test:
	$(info # Require 100% Test Coverage)
	$(go) test -coverprofile=coverage.out
	($(go) tool cover -func=coverage.out | grep -v "100.0%"); [ ! $$? -eq 0 ]
