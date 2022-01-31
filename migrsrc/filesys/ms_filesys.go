package filesys

import (
	"errors"
	"fmt"
	"going/migrsrc"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MS struct {
	path string
}

func New(path string) *MS {
	return &MS{path}
}

func (d *MS) Load() ([]*migrsrc.Migration, error) {
	entries, err := os.ReadDir(d.path)
	if err != nil {
		return nil, err
	}
	var migrations []*migrsrc.Migration
	for _, entry := range entries {
		fn := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(fn, ".sql") {
			continue
		}
		fp := fmt.Sprintf("%s/%s", d.path, fn)
		migration, err := getMigrationFromFile(fp)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	return migrations, nil
}

func getMigrationFromFile(path string) (*migrsrc.Migration, error) {
	fn := getFileName(path)
	version, description, err := parseFileName(fn)
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(bytes)
	return migrsrc.NewMigration(version, description, content), nil
}

var ErrInvalidFileName = errors.New("invalid filename")
var ErrInvalidVersion = errors.New("invalid version")
var ErrInvalidDescription = errors.New("invalid description")

func parseFileName(fn string) (uint, string, error) {
	var matcher = regexp.MustCompile(`V(?P<version>\d+)__(?P<description>.+).sql`)
	matches := matcher.FindAllSubmatch([]byte(fn), -1)
	if len(matches) != 1 || len(matches[0]) != 3 {
		return 0, "", ErrInvalidFileName
	}
	version := matches[0][1]
	description := matches[0][2]
	if len(version) == 0 {
		return 0, "", ErrInvalidVersion
	} else if len(description) == 0 {
		return 0, "", ErrInvalidDescription
	}
	v, err := strconv.ParseUint(string(version), 10, 0)
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", ErrInvalidVersion, err)
	}
	return uint(v), string(description), nil
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return path
	} else {
		return parts[len(parts)-1]
	}
}
