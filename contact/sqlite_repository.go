package contact

import (
    "database/sql"
    "errors"

    "github.com/mattn/go-sqlite3"
)

type Repository interface {
	Migrate() error
	Create(contact Contact) (*Contact, error)
	All() ([]Contact, error)
	SearchByName(name string) ([]Contact, error) 
	GetById(name string) (*Contact, error) 
	Update(id int64, updated Contact) (*Contact, error)
	Delete(id int64) error
}

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("record does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
  CREATE TABLE IF NOT EXISTS contacts(
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	`

	_, err := r.db.Exec(query)

	return err
}

func (r *SQLiteRepository) Create(contact Contact) (*Contact, error) {
	res, err := r.db.Exec("INSERT INTO contacts(name) values(?", contact.Name)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	contact.Id = id

	return &contact, nil
}

func (r *SQLiteRepository) All() ([]Contact, error) {
	rows, err := r.db.Query("SELECT * FROM contacts")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Contact

	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.Id, &contact.Name); err != nil {
			return nil, err
		}

		all = append(all, contact)
	}

	return all, nil
}

func (r *SQLiteRepository) SearchByName(name string) ([]Contact, error) {
	rows, err := r.db.Query("SELECT * FROM contacts WHERE name = ?", name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Contact

	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.Id, &contact.Name); err != nil {
			return nil, err
		}

		results = append(results, contact)
	}

	return results, nil
}

func (r *SQLiteRepository) GetById(id int64) (*Contact, error) {
	row := r.db.QueryRow("SELECT * FROM contacts WHERE id = ?", id)

	var contact Contact
	
	if err := row.Scan(&contact.Id, &contact.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}

		return nil, err
	}

	return &contact, nil
}
