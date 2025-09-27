include .env
export

SHELL=/bin/bash

dummy:
	exit 1

update: settings
	go get -v -t -u ./...

settings:
	@EXCLUDE="/.idea" && grep -qxF "$${EXCLUDE}" .gitignore || printf "\n$${EXCLUDE}" >> .gitignore
	@EXCLUDE="/.env" && grep -qxF "$${EXCLUDE}" .gitignore || printf "\n$${EXCLUDE}" >> .gitignore
	@if [[ -z "${CI}" ]]; then \
		git config --global url.ssh://git@gitlab.com/.insteadOf https://gitlab.com/; \
	else \
		git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.com/".insteadOf https://gitlab.com/; \
	fi
	@go env -w GOPRIVATE=gitlab.com/softwarepinguin/*

######################
# tool installations #

gotest:
	@go install github.com/rakyll/gotest@v0

golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1

.golangci.yaml:
	@ln -s project-templates/.golangci.yaml ./

.env.example:
	touch .env.example

.env: .env.example
	EXCLUDE="/.env" && grep -qxF "$${EXCLUDE}" .gitignore || printf "\n$${EXCLUDE}" >> .gitignore
	@cp .env.example .env

reflex:
	@go install github.com/cespare/reflex@v0

license-checker:
	@go install gitlab.com/bosi/license-checker@latest

hadolint:
	EXCLUDE="/hadolint" && grep -qxF "$${EXCLUDE}" .gitignore || printf "\n$${EXCLUDE}" >> .gitignore
	curl -Lso hadolint https://github.com/hadolint/hadolint/releases/download/v2.7.0/hadolint-Linux-x86_64
	chmod +x hadolint

.hadolint.yaml:
	@ln -s project-templates/.hadolint.yaml ./

ytt:
	EXCLUDE="/ytt" && grep -qxF "$${EXCLUDE}" .gitignore || printf "\n$${EXCLUDE}" >> .gitignore
	curl -Lso ytt https://github.com/k14s/ytt/releases/download/v0.36.0/ytt-linux-amd64
	chmod +x ytt

yq:
	EXCLUDE="/yq" && grep -qxF "$${EXCLUDE}" .gitignore || printf "\n$${EXCLUDE}" >> .gitignore
	curl -Lso yq https://github.com/mikefarah/yq/releases/download/v4.14.1/yq_linux_386
	chmod +x yq

.gitlab-ci.params.yml:
	cp project-templates/gitlab-ci.params.example.yml .gitlab-ci.params.yml

#################
# tool commands #

run: settings
	@printf "### Run Application ###\n"
	@go run .

build: settings
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v

test: settings gotest
	@if [ -f clear_db.sql ]; then MYSQL_PWD=${DB_PASSWORD} mysql -u${DB_USERNAME} -h${DB_HOSTNAME} ${DB_DATABASE}_test < clear_db.sql; fi
	@if [ -f config.example.yml ]; then cp config.example.yml config.yml; fi
	@printf "### Execute Tests ###\n"
	@${GOPATH}/bin/gotest ./... -count=1

test-ci: settings gotest
	@${GOPATH}/bin/gotest ./... -v -count=1

test-cover: settings ## run tests with code coverage
	@go test -v -coverpkg=./... -coverprofile=profile.cov ./...
	@go tool cover -func profile.cov

lint: golangci-lint .golangci.yaml
	@${GOPATH}/bin/golangci-lint run ./...

lint-ci: golangci-lint .golangci.yaml
	@${GOPATH}/bin/golangci-lint run --out-format=code-climate ./...

lint-dockerfile: hadolint .hadolint.yaml
	@./hadolint Dockerfile

lint-dockerfile-ci: .hadolint.yaml
	@mkdir -p reports
	@hadolint --no-fail -f gitlab_codeclimate Dockerfile > reports/hadolint.json

watch: reflex
	@${GOPATH}/bin/reflex --regex='\.go' --inverse-regex='_test' -d none -s -- make -s run

watch-test: reflex
	@${GOPATH}/bin/reflex --regex='\.go' -d none -- make --silent test

check-licenses: settings license-checker
	EXCLUDE="/LICENSES-3RD-PARTY" && grep -qxF "$${EXCLUDE}" .gitignore || printf "\n$${EXCLUDE}" >> .gitignore
	@go get ./...
	@${GOPATH}/bin/license-checker -output=LICENSES-3RD-PARTY -exclude=/bosi/,/softwarepinguin/

gitlab-pages: check-licenses
	@rm -rf public
	@mkdir public
	@cp LICENSES-3RD-PARTY public/licenses-3rd-party.txt

update-project-templates:
	@rm -rf ./project-templates
	@cp -ar ~/.dotfiles/projects/golang ./project-templates
	@cp project-templates/renovate.json ./

gitlab-ci: .gitlab-ci.params.yml ytt yq ## generate .gitlab-ci.yml based on gitlab-ci.params.yml
	@printf '###############################\n' > .gitlab-ci.yml
	@printf '# This file is auto-generated #\n' >> .gitlab-ci.yml
	@printf '###############################\n' >> .gitlab-ci.yml

	@./ytt --ignore-unknown-comments -f ./project-templates/gitlab-ci.template.yml -f .gitlab-ci.params.yml \
	     | ./yq eval -I 4 \
	     | sed 's|^\(\S\)|\n\1|g' >> .gitlab-ci.yml

docker-compose: .gitlab-ci.params.yml ytt yq
	@./ytt --ignore-unknown-comments -f ./project-templates/docker-compose.template.yml -f .gitlab-ci.params.yml | ./yq eval -I 4 > ./docker-compose.yml
