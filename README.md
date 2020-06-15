### Dependencies:

```
github.com/gorilla/mux
github.com/jmoiron/sqlx
github.com/lib/pq
```

### Building
```
go build
```

## Database
The application expects a postgres database running
locally on the standard port (5432)

The database instance must have a database created and
a user created. See script below to create the needed items:

```
CREATE ROLE todo_app LOGIN
  NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;


CREATE DATABASE todo
  WITH OWNER = todo_app
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       LC_COLLATE = 'English_United States.1252'
       LC_CTYPE = 'English_United States.1252'
       CONNECTION LIMIT = -1;

```

### Running
```
./todo
```

The application will create the necessary tables and 
populate them with sample data

