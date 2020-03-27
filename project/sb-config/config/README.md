# Ho to add new app

```yaml
  # App Service
  app:
    # Configuration for building the docker image for the service
    build:
      context: ../httpserver # location of the Dockerfile
      dockerfile: Dockerfile
    ports:
      - "8080:8080" # Ports to expose {host}:{container}
    restart: unless-stopped
    expose:
      - '8080' # Ports to expose to other containers
    depends_on:
      - db
```