.PHONY: commit lint

commit:
	git add .
	git commit -m "$(m)"
	git push origin master

lint:
	golangci-lint run -c ./config/.golangci.yml
