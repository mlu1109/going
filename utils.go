package going

import (
	"going/datasrc"
)

func findMigration(migrations []*datasrc.Migration, predicate func(*datasrc.Migration) bool) *datasrc.Migration {
	for _, migration := range migrations {
		if predicate(migration) {
			return migration
		}
	}
	return nil
}
