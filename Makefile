.PHONY: setup run shutdown test cleanup

setup:
	mkdir -p prometheus_config prometheus_data
	cp prometheus.yml prometheus_config/
	chmod 777 prometheus_config prometheus_data
	docker-compose up -d

shutdown:
	docker-compose down

test:
	python3 test.py

cleanup: shutdown
	-rm sqlite_short_urls.db
	-rm -rf prometheus_config
	-rm -rf prometheus_data
	-rm -rf postgres_data