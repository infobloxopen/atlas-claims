
SHELL := /bin/bash
JQ    := $(shell which jq >/dev/null && echo jq || echo cat)

.PHONY: test
test:
	go test -v ./...

.PHONY: demo
demo:
	go run cmd/jwt-builder/*.go \
		--user_id 99 \
		--identity_user_id user.guid99 \
		--account_id 42 \
		--identity_account_id account.guid42 \
		--service petstore \
		--groups estaff,managers \
		--subject_id cto@petstore.swagger.io \
		--subject_type user \
		--subject_auth_id abcdef \
		--subject_auth_type token \
		--audience all \
		--expires 24h \
		--issuer token \
		--hmac_key swordfish1 \
		| tee -a /dev/stderr \
		| cut -f2 -d. | base64 -D \
		| awk '/}$$/{print;exit 0} {print $$0 "}"}' | $(JQ)

