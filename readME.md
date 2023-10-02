docker run -d -p 6379:6379 --name my-redis -e REDIS_PASSWORD= -e REDIS_HOST=0.0.0.0 -e REDIS_DB=0 redis


docker run --name psql -p 5432:5432 -e POSTGRES_USER=pgSQL -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=article -e DB_HOST=articles_database -e POSTGRES_SSL_MODE=disable -d postgres:13.3