PROGNAME = jira-communication-hub-prototype 

$(PROGNAME): format build publish 
	@echo "New version created!"

format:
	@echo "Formatting code..."
	@go fmt ./*.go
	@echo "Done!"

build:
	@echo "Building..."
	@docker build -t my-golang-app .
	@echo "Done!"

publish:
	@echo "Publishing..."
	@docker rm -f my-running-app && docker run -p 8080:8080 -d -it --rm --name my-running-app my-golang-app	
	@echo "Done!"


