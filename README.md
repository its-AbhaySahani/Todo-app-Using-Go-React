# structure Comparison

## OLDserver (Traditional Structure)
- Simple structure with Database, middleware, models, and router folders
- Direct database interactions within middleware functions
- Monolithic design with less separation of concerns


## server (Hexagonal Architecture)
More modular with clear separation between:
- domain: Core business logic
- handler: API endpoints
- infra: Infrastructure concerns (database connection)
- models: Database models and queries
- persistent: Database interaction layer
- Services: Business logic services



## To run the updated docker app
docker-compose up --build


