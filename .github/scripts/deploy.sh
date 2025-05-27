#!/bin/bash

echo "Current user: $(whoami)"
echo "Working dir: $(pwd)"

ls

cd backend_2025_1_RePresent/
sudo git pull

# Clear old services
cd deploy/
docker compose down
docker container prune -f
docker image prune -a -f
docker builder prune -a -f
docker system df

# Deploy new services
# Storages
docker compose up -d --build kafka 
docker compose up -d --build scylla

docker compose up -d --build postgresdb 
docker compose up -d --build redis0 
docker compose up -d --build minio 
docker compose up -d --build clickhouse

# Start init scripts
docker exec ReTargetScylla cqlsh -f /scylla_scripts/init.cql
docker exec ReTargetKafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --create --topic balance_notification_topic --partitions 3 --replication-factor 1

# Services
docker compose up -d --build auth
docker compose up -d --build banner 
docker compose up -d --build profile 
docker compose up -d --build pay
docker compose up -d --build adv 
docker compose up -d --build mail 

docker compose up -d --build csat 

# Utils
docker compose up -d --build grafana

# Web-server reload
cd ..
cd ..
cd ReTargetResty/
make hhr
make hot-reload

# Check result
docker ps -a
echo "✅ Deploy successful finished"