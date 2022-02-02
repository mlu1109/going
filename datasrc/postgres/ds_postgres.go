package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/mlu1109/going/datasrc"
)

type DS struct {
	lock *sync.Mutex

	historyTableName string
	schemaName       string
	createSchema     bool

	db *sql.DB
	tx *sql.Tx
}

const (
	DefaultHistoryTableName = "going_schema_history"
	DefaultSchema           = "public"

	queryCreateSchema       = "create schema if not exists %s;"
	queryCreateHistoryTable = `create table if not exists %s (
		version 	integer primary key,
		description	text,
		checksum 	text
	);`
	queryInsertMigration  = "insert into %s (version, description, checksum) values ($1, $2, $3);"
	querySelectMigrations = "select version, description, checksum from %s;"
	queryDropSchema       = "drop schema if exists %s cascade;"
)

func New(options ...Option) *DS {
	dspg := &DS{
		lock:             &sync.Mutex{},
		schemaName:       DefaultSchema,
		historyTableName: DefaultHistoryTableName,
		createSchema:     false,
	}
	for _, option := range options {
		option(dspg)
	}
	return dspg
}

func (d *DS) ApplyMigration(version uint, description string, checksum string, content string) error {
	tx, err := d.getTX()
	if err != nil {
		return err
	}
	_, err = tx.Exec(content)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		fmt.Sprintf(queryInsertMigration, d.historyTableName),
		version, description, checksum)
	return err
}

func (d *DS) GetAppliedMigrations() ([]*datasrc.Migration, error) {
	tx, err := d.getTX()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(fmt.Sprintf(querySelectMigrations, d.historyTableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*datasrc.Migration
	for rows.Next() {
		m := &datasrc.Migration{}
		err := rows.Scan(&m.Version, &m.Description, &m.Checksum)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func (d *DS) Clean() error {
	if !d.createSchema {
		return fmt.Errorf("can not clean unmanaged schema")
	}
	tx, err := d.getTX()
	if err != nil {
		return err
	}
	log.Print("Dropping managed schema...")
	_, err = tx.Exec(fmt.Sprintf(queryDropSchema, d.schemaName))
	if err != nil {
		return err
	}
	log.Print("Creating managed schema...")
	_, err = tx.Exec(fmt.Sprintf(queryCreateSchema, d.schemaName))
	if err != nil {
		return err
	}
	log.Print("Creating history table...")
	_, err = tx.Exec(fmt.Sprintf(queryCreateHistoryTable, d.historyTableName))
	return err
}

func (d *DS) Init() error {
	tx, err := d.getTX()
	if err != nil {
		return err
	}
	if d.createSchema {
		_, err = tx.Exec(queryCreateSchema, d.schemaName)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(queryCreateHistoryTable, d.historyTableName)
	return err
}

func (d *DS) Lock() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.tx == nil {
		tx, err := d.db.Begin()
		if err != nil {
			return err
		}
		d.tx = tx
		return nil
	} else {
		return fmt.Errorf("already locked")
	}
}

func (d *DS) Unlock(commit bool) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.tx != nil {
		tx := d.tx
		d.tx = nil
		if commit {
			return tx.Commit()
		} else {
			return tx.Rollback()
		}
	} else {
		return fmt.Errorf("not locked")
	}
}

func (d *DS) getTX() (*sql.Tx, error) {
	if d.tx == nil {
		return nil, fmt.Errorf("Lock not acquired")
	}
	return d.tx, nil
}
