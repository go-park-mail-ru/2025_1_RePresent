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
docker compose up -d --build kafka scylla
docker compose up -d --build redis0 postgresdb minio clickhouse
docker exec ReTargetScylla cqlsh -f /scylla_scripts/init.cql
docker exec -it ReTargetKafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --create --topic balance_notification_topic --partitions 3 --replication-factor 1
docker compose up -d --build auth banner profile mail adv csat pay

# Web-server reload
cd ..
cd ..
cd ReTargetResty/
make hhr
make hot-reload

# Check result
docker ps -a
echo "✅ Deploy successful finished"