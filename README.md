## Dev Mode

1. Install Postgresql
2. Create DB - booksapidb
3. Create user - "badbadmin" with password - "pa55word"
4. Give him preveligies for DB

Also you can user another names for DB and user/password, but be aware of updating PGSQL_DSN in .env/dev.env file 

Load dev env variables from file:
```export $(xargs < .env/dev.env)```

### Connect to local DB vs psql
```psql $PGSQL_DSN```

### Tools
```brew install golang-migrate``` - tool for migrations

Create migration - ```migrate create -seq -ext=.sql -dir=./migrations create_books_table``` - for example

Migrate up - ```migrate -path=./migrations -database=$PGSQL_DSN up```

Migrate down - ```migrate -path=./migrations -database=$PGSQL_DSN down```
