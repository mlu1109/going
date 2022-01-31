package migrsrc

type MS interface {
	Load() ([]*Migration, error)
}
