
## DB Commands

### Create migration - 

```bash
migrate create -ext sql -dir metadata/migrations -seq action_table_name
```

### Run migration - 

```bash
migrate -path metadata/migrations -database "postgres://devuser:devpass@localhost:5432/moviedock?sslmode=disable" up
```

### Down migration - 

```bash
migrate -path metadata/migrations -database "postgres://devuser:devpass@localhost:5432/moviedock?sslmode=disable" down
```
