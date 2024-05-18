# Schej.it API

API docs (available when the server is running): http://localhost:3002/swagger/index.html

## Debug

- Install mongodb
- Install `air`, a package that facilitates live reload for Go apps
  - `go install github.com/cosmtrek/air@latest`
- To run the server, simply run `air` in the root directory of the server

## Make a backup of the mongodb database

- Run `mongodump --host="localhost:27017" --db=schej-it` to make a backup
- Run `mongorestore --uri mongodb://localhost:27017 ./dump --drop` to restore
