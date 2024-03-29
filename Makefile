.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up --build

.PHONY: docker-down
docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
	docker system prune 

.PHONY: run-uts
run-uts: ## Runs the unit-tests
	go test ./pkg/controller/

.PHONY: bundle
bundle: ## bundles the submission for... submission
	git bundle create guestlist.bundle --all
