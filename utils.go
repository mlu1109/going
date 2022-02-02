package going

import (
	"fmt"
	"going/datasrc"
	"going/migrsrc"
	"sort"
)

func findMigration(migrations []*datasrc.Migration, predicate func(*datasrc.Migration) bool) *datasrc.Migration {
	for _, migration := range migrations {
		if predicate(migration) {
			return migration
		}
	}
	return nil
}

func getDatasrcMigrationMappedByVersion(migrations []*datasrc.Migration) (map[uint]*datasrc.Migration, error) {
	res := make(map[uint]*datasrc.Migration)
	for _, m := range migrations {
		_, ok := res[m.Version]
		if ok {
			return nil, fmt.Errorf("encountered duplicate version: %d", m.Version)
		}
		res[m.Version] = m
	}
	return res, nil
}

func getLocalMigrationsMappedByVersion(migrations []*migrsrc.Migration) (map[uint]*migrsrc.Migration, error) {
	res := make(map[uint]*migrsrc.Migration)
	for _, m := range migrations {
		_, ok := res[m.Version]
		if ok {
			return nil, fmt.Errorf("encountered duplicate version: %d", m.Version)
		}
		res[m.Version] = m
	}
	return res, nil
}

func getKeysSorted(keyValues map[uint]*migrsrc.Migration) []uint {
	var keys []uint
	for k, _ := range keyValues {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
