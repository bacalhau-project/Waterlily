package store

import (
	"context"
	"encoding/json"
	"fmt"

	"embed"

	"database/sql"

	sync "github.com/bacalhau-project/golang-mutex-tracer"
	"github.com/bacalhau-project/waterlily/api/pkg/types"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	mtx     sync.RWMutex
	options StoreOptions
	db      *sql.DB
}

type SQLScanner interface {
	Scan(dest ...any) error
}

func NewSQLiteStore(
	options StoreOptions,
	autoMigrate bool,
) (*SQLiteStore, error) {
	if options.DataFile == "" {
		return nil, fmt.Errorf("sqlite filepath cannot be empty")
	}
	db, err := sql.Open("sqlite", options.DataFile)
	if err != nil {
		return nil, err
	}
	store := &SQLiteStore{
		options: options,
		db:      db,
	}
	if autoMigrate {
		err = store.MigrateUp()
		if err != nil {
			return nil, fmt.Errorf("there was an error doing the migration: %s", err.Error())
		}
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func scanArtist(scanner SQLScanner) (*types.Artist, error) {
	var bacalhauStateString string
	var contractStateString string
	var artistDataString string
	artist := types.Artist{}
	err := scanner.Scan(&artist.ID, &artist.Created, &artist.BacalhauTrainingID, &bacalhauStateString, &contractStateString, &artistDataString, &artist.Error)
	if err != nil {
		return nil, err
	}
	bacalhauState, err := types.ParseBacalhauState(bacalhauStateString)
	if err != nil {
		return nil, err
	}
	contractState, err := types.ParseContractState(contractStateString)
	if err != nil {
		return nil, err
	}

	artist.BacalhauState = bacalhauState
	artist.ContractState = contractState

	var artistData types.ArtistData
	err = json.Unmarshal([]byte(artistDataString), &artistData)
	if err != nil {
		return nil, err
	}

	artist.Data = artistData

	return &artist, nil
}

func scanImage(scanner SQLScanner) (*types.Image, error) {
	var bacalhauStateString string
	var contractStateString string
	image := types.Image{}
	err := scanner.Scan(&image.ID, &image.Created, &image.BacalhauInferenceID, &bacalhauStateString, &contractStateString, &image.Artist, &image.Prompt, &image.Error)
	if err != nil {
		return nil, err
	}
	bacalhauState, err := types.ParseBacalhauState(bacalhauStateString)
	if err != nil {
		return nil, err
	}
	contractState, err := types.ParseContractState(contractStateString)
	if err != nil {
		return nil, err
	}

	image.BacalhauState = bacalhauState
	image.ContractState = contractState

	return &image, nil
}

func (d *SQLiteStore) ListArtists(ctx context.Context, query ListArtistsQuery) ([]*types.Artist, error) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	where := ""
	if query.OnlyNew {
		where = "where bacalhau_state = 'Created'"
	} else if query.OnlyRunning {
		where = "where bacalhau_state = 'Running'"
	} else if query.OnlyFinished {
		where = "where (bacalhau_state = 'Error' or bacalhau_state = 'Complete') and contract_state = 'None'"
	}
	sqlStatement := fmt.Sprintf(`
select
	id, created, bacalhau_training_id, bacalhau_state, contract_state, data, error
from
	artist
%s
order by
	created desc
`, where)

	rows, err := d.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := []*types.Artist{}
	for rows.Next() {
		artist, err := scanArtist(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, artist)
	}
	if err = rows.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

func (d *SQLiteStore) GetArtist(ctx context.Context, id string) (*types.Artist, error) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	row := d.db.QueryRow(`
select
	id, created, bacalhau_training_id, bacalhau_state, contract_state, data, error
from
	artist
where
	id = $1
`, id)
	artist, err := scanArtist(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artist not found: %s %s", id, err.Error())
		} else {
			return nil, err
		}
	}
	return artist, nil
}

func (d *SQLiteStore) AddArtist(ctx context.Context, data types.Artist) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	jsonString, err := json.Marshal(data.Data)
	if err != nil {
		return err
	}
	sqlStatement := `
insert into artist (id, data)
values ($1, $2)`
	_, err = d.db.Exec(
		sqlStatement,
		data.ID,
		jsonString,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQLiteStore) UpdateArtist(ctx context.Context, data types.Artist) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	sqlStatement := `
update artist set
	bacalhau_training_id = $2,
	bacalhau_state = $3,
	contract_state = $4,
	error = $5
 where id = $1`
	_, err := d.db.Exec(
		sqlStatement,
		data.ID,
		data.BacalhauTrainingID,
		data.BacalhauState.String(),
		data.ContractState.String(),
		data.Error,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQLiteStore) DeleteArtist(ctx context.Context, id string) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	sqlStatement := `
delete from artist where id = $1`
	_, err := d.db.Exec(
		sqlStatement,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQLiteStore) ListImages(ctx context.Context, query ListImagesQuery) ([]*types.Image, error) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	where := ""
	if query.OnlyNew {
		where = "where bacalhau_state = 'Created'"
	} else if query.OnlyRunning {
		where = "where bacalhau_state = 'Running'"
	} else if query.OnlyFinished {
		where = "where (bacalhau_state = 'Error' or bacalhau_state = 'Complete') and contract_state = 'None'"
	}
	sqlStatement := fmt.Sprintf(`
select
	id, created, bacalhau_inference_id, bacalhau_state, contract_state, artist_id, prompt, error
from
	image
%s
order by
	created desc
`, where)

	rows, err := d.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := []*types.Image{}
	for rows.Next() {
		image, err := scanImage(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, image)
	}
	if err = rows.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

func (d *SQLiteStore) GetImage(ctx context.Context, id int) (*types.Image, error) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	row := d.db.QueryRow(`
select
	id, created, bacalhau_inference_id, bacalhau_state, contract_state, artist_id, prompt, error
from
	image
where
	id = $1
`, id)
	image, err := scanImage(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image not found: %d %s", id, err.Error())
		} else {
			return nil, err
		}
	}
	return image, nil
}

func (d *SQLiteStore) AddImage(ctx context.Context, data types.Image) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	sqlStatement := `
insert into image (id, artist_id, prompt)
values ($1, $2, $3)`
	_, err := d.db.Exec(
		sqlStatement,
		data.ID,
		data.Artist,
		data.Prompt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQLiteStore) UpdateImage(ctx context.Context, data types.Image) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	sqlStatement := `
update image set
	bacalhau_inference_id = $2,
	bacalhau_state = $3,
	contract_state = $4,
	error = $5
 where id = $1`
	_, err := d.db.Exec(
		sqlStatement,
		data.ID,
		data.BacalhauInferenceID,
		data.BacalhauState.String(),
		data.ContractState.String(),
		data.Error,
	)
	if err != nil {
		return err
	}
	return nil
}

//go:embed migrations/*.sql
var fs embed.FS

func (d *SQLiteStore) GetMigrations() (*migrate.Migrate, error) {
	files, err := iofs.New(fs, "migrations")
	if err != nil {
		return nil, err
	}
	migrations, err := migrate.NewWithSourceInstance(
		"iofs",
		files,
		fmt.Sprintf("sqlite://%s", d.options.DataFile),
	)
	if err != nil {
		return nil, err
	}
	return migrations, nil
}

func (d *SQLiteStore) MigrateUp() error {
	migrations, err := d.GetMigrations()
	if err != nil {
		return err
	}
	err = migrations.Up()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (d *SQLiteStore) MigrateDown() error {
	migrations, err := d.GetMigrations()
	if err != nil {
		return err
	}
	err = migrations.Down()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}
