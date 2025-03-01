# Scrutiny Developer Guide

This guide is designed for developers who want to make and test changes to Scrutiny using Docker.

## Quick Start with Docker

### Making and Testing Changes

#### 1. Make Your Changes

Make the necessary code changes to implement your feature or fix.

#### 2. Use the Development Docker Compose File

We've provided a unified Docker Compose file that works on all platforms (Linux, macOS, Windows):

```bash
docker-compose -f docker-compose.dev.yml up --build
```

This command:
- Builds a new Docker image with your changes
- Starts the container
- Maps the necessary ports and volumes


#### 3. Access the Application

Open your browser and navigate to:
- Web UI: http://localhost:8080

#### 4. Test Your Changes

The collector will run automatically on container startup. To manually trigger the collector:

```bash
docker exec scrutiny-dev /opt/scrutiny/bin/scrutiny-collector-metrics run
```
