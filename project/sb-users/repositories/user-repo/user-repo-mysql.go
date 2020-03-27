package userrepo

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MartinNikolovMarinov/sb-users/models"
	"golang.org/x/crypto/bcrypt"
)

type userRepoMysql struct {
	db *sql.DB
}

// NewMysqlUserRepo is a UserRepo constructor
func NewMysqlUserRepo(connectionString string) UserRepo {
	repo := &userRepoMysql{}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func (u *userRepoMysql) Find(name string) (*models.User, error) {
	const findUserByNameQuery = `
		select u.ID, u.Name, u.Password, u.Roles from User as u
		where u.IsDeleted = 0 and u.Name = ?
	`

	var password string
	var id int64
	var roles int

	err := u.db.
		QueryRow(findUserByNameQuery, name).
		Scan(&id, &name, &password, &roles)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var user = &models.User{ID: id, Name: name, Password: password, Roles: models.RoleFlags(roles)}
	return user, nil
}

func (u *userRepoMysql) MatchPassword(password, hashedPass string) error {
	var err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	return err
}

func (u *userRepoMysql) Create(
	name string,
	password string,
	roles models.RoleFlags,
	siteID int64) (*models.User, error) {

	const createUserQuery = `
		insert into User(name, password, roles) values (?, ?, ?)
	`
	const insertInSiteToUserQuery = `
		insert into Site2User(SiteID, UserID) values (?, ?)
	`

	// Hash the pasword with bcrypt
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	password = string(pass)

	result, err := u.db.Exec(createUserQuery, name, password, int(roles))
	if err != nil {
		return nil, err
	}

	var user = &models.User{Name: name, Password: password, Roles: roles}
	user.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	result, err = u.db.Exec(insertInSiteToUserQuery, siteID, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepoMysql) AllPerSite(siteID int64) ([]models.User, error) {
	const allUsersPerSiteQuery = `
		select  u.ID, u.Name, u.Roles from User as u
		join Site as s on s.ID = ?
		where u.IsDeleted = 0 and s.IsDeleted = 0
	`
	users := make([]models.User, 0)
	rows, err := u.db.Query(allUsersPerSiteQuery, siteID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var name string
		var roles int64
		err := rows.Scan(&id, &name, &roles)
		if err != nil {
			return nil, err
		}

		var s = models.User{ID: id, Name: name, Roles: models.RoleFlags(roles), Password: ""}
		users = append(users, s)
	}

	return users, nil
}

// Delete commnet
func (u *userRepoMysql) Delete(userID int64) error {
	const deleteUserQuery = `
		update User
		set IsDeleted = 1
		where ID = ?
	`

	_, err := u.db.Exec(deleteUserQuery, userID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUsername commnet
func (u *userRepoMysql) UpdateUsername(userID int64, userName string) error {
	const updateUsernameQuery = `
		update User
		set Name = ?
		where ID = ?
	`
	const doesUsernameExistQuery = `
		select ID from User where Name = ?
	`

	var id int64
	row := u.db.QueryRow(doesUsernameExistQuery, userName)
	err := row.Scan(&id)
	if err != sql.ErrNoRows {
		return errors.New("Username already exists")
	}

	result, err := u.db.Exec(updateUsernameQuery, userName, userID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil || n != 1 {
		return errors.New("User not found")
	}

	return nil
}

// UpdatePassword commnet
func (u *userRepoMysql) UpdatePassword(userID int64, password string) error {
	const updateUsernameQuery = `
		update User
		set Password = ?
		where ID = ?
	`

	// Hash the pasword with bcrypt
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result, err := u.db.Exec(updateUsernameQuery, string(hashedPass), userID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil || n != 1 {
		return errors.New("User not found")
	}

	return nil
}
