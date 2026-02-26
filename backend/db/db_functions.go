package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Storage interface {
	Init() error
	GetUserWithEmail(string) *User
	GetUserWithAuthToken(string) *User
	SetNewAuthToken(string, string, string, string)
	SetClientAuthToken(string, string, string)
	RemoveAuthToken(string)
	CreateNewUser(string, string, string, string)
	AddNewPassord(int, string, string) error
	GetHostNames(int) []string
	GetPassword(int, string) (string, error)
	RemovePassword(int, string) error
	Migrate()
	EditPassword(int, string, string) error
}

type LocalStorage struct {
	db *sql.DB
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (s *LocalStorage) Init() error {
	database, err := sql.Open("sqlite3", "./db/.db")
	if err != nil {
		log.Fatal("Could not connect to database")
		return err
	}
	log.Println("Connected to './db/.db'")
	s.db = database
	return nil
}

func (s *LocalStorage) GetUserWithEmail(userEmail string) *User {
	row, err := s.db.Query(fmt.Sprintf("SELECT * FROM user WHERE email = '%s';", userEmail))
	if err != nil {
		log.Fatalf("Could not retrieve user with email '%s'", userEmail)
		return nil
	}

	defer func() {
		if err := row.Close(); err != nil {
			log.Fatal("Could not close database row")
		}
	}()

	user := DbEntryToUser(row)
	if user == nil {
		log.Println("User not found(GetUserWithEmail)")
	}
	return user
}

func (s *LocalStorage) GetUserWithAuthToken(authToken string) *User {
	row, err := s.db.Query(fmt.Sprintf("SELECT * FROM user where auth_token = '%s';", authToken))
	if err != nil {
		log.Fatalf("Could not retrieve user with auth_token '%s'", authToken)
		return nil
	}

	defer func() {
		if err := row.Close(); err != nil {
			log.Fatal("Could not close database row")
		}
	}()

	selectedUser := DbEntryToUser(row)
	if selectedUser == nil {
		log.Println("User not found(GetUserWithAuthToken)")
		return nil
	}

	if selectedUser.TokenExpiryDate <= time.Now().Unix() {
		log.Println("auth_token not valid")
		s.RemoveAuthToken(selectedUser.Email)
		return nil
	}

	return selectedUser
}

func (s *LocalStorage) SetNewAuthToken(userEmail, userToken, clientToken, ipAddr string) {
	expiryDate := time.Now().AddDate(0, 1, 0).Unix()
	updateUserQuery := fmt.Sprintf(`UPDATE user
    SET auth_token = '%s'
    WHERE 
    email = '%s';`, userToken, userEmail)
	userStmt, err := s.db.Prepare(updateUserQuery)
	if err != nil {
		log.Fatalf("Could not update AUTH on '%s'(SetAuthToken)", userEmail)
		return
	}
	userStmt.Exec()
	clientAuthQuery := `INSERT INTO auth_key(
		ip_addr,
		user_token,
		client_token,
		token_expiry_date
	) VALUES (
		?,?,?,?
	);`
	tokenStmt, err := s.db.Prepare(clientAuthQuery)
	if err != nil {
		log.Fatalf("Could not add CLIENT TOKEN on '%s', '%s'", userEmail, ipAddr)
		return
	}
	tokenStmt.Exec(ipAddr, userToken, clientToken, expiryDate)

	log.Printf("Updated AUTH on '%s'", userEmail)
}

func (s *LocalStorage) SetClientAuthToken(userEmail, clientToken, ipAddr string) {

}

func (s *LocalStorage) RemoveAuthToken(userEmail string) {
	user := s.GetUserWithEmail(userEmail)
	if user == nil {
		log.Fatal("User not found(RemoveAuthToken)")
		return
	}
	removeTokenQuery := fmt.Sprintf(`UPDATE user
    SET auth_token = '',
        token_expiry_date = 0
    WHERE
    email = '%s';`, userEmail)

	statement, err := s.db.Prepare(removeTokenQuery)
	if err != nil {
		log.Fatalf("Could not remove AUTH token on '%s'", userEmail)
		return
	}
	statement.Exec()
	log.Printf("Removed AUTH token on '%s'", userEmail)
}

func (s *LocalStorage) CreateNewUser(email, password, name, surname string) {
	createNewUserQuery := `INSERT INTO user(
        email, 
        password, 
        name, 
        surname,
        auth_token) 
        VALUES(?,?,?,?,?);`
	log.Println("Inserting new user")
	statement, err := s.db.Prepare(createNewUserQuery)
	if err != nil {
		log.Fatalf("Error inserting new user '%s'", email)
		return
	}
	statement.Exec(email, password, name, surname, "")
	log.Printf("Created user '%s %s' - '%s'", name, surname, email)
}

func (s *LocalStorage) AddNewPassord(userId int, password, hostName string) error {
	insertNewPasswordQuery := `INSERT INTO password(
		user_id,
		password,
		host_name)
		VALUES(?,?,?);`
	log.Printf("Inserting new password for user '%d'", userId)
	statement, err := s.db.Prepare(insertNewPasswordQuery)
	if err != nil {
		log.Fatalf("Error inserting new password for user '%d' - %v", userId, err)
		return err
	}

	_, err = statement.Exec(userId, password, hostName)
	if err != nil {
		return err
	}
	log.Printf("Inserted new password for user '%d'", userId)
	return nil
}

func (s *LocalStorage) GetHostNames(userId int) []string {
	getHostsQuery := fmt.Sprintf("SELECT host_name FROM password where user_id = %d", userId)

	row, err := s.db.Query(getHostsQuery)
	if err != nil {
		log.Printf("Could not get host names for user: %d - %v", userId, err)
		return []string{}
	}

	return DbEntryToHostNames(row)
}

func (s *LocalStorage) GetPassword(userId int, hostname string) (string, error) {
	getPasswordQuery := fmt.Sprintf(`SELECT password 
		FROM password 
		WHERE 
			user_id = %d
		AND 
			host_name = '%s';`, userId, hostname)

	row, err := s.db.Query(getPasswordQuery)
	if err != nil {
		log.Fatalf("Could not get host names for user: %d - %v", userId, err)
		return "", err
	}

	return DbEntryToPassword(row), nil
}

func (s *LocalStorage) RemovePassword(userId int, hostname string) error {
	removePasswordQuery := `DELETE FROM password
		WHERE
			user_id = ?
		AND
			host_name = ?;`
	statement, err := s.db.Prepare(removePasswordQuery)
	if err != nil {
		log.Fatalf("Could not remove password (userId: %d - hostname: %s). Error: %v", userId, hostname, err)
		return err
	}

	_, err = statement.Exec(userId, hostname)
	return err
}

func (s *LocalStorage) EditPassword(userId int, hostname, newPassword string) error {
	editPasswordQuery := `UPDATE password
		SET password = ?
		WHERE user_id = ? AND host_name = ?;`
	log.Printf("Updating password for user '%d'", userId)
	statement, err := s.db.Prepare(editPasswordQuery)
	if err != nil {
		log.Fatalf("Could not update password ( userId: %d - hostname: %s). Error: %v", userId, hostname, err)
		return err
	}
	_, err = statement.Exec(newPassword, userId, hostname)
	return err
}
