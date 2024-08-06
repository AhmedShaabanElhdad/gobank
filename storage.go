package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// need to think of way to better solution as interface increase in size
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int64) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgressStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgressStore, error) {
	// connStr := "user=pqgotest dbname=gobankdb password=gobank sslmode=verify-full"
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// defer db.Close()
	return &PostgressStore{
		db: db,
	}, nil
}

func (s *PostgressStore) init() error {
	return s.createAccountTable()
}

func (s *PostgressStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		encryptedPass varchar(100),
		number serial,
		balance serial,
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgressStore) CreateAccount(account *Account) error {
	query := `insert into account 
	(
		first_name, last_name, number, balance, created_at, encryptedPass
	) values($1, $2, $3, $4, $5, $6)`
	sqlRes, err := s.db.Query(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt, account.EncryptedPassword)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", sqlRes)
	return nil
}

func (s *PostgressStore) DeleteAccount(id int) error {
	// in production we don't delete account we make it with flag false or make it inActive
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *PostgressStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgressStore) GetAccountByID(number int) (*Account, error) {
	rows, err := s.db.Query("select * from account where number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account number  %d  not found", number)
}

func (s *PostgressStore) GetAccountByNumber(number int64) (*Account, error) {
	rows, err := s.db.Query("select * from account where number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d  not found", number)
}

func (s *PostgressStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.EncryptedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	return account, err
}
