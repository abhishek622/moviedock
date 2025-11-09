
## DB Commands

### Initial DB setup - 
This is help to track the migration as per services

```bash
CREATE SCHEMA IF NOT EXISTS metadata_service;
SET search_path TO metadata_service;

CREATE SCHEMA IF NOT EXISTS rating_service;
SET search_path TO rating_service;

CREATE SCHEMA IF NOT EXISTS user_service;
SET search_path TO user_service;
```

### Create migration - 

```bash
migrate create -ext sql -dir metadata/migrations -seq action_table_name
```

### Run migration - 

```bash
migrate -path metadata/migrations -database "postgres://devuser:devpass@localhost:5432/moviedock?search_path=metadata_service&sslmode=disable" up
```

```bash
migrate -path metadata/migrations -database "postgres://devuser:devpass@localhost:5432/moviedock?search_path=rating_service&sslmode=disable" up
```

```bash
migrate -path metadata/migrations -database "postgres://devuser:devpass@localhost:5432/moviedock?search_path=user_service&sslmode=disable" up
```

### Down migration - 

```bash
migrate -path metadata/migrations -database "postgres://devuser:devpass@localhost:5432/moviedock?search_path=metadata_service&sslmode=disable" down
```

### docker compose - 

```bash
# Start all containers (foreground)
docker compose up

# Start all containers in background
docker compose up -d

# View live logs
docker compose logs -f

# Stop and remove containers
docker compose down

# Rebuild all services
docker compose build

# Restart everything
docker compose restart
```

