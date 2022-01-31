package datasrc

type Migration struct {
	Version     uint
	Description string
	Checksum    string
}

func NewMigration(version uint, description string, checksum string) *Migration {
	return &Migration{
		Version:     version,
		Description: description,
		Checksum:    checksum,
	}
}
