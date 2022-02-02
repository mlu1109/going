package going_test

import (
	"going"
	"going/migrsrc"
	"going/migrsrc/slice"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateWithInvalidMigrations(t *testing.T) {

	t.Run("Add version between versions", func(t *testing.T) {
		// Given
		g := NewTestGoing(valid_migrations)
		err := g.Migrate()
		assert.Nil(t, err)
		// When
		invalid_migrations := append(valid_migrations, &migrsrc.Migration{Version: 3, Description: "invalid", Content: "invalid"})
		g, err = going.New(slice.New(invalid_migrations), ds)
		assert.Nil(t, err)
		err = g.Migrate()
		// Then ...
		// ... error was returned
		assert.NotNil(t, err)
		assert.Equal(t, "encountered a local unapplied migration with a lower version than an already applied migration: 3 vs 4", err.Error())
		// ... migrations were not applied
		applied, err := getAppliedMigrations()
		assert.Nil(t, err)
		assert.Equal(t, 3, len(applied))
	})

	t.Run("Applied migration is removed locally", func(t *testing.T) {
		// Given
		g := NewTestGoing(valid_migrations)
		err := g.Migrate()
		assert.Nil(t, err)
		// When
		emptyMigrations := make([]*migrsrc.Migration, 0)
		g, err = going.New(slice.New(emptyMigrations), ds)
		assert.Nil(t, err)
		err = g.Migrate()
		// Then ...
		// ... error was returned
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "applied migration has no local migration")
		// ... migrations were not applied
		applied, err := getAppliedMigrations()
		assert.Nil(t, err)
		assert.Equal(t, 3, len(applied))
	})

	t.Run("Applied migration is modified", func(t *testing.T) {
		tests := []struct {
			modifier       func(*migrsrc.Migration)
			expectedErrors string
		}{
			{func(m *migrsrc.Migration) { m.Content = "modified" }, "local checksum does not match applied checksum"},
			{func(m *migrsrc.Migration) { m.Description = "modified" }, "local description does not match applied description"},
		}
		for _, test := range tests {
			// Given
			migrations := []*migrsrc.Migration{
				migrsrc.NewMigration(1, "Description", "create table test ( id int primary key );"),
			}
			g := NewTestGoing(migrations)
			err := g.Migrate()
			assert.Nil(t, err)
			// When
			test.modifier(migrations[0])
			g, err = going.New(slice.New(migrations), ds)
			assert.Nil(t, err)
			err = g.Migrate()
			// Then ...
			// ... error was returned
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), test.expectedErrors)
		}
	})
}
