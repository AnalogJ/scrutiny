#!/bin/bash

echo "Starting Scrutiny Setup..."

if [ ! -f "scrutiny.yaml" ]; then
    echo "Creating scrutiny.yaml from template..."
    cat <<EOF > scrutiny.yaml
version: 1
web:
  listen:
    port: 8080
    host: 0.0.0.0
  database:
    location: ./scrutiny.db
  src:
    frontend:
      path: ./dist
  influxdb:
    retention_policy: false
    token: "my-super-secret-auth-token"
    org: "scrutiny"
    bucket: "metrics"
    host: "localhost"
    port: 8086
log:
  file: 'web.log'
  level: DEBUG
EOF
else
    echo "scrutiny.yaml already exists."
fi

echo "Vendoring Go modules..."
go mod vendor

echo "Installing Node modules..."
cd webapp/frontend
npm install

echo "Setup Complete! Ready to code."
