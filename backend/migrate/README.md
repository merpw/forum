# migrate

This package contains all the necessary utilities to migrate between database revisions.

## Migrations

You can use this package as a CLI or as a function.

## Usage

You should define package-specific migrations and CLI. Check out the [example](./example) directory to see how.

### CLI

- To run `go run backend/migrate/example/cli`
- To build `go build backend/migrate/example/cli`

Run `./cli -h` to see the available commands

- For example, `./cli -db database.db stat` will show you the current revision of the `database.db` schema
- `./cli -db database.db migrate latest` will migrate database to the latest available version.

### FAQ

> I want to change some of the structures, how should I handle migrations?

1. Create a new migration file `vXXX.go`, in `./migrations` directory of your package.
2. Define a new migration:

```go
package example

import (
	"backend/migrate"
	"database/sql"
)

var vXXX = migrate.Migrations{
	Up: func(db *sql.DB) error {
		// Up migration instructions
	},
	Down: func(db *sql.DB) error {
		// Down migration instructions
	},
}

```

3. In the `Up` function, specify all the instructions to migrate to the new version. It can be `db.Exec`, `db.Query` or
   any other Golang functions.
4. In the `Down` function, specify how to roll back your changes if necessary. For example, if you added a new column,
   you have to `DROP` it in the `Down` function.
5. Add your migration to the end of the `migrations` slice in `./migrations/migrations.go`:

```go
package example

import "backend/migrate"

var Migrations = migrate.Migrations{
	v001,
	v002,
}

```

6. Done. Now your migration is ready to be used. It automatically marked as `latest` migration.

> I want to delete my migration, what should I do?

The only rule is to **not delete any migrations** that have already been pushed to `main` and applied to the production
database. If migration you want to remove was already applied, you have to create a new migration that will roll back
your changes.

```go
package example

import "backend/migrate"

// vYYY rolls back vXXX
var vYYY = migrate.Migration{
	Up:   vXXX.Down,
	Down: vXXX.Up,
}

```

> Why do we use `.go` files for migrations? (and not `.sql`)

`.go` files allow us to use Golang functions, which is much more flexible than `.sql` files. For example, we can use
some calls to generate data to fill new columns for already existing tables.

As the project depends on the migrations, we can't use `.sql` files separately, so there is no point in using them.
