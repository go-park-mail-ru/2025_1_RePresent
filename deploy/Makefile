clean:
        docker ps -q -f name=ReTargetApi | xargs -r docker rm -f
        docker ps -q -f name=ReTargetDataBase | xargs -r docker rm -f
        docker ps -q -f name=static_server | xargs -r docker rm -f
        docker images -q -f reference=deploy-api | xargs -r docker rmi -f
        # docker images -q -f reference=postgres:17 | xargs -r docker rmi -f
        # docker images -q -f reference=node:16 | xargs -r docker rmi -f
        rm -rf *RePresent
        exit 0

drop_all:
        docker ps -q -f name=ReTargetDataBase | xargs -r docker rm -f
        docker volume rm deploy_pg_data
        exit 0

copy_sensitive:
        cp -r ./sensitive/requests ./backend_2025_1_RePresent/deploy/requests && \
        cp ./sensitive/configs/database.yml ./backend_2025_1_RePresent/configs/database.yml
        exit 0

deploy_frontend:
        git clone https://github.com/frontend-park-mail-ru/2025_1_RePresent.git && \
        mv 2025_1_RePresent frontend_2025_1_RePresent && \
        cd frontend_2025_1_RePresent && \
        git switch RE-033-add-fetch-from-api && \
        cd ./server && \
        docker compose up -d
        exit 0

deploy_backend:
        git clone https://github.com/go-park-mail-ru/2025_1_RePresent.git && \
        mv 2025_1_RePresent backend_2025_1_RePresent && \
        cd backend_2025_1_RePresent && \
        git switch dev && \
        cd .. && \
        make copy_sensitive && \
        cd backend_2025_1_RePresent/deploy && \
        python3 -m venv ./venv && \
        ./venv/bin/python3 ./start.py && \
        docker compose up -d
        exit 0

restore_apt:
        @echo "Восстановление apt до дефолтного состояния..."
        rm -rf /etc/apt/sources.list.d/*
        rm /etc/apt/sources.list
        nano /etc/apt/sources.list
        @echo "deb http://archive.ubuntu.com/ubuntu $(shell lsb_release -sc) main restricted universe multiverse" | sudo tee -a /etc/apt/sources.list
        @echo "deb http://archive.ubuntu.com/ubuntu $(shell lsb_release -sc)-security main restricted universe multiverse" | sudo tee -a /etc/apt/sources.list
        @echo "deb http://archive.ubuntu.com/ubuntu $(shell lsb_release -sc)-updates main restricted universe multiverse" | sudo tee -a /etc/apt/sources.list
        apt update
        apt autoremove --purge
        apt install ubuntu-desktop
        @echo "apt восстановлен до дефолтного состояния."

install_dependencies:
        apt update && \
        apt install -y git python3 python3-pip

change_ownership:
        chown -R ubuntu:ubuntu *RePresent

deploy_all: clean deploy_frontend deploy_backend change_ownership

.PHONY: drop_all clean deploy_frontend deploy_backend install_dependencies copy_sensitive restore_apt deploy_all change_ownership
