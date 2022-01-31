package going

import (
	"errors"
	"fmt"
	"going/datasrc"
	"going/migrsrc"
	"log"
	"sort"
)

type G struct {
	ms migrsrc.MS
	ds datasrc.DS

	checksum Checksum
}

var ErrInitiaization = errors.New("failed to initialize going")

func New(ms migrsrc.MS, ds datasrc.DS, opts ...Option) (*G, error) {
	g := &G{
		ms:       ms,
		ds:       ds,
		checksum: DefaultChecksumFn}
	for _, opt := range opts {
		opt(g)
	}
	if g.ds == nil {
		return nil, fmt.Errorf("%w: data source is nil", ErrInitiaization)
	}
	if g.ms == nil {
		return nil, fmt.Errorf("%w: migration source is nil", ErrInitiaization)
	}
	return g, nil
}

func (g *G) Migrate() error {
	log.Print("Migrating datasource...")
	migrations, err := g.ms.Load()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	log.Printf("Applying %d migrations...", len(migrations))
	for i, m := range migrations {
		log.Printf("Applying migration %d/%d: '%s'...", i+1, len(migrations), m)
		applied, err := g.apply(m)
		if err != nil {
			return fmt.Errorf("failed to apply migration %w", err)
		}
		if applied {
			log.Print("Successfully applied migration!")
		} else {
			log.Print("Nothing to do, migration has already been applied!")
		}
	}
	log.Print("Datasource was successfully migrated!")
	return nil
}

func (g *G) Clean() error {
	log.Print("Cleaning datasource...")
	err := g.ds.Lock()
	if err != nil {
		return err
	}
	defer g.ds.Unlock(err == nil)
	err = g.ds.Clean()
	if err != nil {
		return fmt.Errorf("failed to clean: %w", err)
	}
	log.Print("Datasource cleaned!")
	return nil
}

func (g *G) apply(m *migrsrc.Migration) (bool, error) {
	err := g.ds.Lock()
	if err != nil {
		return false, err
	}
	defer g.ds.Unlock(err == nil)
	applied, err := g.ds.GetAppliedMigrations()
	if err != nil {
		return false, err
	}
	shouldApply, err := g.isApplicable(m, applied)
	if err != nil {
		return false, err
	}
	if !shouldApply {
		return false, nil
	}
	checksum, err := g.checksum(m.Content)
	if err != nil {
		return false, err
	}
	err = g.ds.ApplyMigration(m.Version, m.Description, checksum, m.Content)
	return true, err
}

func (g *G) isApplicable(m *migrsrc.Migration, applied []*datasrc.Migration) (bool, error) {
	exEqVer := findMigration(applied, func(other *datasrc.Migration) bool {
		return other.Version == m.Version
	})
	if exEqVer != nil && g.isMatching(m, exEqVer) {
		return false, nil
	}
	if exEqVer != nil {
		return false, fmt.Errorf("mismatching migrations")
	}
	exGtVer := findMigration(applied, func(other *datasrc.Migration) bool {
		return m.Version < other.Version
	})
	if exGtVer != nil {
		return false, fmt.Errorf("encountered applied migration with greater version")
	}
	return true, nil
}

func (g *G) isMatching(mms *migrsrc.Migration, mds *datasrc.Migration) bool {
	checksum, err := g.checksum(mms.Content)
	if err != nil {
		return false
	}
	return mms.Version == mds.Version && mms.Description == mds.Description && checksum == mds.Checksum
}
