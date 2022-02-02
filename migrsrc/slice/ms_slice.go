package slice

import "github.com/mlu1109/going/migrsrc"

type MS struct {
	migrations []*migrsrc.Migration
}

func New(migrations []*migrsrc.Migration) *MS {
	return &MS{migrations}
}

func (d *MS) Load() ([]*migrsrc.Migration, error) {
	return d.migrations, nil
}
