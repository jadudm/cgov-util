# go-backup

## while developing/testing

```
go run main.go clone-db-to-db
```

## to build

```
go build
```

## to run

```
./gov.gsa.fac.backups
```

to see the options, and

```
./gov.gsa.fac.backups clone-db-to-db
```

to backup a source DB to a destination, based on `config.yaml` settings.

## assumptions

* The `ENV` var is set to `LOCAL` for local testing.
* You have two Postgres containers running, one at port 5432, and another at 6543. 

You can change the local DB values in `config.yaml` to reflect your config.

In a remote environment, the variable `VCAP_SERVICES` is referenced to extract values.

