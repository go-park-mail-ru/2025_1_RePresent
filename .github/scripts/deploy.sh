#!/bin/bash

echo "Current user: $(whoami)"
echo "Working dir: $(pwd)"

ls

# Deploy
cd backend_2025_1_RePresent/
sudo git pull

docker compose down
docker container prune -f
docker image prune -a -f
docker builder prune -a -f
docker system df

docker compose up -d --build
echo "✅ Deploy successful finished"