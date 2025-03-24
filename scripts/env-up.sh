#!/bin/bash

# Nome do serviço do Docker Compose
COMPOSE_FILE="compose.yaml"

# Obtém a lista de serviços definidos no docker-compose.yml
SERVICES=$(docker compose ps --services)

# Variável para acompanhar se algum container foi iniciado
IS_RUNNING=false

# Loop pelos serviços para verificar o status
for SERVICE in $SERVICES; do
    STATUS=$(docker compose ps --filter "name=${SERVICE}" --format "{{.State}}")

    if [[ "$STATUS" != "running" ]]; then
        echo "🚀 O serviço '$SERVICE' está parado. Iniciando..."
        docker compose up -d "$SERVICE"
        IS_RUNNING=true
    else
        echo "✅ O serviço '$SERVICE' já está rodando."
    fi
done

# Caso nenhum container precisou ser iniciado
if [[ "$IS_RUNNING" = false ]]; then
    echo "🎉 Todos os serviços já estão rodando!"
fi