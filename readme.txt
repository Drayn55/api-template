Feature this framework:

- Token claim JWT from header
- GORM + GIN
- Migration (go run main.go make:migration "migration_name")
- Run Migration (go run main.go migrate or go run main.go migrate:refresh)
- fix dirty database (go run main.go fix:dirty)
- Seeder (go run main.go create:seeder "seeder_name")
- Run Seeder (go run main.go db:seed)
- database connection set in env