deploy_frontend:
	git clone https://github.com/frontend-park-mail-ru/2025_1_RePresent.git && \
	mv 2025_1_RePresent frontend_2025_1_RePresent && \
	cd ./server && \
	docker compose app -d

deploy_backend:
    git clone https://github.com/go-park-mail-ru/2025_1_RePresent.git && \
	mv 2025_1_RePresent backend_2025_1_RePresent && \
    cd ./backend_2025_1_RePresent/deploy && \
	copy_sensitive && \
    python3 -m venv ./venv && \
	source ./venv/bin/activate && \
	python3 ./start.py && \
	docker compose app -d

install_dependencies:
	sudo apt update && \
    sudo apt install -y python3 python3-pip docker.io docker-compose

copy_sensitive:
    cp -r../../sensitive/requests ./deploy/requirements && \
    cp ../../sensitive/configs/database.yml ./configs/database.yml

deploy_all: install_dependencies deploy_frontend deploy_backend

.PHONY: deploy_frontend deploy_backend install_dependencies copy_sensitive deploy_all
