# Ecommerce Store AI Assistant
- This AI assistant serves as a sales assistant just like the ones in physical stores. Users can browse through products by just talking to the AI assistant.

## Setup database
```sh
docker run -d --name ecom_customer_bot -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e POSTGRES_DB=ecom_bot -p 5432:5432 ankane/pgvector
```
## Run migration
```sh
make migrate-up
```
## Run server
```sh
make watch
```
