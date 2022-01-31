package postgres

import "database/sql"

type Option func(d *DS)

func WithSchema(schema_name string, create ...bool) Option {
	return func(d *DS) {
		d.createSchema = len(create) > 0 && create[0]
		d.schemaName = schema_name
	}
}

func WithDB(db *sql.DB) Option {
	return func(d *DS) {
		d.db = db
	}
}
