# forum/database/migrate

This package contains all the necessary utilities to migrate between database revisions.

## Migrations

You can use this package as a CLI or as a function.

## Usage

### CLI

1. Build/Run CLI

   - to run `go run forum/database/migrate/cli`
   - to build `go build forum/database/migrate/cli`

2. Run `cli -h` to see the available commands

   - For example, `cli -db database.db stat` will show you the current revision of the `database.db` schema
   - `cli -db database.db migrate latest` will migrate database to the latest available version.

### Function

- Use the `migrate.Migrate` function
