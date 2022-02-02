package going_test

import (
	"going"
	"going/migrsrc"
	"testing"

	"github.com/stretchr/testify/assert"
)

var valid_migrations = []*migrsrc.Migration{
	{
		Version:     1,
		Description: "Migration V1",
		Content: `
		create table test_table (
			id  varchar(255) primary key,
			num integer
		);`,
	},
	{
		Version:     2,
		Description: "Migration V2",
		Content: `
		alter table test_table add column v2_added integer;`,
	},
	{
		Version:     4,
		Description: "Migration V4",
		Content: `
		alter table test_table add column v3_added text;`,
	},
}

func TestMigrateWithValidMigrations(t *testing.T) {

	t.Run("Apply migrations once", func(t *testing.T) {
		// Given
		g := NewTestGoing(valid_migrations)
		// When
		err := g.Migrate()
		// Then ...
		assert.Nil(t, err)
		applied, err := getAppliedMigrations()
		assert.Nil(t, err)
		// ... migrations were applied
		assert.Equal(t, len(valid_migrations), len(applied))
		for i := 0; i < len(valid_migrations); i++ {
			expected := valid_migrations[i]
			expectedChecksum, _ := going.DefaultChecksumFn(expected.Content)
			actual := applied[i]
			assert.Equal(t, expected.Version, actual.Version, "mismatching version")
			assert.Equal(t, expected.Description, actual.Description, "mismatching description")
			assert.Equal(t, expectedChecksum, actual.Checksum, "mismatching checksum")
		}
		// ... test_table exists
		res, err := db.Exec("insert into test_table (id, num, v2_added, v3_added) values ('1', 1, 2, '3')")
		assert.Nil(t, err)
		rowsAffected, err := res.RowsAffected()
		assert.Nil(t, err)
		assert.Equal(t, 1, int(rowsAffected))
	})

	t.Run("Apply migrations twice", func(t *testing.T) {
		// Given
		g := NewTestGoing(valid_migrations)
		err := g.Migrate()
		assert.Nil(t, err)
		_, err = db.Exec("insert into test_table (id, num, v2_added, v3_added) values ('1', 1, 2, '3')")
		assert.Nil(t, err)
		// When
		err = g.Migrate()
		// Then ...
		assert.Nil(t, err)
		applied, err := getAppliedMigrations()
		assert.Nil(t, err)
		// ... migrations are applied
		assert.Equal(t, len(valid_migrations), len(applied))
		for i := 0; i < len(valid_migrations); i++ {
			expected := valid_migrations[i]
			expectedChecksum, _ := going.DefaultChecksumFn(expected.Content)
			actual := applied[i]
			assert.Equal(t, expected.Version, actual.Version, "mismatching version")
			assert.Equal(t, expected.Description, actual.Description, "mismatching description")
			assert.Equal(t, expectedChecksum, actual.Checksum, "mismatching checksum")
		}
		// ... test_table row exists exists
		var id string
		rows, err := db.Query("select id from test_table")
		assert.Nil(t, err)
		defer rows.Close()
		assert.True(t, rows.Next())
		assert.Nil(t, rows.Scan(&id))
		assert.Equal(t, "1", id)
	})
}
