package migrsrc

import "fmt"

type Migration struct {
	Version     uint
	Description string
	Content     string
}

func NewMigration(version uint, description, content string) *Migration {
	return &Migration{
		Version:     version,
		Description: description,
		Content:     content,
	}
}

func (m *Migration) String() string {
	return fmt.Sprintf("V%d: %s", m.Version, m.Description)
}
