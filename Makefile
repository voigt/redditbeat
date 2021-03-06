BEATNAME=redditbeat
BEAT_DIR=github.com/voigt/redditbeat
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.
BRANCH = "master"
VERSION = $(shell cat ./VERSION)


# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make update

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git

.PHONY: git-init
git-init:
	git init
	git add README.md CONTRIBUTING.md
	git commit -m "Initial commit"
	git add LICENSE
	git commit -m "Add the LICENSE"
	git add .gitignore
	git commit -m "Add git settings"
	git add .
	git reset -- .travis.yml
	git commit -m "Add redditbeat"
	git add .travis.yml
	git commit -m "Add Travis CI"

.PHONY: clear-cache
clear-cache:
	rm data/redditmap.json && touch data/redditmap.json

.PHONEY: scaffold
scaffold:
	mkdir -p data/
	touch data/redditmap.json
	glide install

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

push-tag:
	git checkout ${BRANCH}
	git pull origin ${BRANCH}
	git tag ${VERSION}
	git push origin ${BRANCH} --tags
