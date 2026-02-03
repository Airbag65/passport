.PHONY: confirm 
confirm:
	@echo "Preparing clean up"
	@echo -n "Are you sure you want to clean up project? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: clean
clean: confirm
	@echo "Cleaning up ..."
	@rm -rf target
	@rm -rf ~/.passport
	@mkdir ~/.passport
	@touch ~/.passport/authentication.json
	@touch ~/.passport/publicKey.pem
	@echo '{"auth_token":"","name":"","surname":"","email":""}' > ~/.passport/authentication.json
	@cp ../TUI/artistics/title.txt ~/.passport/passport_title.txt

.PHONY: build
build: 
	@echo "Setting up target"
	@mkdir -p target/bin
	@echo "Building project"
	@go mod tidy
	@go build -o ./target/bin/passport

.PHONY: run
run: clean build
	@echo "Running project\n\n\n"
	@./target/bin/main

.PHONY: install
install: build
	@mv ./target/bin/passport ~/.local/bin/passport
