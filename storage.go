package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	conStr := "user=postgres dbname=gobank password=postgres sslmode=disable"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	accTable := s.CreateAccountTable()
    if accTable !=nil {
        return accTable
    }

	transferTable := s.CreateTransferTable()
    if transferTable !=nil {
        return transferTable
    }

	transactionTable := s.CreateTransactionTable()
    if transactionTable !=nil {
        return transactionTable
    }

	return nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
        id serial primary key,
        email varchar(50) unique not null,
        password varchar(100) not null,
        first_name varchar(50) not null,
        last_name varchar(50) not null,
        user_name varchar(50) not null,
        phone_number varchar(11),
        account_number serial, 
        balance bigint not null,
        created_at timestamp,
        role int,
        is_active boolean
    )`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `insert into account 
    (email, password, first_name, last_name, user_name, phone_number, account_number, balance, created_at, role, is_active) 
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`


	_, err := s.db.Query(
		query,
		account.Email,
		account.Password,
		account.FirstName,
		account.LastName,
        account.UserName,
		account.PhoneNumber,
		account.AccountNumber,
		account.Balance,
		account.CreatedAt,
		account.Role,
        account.IsActive,
	)
	if err != nil {
        if strings.Contains(err.Error(), "duplicate") {
            return fmt.Errorf("Email in use")
        }
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("update account set is_active = false where id = $1", id)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account where is_active = true")
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

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where is_active = true and id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found", id)
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	rows, err := s.db.Query("select * from account where is_active = true email = $1", email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %v not found", email)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.Email,
		&account.Password,
		&account.FirstName,
		&account.LastName,
		&account.UserName,
		&account.PhoneNumber,
		&account.AccountNumber,
		&account.Balance,
		&account.CreatedAt,
		&account.Role,
		&account.IsActive,
	)

	return account, err
}


func (s *PostgresStore) CreateTransferTable() error {
	query := `create table if not exists transfer (
        id serial primary key,
        from_account int references account(id) not null,
        to_account int references account(id) not null,
        amount bigint not null,
        description varchar(200),
        created_at timestamp
   )`

	_, err := s.db.Query(query)

	return err
}

func (s *PostgresStore) CreateTransfer() {}

func (s *PostgresStore) CreateTransactionTable() error {
	query := `create table if not exists transaction (
        id serial primary key,
        account int references account(id) not null,
        transaction_source varchar(50) not null,
        amount bigint not null,
        description varchar(200),
        created_at timestamp,
        transaction_type int
    )`

	_, err := s.db.Query(query)

	return err
}
