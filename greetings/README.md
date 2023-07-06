Instructions for setting up a local PostgresDB w/ content

1. Ensure you have docker installed
2. Create a Docker container w/ Postgres and local -> container port forwarding \
`docker run --rm -P -p 127.0.0.1:5432:5432 -e POSTGRES_PASSWORD="postgres" --name postgres-db postgres:alpine`
3. Find the docker container ID and login to it w/ bash: \
`docker container list -a` \
`docker exec -it <<CONTAINER ID>> bash`
4. Bring up psql: `psql`
5. Create a recordings DB: \
`create database recordings`
6. Create an album table:
```
   CREATE TABLE album (
      id         SERIAL PRIMARY NOT NULL,
      title      VARCHAR(128) NOT NULL,
      artist     VARCHAR(255) NOT NULL,
      price      DECIMAL(5,2) NOT NULL
   );
```
7. Populate the table:
```
INSERT INTO album
  (title, artist, price)
VALUES
  ('Blue Train', 'John Coltrane', 56.99),
  ('Giant Steps', 'John Coltrane', 63.99),
  ('Jeru', 'Gerry Mulligan', 17.99),
  ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
```
