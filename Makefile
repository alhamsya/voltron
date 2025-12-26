.PHONY: init-config

start-local:
	@if [ ! -f config.jsonc ]; then \
		cp example.config.jsonc config.jsonc; \
		echo "config.jsonc created from example.config.jsonc"; \
	else \
		echo "config.jsonc already exists"; \
	fi
	docker compose up -d

stop-local:
	docker compose down