verify:
	hack/verify-gofmt.sh

build-dependabot:
	python3 hack/create_dependabot.py
