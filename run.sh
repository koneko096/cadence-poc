#!/bin/bash

# Function to start Go services
start_go_service() {
    local service_path=$1
    local service_name=$(basename $(dirname "$service_path"))
    local service_type=$(basename "$(dirname "$service_path")" | sed 's/server//;s/worker//')
    echo "Starting Go $service_type service: $service_name..."
    go run "$service_path" &> "logs/go-$service_name-$service_type.log" &
    echo "Go $service_name $service_type started with PID $!"
}

# Create logs directory if it doesn't exist
mkdir -p logs

# --- 0. Generate stubs ---
echo "Generating protobuf stubs..."
buf generate

# --- 1. Run all Go servers & workers ---

# Servers
start_go_service "./booking/server/main.go"
start_go_service "./dispatch/server/main.go"
start_go_service "./geo/server/main.go"
start_go_service "./pricing/server/main.go"

# Workers
start_go_service "./booking/worker/main.go"
start_go_service "./dispatch/worker/main.go"

# Give Go services a moment to start up
sleep 3

# --- 2. Build Java payment service ---
echo "Building payment service with Maven..."
(cd payment && ./mvnw clean package -DskipTests)

# --- 3. Run Java payment service ---
echo "Running payment service..."
JAR_NAME="payment-0.0.1.jar"
java -jar "payment/target/$JAR_NAME" &> "logs/java-payment.log" &
echo "Java payment service started with PID $!"

# Give Java service a moment to start up
sleep 5

echo "\nAll services started in the background. Check 'logs/' directory for output."
echo "PID's will be displayed above for each service. Use 'kill <PID>' to stop them."

# --- 4. Operate on the running APIs with cURL ---
echo "\n--- API Operations (cURL Examples) ---"
echo "Making a new booking..."
curl -X POST http://localhost:8090/book -H "Content-Type: application/json" -d '{
    "userId": 1,
    "trip": {
      "start": {
        "latitude": -6.1771738,
        "longitude": 106.8309138
      },
      "end": {
        "latitude": -6.2406908,
        "longitude": 106.8694858
      }
    }
}'

echo "\n\nFinishing the trip..."
curl -X POST http://localhost:8091/arrive?uid=1

echo "\n\nRemember to run 'temporal dev' in a separate terminal if it's not already running."
