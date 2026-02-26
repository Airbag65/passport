package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type Storage interface {
	Init() error
	GetUserWithEmail(string) *User
	GetUserWithAuthToken(string) *User
	ValidateToken(string, string, string) bool
	SetNewAuthToken(string, string, string, string)
	SetClientAuthToken(string, string, string) (string, error)
	RemoveAuthToken(string, string)
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
	splitToken := strings.Split(authToken, "+")
	userToken := splitToken[0]
	// clientToken := splitToken[1]
	row, err := s.db.Query(fmt.Sprintf("SELECT * FROM user where auth_token = '%s';", userToken))
	if err != nil {
		log.Fatalf("Could not retrieve user with auth_token '%s'", userToken)
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

	return selectedUser
}

func (s *LocalStorage) ValidateToken(authToken, ipAddr, userEmail string) bool {
	tokens := strings.Split(authToken, "+")

	row, err := s.db.Query(fmt.Sprintf(`SELECT * FROM auth_key
		WHERE
			user_token = '%s'
		AND
			ip_addr = '%s'
		AND
			client_token = '%s';`, tokens[0], ipAddr, tokens[1]))
	if err != nil {
		log.Fatalf("Could not retrieve auth_key with auth_token '%s'", tokens[0])
		return false
	}
	authKey := DbEntryToAuthKey(row)
	if authKey == nil {
		log.Println("AuthKey not found")
		return false
	}

	if authKey.TokenExpiryDate <= time.Now().Unix() {
		log.Println("auth_token not valid")
		s.RemoveAuthToken(userEmail, ipAddr)
		return false
	}
	return true
}

func (s *LocalStorage) SetNewAuthToken(userEmail, userToken, clientToken, ipAddr string) {
	expiryDate := time.Now().AddDate(0, 1, 0).Unix()
	updateUserQuery := fmt.Sprintf(`UPDATE user
    SET auth_token = '%s',
    	logged_in_count = 1
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

func (s *LocalStorage) SetClientAuthToken(userEmail, clientToken, ipAddr string) (string, error) {
	expiryDate := time.Now().AddDate(0, 1, 0).Unix()
	user := s.GetUserWithEmail(userEmail)
	if user == nil {
		log.Fatalf("Could not get user with email '%s'", userEmail)
	}
	userQuery := `UPDATE user
	SET logged_in_count = ?
	WHERE email = ?;`
	userStmt, err := s.db.Prepare(userQuery)
	if err != nil {
		log.Fatalf("Could not update logged_in_count for user '%s'", userEmail)
		return "", nil
	}
	userStmt.Exec(user.LoggedInCount+1, userEmail)
	clientQuery := `INSERT INTO auth_key(
		ip_addr,
		user_token,
		client_token,
		token_expiry_date
	) VALUES (
		?,?,?,?		
	);`
	clientStmt, err := s.db.Prepare(clientQuery)
	if err != nil {
		log.Fatalf("Could not add client token for '%s' '%s'", userEmail, ipAddr)
		return "", err
	}
	clientStmt.Exec(ipAddr, user.AuthToken, clientToken, expiryDate)
	return user.AuthToken + "+" + clientToken, nil
}

func (s *LocalStorage) RemoveAuthToken(userEmail, ipAddr string) {
	user := s.GetUserWithEmail(userEmail)
	if user == nil {
		log.Fatal("User not found(RemoveAuthToken)")
		return
	}
	removeClientTokenQuery := `DELETE FROM auth_key WHERE
		ip_addr = ?
	AND
		user_token = ?;`
	clientStmt, err := s.db.Prepare(removeClientTokenQuery)
	if err != nil {
		log.Fatalf("Could not remove Client Auth Token on %s", ipAddr)
		return
	}
	clientStmt.Exec(ipAddr, user.AuthToken)
	userQuery := ""
	if user.LoggedInCount-1 == 0 {
		userQuery += `UPDATE user
		SET auth_token = '',
			logged_in_count = 0
		WHERE
			email = ?;`
	} else {
		userQuery += fmt.Sprintf(`UPDATE user
		SET logged_in_count = %d
		WHERE email = ?;`, user.LoggedInCount-1)
	}
	userStmt, err := s.db.Prepare(userQuery)
	if err != nil {
		log.Fatalf("Could not remove AUTH token on '%s'", userEmail)
		return
	}
	userStmt.Exec(userEmail)
	log.Printf("Removed AUTH token on '%s'", userEmail)
}

func (s *LocalStorage) CreateNewUser(email, password, name, surname string) {
	createNewUserQuery := `INSERT INTO user(
        email, 
        password, 
        name, 
        surname,
        auth_token,
        logged_in_count) 
        VALUES(?,?,?,?,?,?);`
	log.Println("Inserting new user")
	statement, err := s.db.Prepare(createNewUserQuery)
	if err != nil {
		log.Fatalf("Error inserting new user '%s'", email)
		return
	}
	statement.Exec(email, password, name, surname, "", 0)
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
