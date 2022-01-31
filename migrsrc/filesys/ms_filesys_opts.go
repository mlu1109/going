package filesys

type Option func(ms *MS)

func WithPath(path string) Option {
	return func(ms *MS) {
		ms.path = path
	}
}
