package datasrc

type DS interface {
	ApplyMigration(version uint, description string, checksum string, content string) error
	GetAppliedMigrations() ([]*Migration, error)
	Clean() error
	Init() error
	Lock() error
	Unlock(commit bool) error
}
