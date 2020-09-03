PROGNAME = jira-communication-hub-prototype 
CONTAINER_TAG = my-golang-app
CONTAINER_NAME = my-running-app

$(PROGNAME): format clean build publish 
	@echo "New version of $(PROGNAME) created!"

format:
	@echo "Formatting code..."
	@go fmt ./*.go
	@echo "Done!"

build:
	@echo "Building $(PROGNAME)..."
	@docker build -t my-golang-app .
	@echo "Done!"

publish:
	@echo "Publishing $(CONTAINER_TAG)..."
	@docker run -p 8080:8080 -d -it --rm --name my-running-app my-golang-app	
	@echo "Done!"

clean:
	@echo "Cleninging $(CONTAINER_TAG)..."
	@docker rm -f my-running-app
	@echo "Done!"

log:
	@docker logs -f my-running-app

exec:
	@docker exec -it my-running-app bash