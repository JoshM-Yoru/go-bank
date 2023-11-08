package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User, *Account) error
	DeleteAccount(int) error
	UpdateUser(*User) error
	GetUsers() ([]*User, error)
	GetAccounts() ([]*Account, error)
	GetUserByID(int) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetAccountByUserID(int) (*FullAccount, error)
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
	roleTable := s.CreateRoleTable()
	if roleTable != nil {
		return roleTable
	}
	userTable := s.CreateUserTable()
	if userTable != nil {
		return userTable
	}
	accTable := s.CreateAccountTable()
	if accTable != nil {
		return accTable
	}

	transferTable := s.CreateTransferTable()
	if transferTable != nil {
		return transferTable
	}

	transactionTable := s.CreateTransactionTable()
	if transactionTable != nil {
		return transactionTable
	}

	return nil
}

func (s *PostgresStore) CreateUserTable() error {
	//need to add role reference
	query := `create table if not exists user_profile (
        user_id serial primary key,
        email varchar(50) unique not null,
        password varchar(100) not null,
        first_name varchar(50) not null,
        last_name varchar(50) not null,
        user_name varchar(50) unique not null,
        phone_number varchar(10),
        created_at timestamp,
        last_login timestamp,
        role int,
        is_active_user boolean
    )`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccountTable() error {
	//need to add fk_account_type reference
	query := `create table if not exists account (
        account_id serial primary key,
        fk_user serial references user_profile(user_id), 
        account_number serial unique, 
        balance bigint,
        created_at timestamp,
        fk_account_type int,
        is_active_account boolean
    )`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateRoleTable() error {
	query := `create table if not exists role (
        role_id serial primary key,
        role_name varchar(10)
    )`

	_, err := s.db.Query(query)

	return err
}

func (s *PostgresStore) CreateTransferTable() error {
	query := `create table if not exists transfer (
        id serial primary key,
        from_account int references account(account_id) not null,
        to_account int references account(account_id) not null,
        amount bigint not null,
        description varchar(200),
        created_at timestamp
   )`

	_, err := s.db.Query(query)

	return err
}

func (s *PostgresStore) CreateTransactionTable() error {
	query := `create table if not exists transaction (
        id serial primary key,
        account int references account(account_id) not null,
        transaction_source varchar(50) not null,
        amount bigint not null,
        description varchar(200),
        created_at timestamp,
        transaction_type int
    )`

	_, err := s.db.Query(query)

	return err
}

func (s *PostgresStore) CreateUser(user *User, account *Account) error {

	if user.Role == Admin {
		query := `insert into user_profile (email, password, first_name, last_name, user_name, phone_number, created_at, last_login, role, is_active_user) 
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		_, err := s.db.Query(
			query,
			user.Email,
			user.Password,
			user.FirstName,
			user.LastName,
			user.UserName,
			user.PhoneNumber,
			user.CreatedAt,
			user.LastLogin,
			user.Role,
			user.IsActive,
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				return fmt.Errorf("Email in use")
			}
			return err
		}
	} else {
		query := `with x as (
            insert into user_profile (email, password, first_name, last_name, user_name, phone_number, created_at, last_login, role, is_active_user) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            returning user_id
        )
    insert into account (fk_user, account_number, balance, created_at, fk_account_type, is_active_account)
    select x.user_id, $11, $12, $13, $14, $15
    from x;
    `

		_, err := s.db.Query(
			query,
			user.Email,
			user.Password,
			user.FirstName,
			user.LastName,
			user.UserName,
			user.PhoneNumber,
			user.CreatedAt,
			user.LastLogin,
			user.Role,
			user.IsActive,
			account.AccountNumber,
			account.Balance,
			account.CreatedAt,
			account.AccountType,
			account.IsActiveAccount,
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "user") {
				return fmt.Errorf("Email in use")
			} else if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "account") {
				return fmt.Errorf("Something went wrong, please try registering again")
            }
			return err
		}
	}

	return nil
}

func (s *PostgresStore) UpdateUser(*User) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, accErr := s.db.Query(`update account
    set is_active_account = false
    where fk_user = $1`, id)
	if accErr != nil {
		return accErr
	}

	_, err := s.db.Query(`update user_profile 
        set is_active_user = false 
        where user_id = $1`, id)
	return err
}

func (s *PostgresStore) GetUsers() ([]*User, error) {
	rows, err := s.db.Query("select * from user where is_active = true")
	if err != nil {
		return nil, err
	}

	users := []*User{}

	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
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

func (s *PostgresStore) GetUserByID(id int) (*User, error) {
	rows, err := s.db.Query("select * from user where is_active = true and id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("Account %d not found", id)
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	rows, err := s.db.Query("select * from user where is_active = true email = $1", email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("User %v not found", email)
}

func (s *PostgresStore) GetAccountByUserID(id int) (*FullAccount, error) {
	return nil, nil
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.LastLogin,
		&user.Role,
		&user.IsActive,
	)

	return user, err
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.AccountNumber,
		&account.Balance,
		&account.CreatedAt,
		&account.AccountType,
		&account.IsActiveAccount,
	)

	return account, err
}

func (s *PostgresStore) CreateTransfer() {}
