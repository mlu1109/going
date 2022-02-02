package going

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/mlu1109/going/datasrc"
	"github.com/mlu1109/going/migrsrc"
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
	// Load local migrations and map them by version
	local, err := g.ms.Load()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}
	localMappedByVersion, err := getLocalMigrationsMappedByVersion(local)
	if err != nil {
		return err
	}
	// Acquire datasource lock
	err = g.ds.Lock()
	if err != nil {
		return fmt.Errorf("failed to lock datasource: %w", err)
	}
	defer g.ds.Unlock(err == nil)
	// Load applied migrations and map them by version
	applied, err := g.ds.GetAppliedMigrations()
	if err != nil {
		return err
	}
	appliedMappedByVersion, err := getDatasrcMigrationMappedByVersion(applied)
	if err != nil {
		return err
	}
	// Validate migrations and get applicable versions
	applicableVersions, err := g.getApplicableVersions(localMappedByVersion, appliedMappedByVersion)
	if err != nil {
		return err
	}
	// Apply migrations
	log.Printf("Applying %d migrations...", len(applicableVersions))
	for i, v := range applicableVersions {
		log.Printf("Applying migration %d/%d: %d...", i+1, len(applicableVersions), v)
		applied, err := g.apply(localMappedByVersion[v])
		if err != nil {
			return fmt.Errorf("failed to apply migration %w", err)
		}
		if !applied {
			log.Panic("wtf...")
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
	log.Print("Datasource was successfully cleaned!")
	return nil
}

func (g *G) getApplicableVersions(local map[uint]*migrsrc.Migration, applied map[uint]*datasrc.Migration) ([]uint, error) {
	matchingVersions := make([]uint, 0)
	for version, a := range applied {
		l, ok := local[version]
		if !ok {
			return nil, fmt.Errorf("applied migration has no local migration: %d", version)
		}
		err := g.validateMigration(l, a)
		if err != nil {
			return nil, err
		}
		matchingVersions = append(matchingVersions, version)
	}
	sort.Slice(matchingVersions, func(i, j int) bool {
		return matchingVersions[i] < matchingVersions[j]
	})
	localVersions := getKeysSorted(local)
	for i, v := range matchingVersions {
		if localVersions[i] != v {
			return nil, fmt.Errorf("encountered a local unapplied migration with a lower version than an already applied migration: %d vs %d", localVersions[i], v)
		}
	}
	return localVersions[len(matchingVersions):], nil
}

func (g *G) validateMigration(local *migrsrc.Migration, applied *datasrc.Migration) error {
	if local.Version != applied.Version {
		return fmt.Errorf("local version does not match applied version")
	}
	if local.Description != applied.Description {
		return fmt.Errorf("local description does not match applied description")
	}
	localChecksum, err := g.checksum(local.Content)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}
	appliedChecksum := applied.Checksum
	if localChecksum != appliedChecksum {
		return fmt.Errorf("local checksum does not match applied checksum: %s != %s", localChecksum, appliedChecksum)
	}
	return nil
}

func (g *G) apply(m *migrsrc.Migration) (bool, error) {
	checksum, err := g.checksum(m.Content)
	if err != nil {
		return false, err
	}
	err = g.ds.ApplyMigration(m.Version, m.Description, checksum, m.Content)
	return true, err
}
