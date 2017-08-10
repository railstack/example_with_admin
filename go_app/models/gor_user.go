// Package models includes the functions on the model User.
package models

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

// set flags to output more detailed log
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type User struct {
	Id                  int64     `json:"id,omitempty" db:"id" valid:"-"`
	Email               string    `json:"email,omitempty" db:"email" valid:"required,matches(\A[^@\s]+@[^@\s]+\z)"`
	EncryptedPassword   string    `json:"encrypted_password,omitempty" db:"encrypted_password" valid:"-"`
	ResetPasswordToken  string    `json:"reset_password_token,omitempty" db:"reset_password_token" valid:"-"`
	ResetPasswordSentAt time.Time `json:"reset_password_sent_at,omitempty" db:"reset_password_sent_at" valid:"-"`
	RememberCreatedAt   time.Time `json:"remember_created_at,omitempty" db:"remember_created_at" valid:"-"`
	SignInCount         int64     `json:"sign_in_count,omitempty" db:"sign_in_count" valid:"-"`
	CurrentSignInAt     time.Time `json:"current_sign_in_at,omitempty" db:"current_sign_in_at" valid:"-"`
	LastSignInAt        time.Time `json:"last_sign_in_at,omitempty" db:"last_sign_in_at" valid:"-"`
	CurrentSignInIp     string    `json:"current_sign_in_ip,omitempty" db:"current_sign_in_ip" valid:"-"`
	LastSignInIp        string    `json:"last_sign_in_ip,omitempty" db:"last_sign_in_ip" valid:"-"`
	CreatedAt           time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
	UpdatedAt           time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
	Role                string    `json:"role,omitempty" db:"role" valid:"-"`
	Posts               []Post    `json:"posts,omitempty" db:"posts" valid:"-"`
}

// FindUser find a single user by an id
func FindUser(id int64) (*User, error) {
	if id == 0 {
		return nil, errors.New("Invalid id: it can't be zero")
	}
	_user := User{}
	err := db.Get(&_user, db.Rebind(`SELECT * FROM users WHERE id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_user, nil
}

// FirstUser find the first one user by id ASC order
func FirstUser() (*User, error) {
	_user := User{}
	err := db.Get(&_user, db.Rebind(`SELECT * FROM users ORDER BY id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_user, nil
}

// FirstUsers find the first N users by id ASC order
func FirstUsers(n uint32) ([]User, error) {
	_users := []User{}
	sql := fmt.Sprintf("SELECT * FROM users ORDER BY id ASC LIMIT %v", n)
	err := db.Select(&_users, db.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _users, nil
}

// LastUser find the last one user by id DESC order
func LastUser() (*User, error) {
	_user := User{}
	err := db.Get(&_user, db.Rebind(`SELECT * FROM users ORDER BY id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_user, nil
}

// LastUsers find the last N users by id DESC order
func LastUsers(n uint32) ([]User, error) {
	_users := []User{}
	sql := fmt.Sprintf("SELECT * FROM users ORDER BY id DESC LIMIT %v", n)
	err := db.Select(&_users, db.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _users, nil
}

// FindUsers find one or more users by one or more ids
func FindUsers(ids ...int64) ([]User, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_users := []User{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := db.Rebind(fmt.Sprintf(`SELECT * FROM users WHERE id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := db.Select(&_users, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _users, nil
}

// FindUserBy find a single user by a field name and a value
func FindUserBy(field string, val interface{}) (*User, error) {
	_user := User{}
	sqlFmt := `SELECT * FROM users WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := db.Get(&_user, db.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_user, nil
}

// FindUsersBy find all users by a field name and a value
func FindUsersBy(field string, val interface{}) (_users []User, err error) {
	sqlFmt := `SELECT * FROM users WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = db.Select(&_users, db.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _users, nil
}

// AllUsers get all the User records
func AllUsers() (users []User, err error) {
	err = db.Select(&users, "SELECT * FROM users")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

// UserCount get the count of all the User records
func UserCount() (c int64, err error) {
	err = db.Get(&c, "SELECT count(*) FROM users")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// UserCountWhere get the count of all the User records with a where clause
func UserCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM users"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	err = stmt.Get(&c, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// UserIncludesWhere get the User associated models records, it's just the eager_load function
func UserIncludesWhere(assocs []string, sql string, args ...interface{}) (_users []User, err error) {
	_users, err = FindUsersWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _users, err
	}
	if len(_users) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_users))
	for _, v := range _users {
		ids = append(ids, interface{}(v.Id))
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	for _, assoc := range assocs {
		switch assoc {
		case "posts":
			where := fmt.Sprintf("user_id IN (?%s)", idsHolder)
			_posts, err := FindPostsWhere(where, ids...)
			if err != nil {
				log.Printf("Error when query associated objects: %v\n", assoc)
				continue
			}
			for _, vv := range _posts {
				for i, vvv := range _users {
					if vv.UserId == vvv.Id {
						vvv.Posts = append(vvv.Posts, vv)
					}
					_users[i].Posts = vvv.Posts
				}
			}
		}
	}
	return _users, nil
}

// UserIds get all the Ids of User records
func UserIds() (ids []int64, err error) {
	err = db.Select(&ids, "SELECT id FROM users")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// UserIdsWhere get all the Ids of User records by where restriction
func UserIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := UserIntCol("id", where, args...)
	return ids, err
}

// UserIntCol get some int64 typed column of User by where restriction
func UserIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM users"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&intColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return intColRecs, nil
}

// UserStrCol get some string typed column of User by where restriction
func UserStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM users"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&strColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strColRecs, nil
}

// FindUsersWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindUsersWhere(where string, args ...interface{}) (users []User, err error) {
	sql := "SELECT * FROM users"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&users, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

// FindUserBySql query use a complete SQL clause
// with placeholders, eg: FindUserBySql("SELECT * FROM users WHERE first_name = ? AND age > ? ORDER BY DESC LIMIT 1", "John", 18)
// will return only One record in the table "users" whose first_name is "John" and age elder than 18
func FindUserBySql(sql string, args ...interface{}) (*User, error) {
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_user := &User{}
	err = stmt.Get(_user, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return _user, nil
}

// FindUsersBySql query use a complete SQL clause
// with placeholders, eg: FindUsersBySql("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindUsersBySql(sql string, args ...interface{}) (users []User, err error) {
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&users, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

// CreateUser use a named params to create a single User record.
// A named params is key-value map like map[string]interface{}{"first_name": "John", "age": 23} .
func CreateUser(am map[string]interface{}) (int64, error) {
	if len(am) == 0 {
		return 0, fmt.Errorf("Zero key in the attributes map!")
	}
	t := time.Now()
	for _, v := range []string{"created_at", "updated_at"} {
		if am[v] == nil {
			am[v] = t
		}
	}
	keys := make([]string, len(am))
	i := 0
	for k := range am {
		keys[i] = k
		i++
	}
	sqlFmt := `INSERT INTO users (%s) VALUES (%s)`
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(keys, ","), ":"+strings.Join(keys, ",:"))
	result, err := db.NamedExec(sqlStr, am)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// Create is a method for User to create a record
func (_user *User) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_user)
	if !ok {
		errMsg := "Validate User struct error: Unknown error"
		if err != nil {
			errMsg = "Validate User struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_user.CreatedAt = t
	_user.UpdatedAt = t
	sql := `INSERT INTO users (email,encrypted_password,reset_password_token,reset_password_sent_at,remember_created_at,sign_in_count,current_sign_in_at,last_sign_in_at,current_sign_in_ip,last_sign_in_ip,created_at,updated_at,role) VALUES (:email,:encrypted_password,:reset_password_token,:reset_password_sent_at,:remember_created_at,:sign_in_count,:current_sign_in_at,:last_sign_in_at,:current_sign_in_ip,:last_sign_in_ip,:created_at,:updated_at,:role)`
	result, err := db.NamedExec(sql, _user)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// PostsCreate is used for User to create the associated objects Posts
func (_user *User) PostsCreate(am map[string]interface{}) error {
	am["user_id"] = _user.Id
	_, err := CreatePost(am)
	return err
}

// GetPosts is used for User to get associated objects Posts
// Say you have a User object named user, when you call user.GetPosts(),
// the object will get the associated Posts attributes evaluated in the struct
func (_user *User) GetPosts() error {
	_posts, err := UserGetPosts(_user.Id)
	if err == nil {
		_user.Posts = _posts
	}
	return err
}

// UserGetPosts a helper fuction used to get associated objects for UserIncludesWhere()
func UserGetPosts(id int64) ([]Post, error) {
	_posts, err := FindPostsBy("user_id", id)
	return _posts, err
}

// Destroy is method used for a User object to be destroyed.
func (_user *User) Destroy() error {
	if _user.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyUser(_user.Id)
	return err
}

// DestroyUser will destroy a User record specified by id parameter.
func DestroyUser(id int64) error {
	stmt, err := db.Preparex(db.Rebind(`DELETE FROM users WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// DestroyUsers will destroy User records those specified by the ids parameters.
func DestroyUsers(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM users WHERE id IN (?%s)`, idsHolder)
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	result, err := stmt.Exec(idsT...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// DestroyUsersWhere delete records by a where clause
// like: DestroyUsersWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyUsersWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM users WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// Save method is used for a User object to update an existed record mainly.
// If no id provided a new record will be created. A UPSERT action will be implemented further.
func (_user *User) Save() error {
	ok, err := govalidator.ValidateStruct(_user)
	if !ok {
		errMsg := "Validate User struct error: Unknown error"
		if err != nil {
			errMsg = "Validate User struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _user.Id == 0 {
		_, err = _user.Create()
		return err
	}
	_user.UpdatedAt = time.Now()
	sqlFmt := `UPDATE users SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "email = :email, encrypted_password = :encrypted_password, reset_password_token = :reset_password_token, reset_password_sent_at = :reset_password_sent_at, remember_created_at = :remember_created_at, sign_in_count = :sign_in_count, current_sign_in_at = :current_sign_in_at, last_sign_in_at = :last_sign_in_at, current_sign_in_ip = :current_sign_in_ip, last_sign_in_ip = :last_sign_in_ip, updated_at = :updated_at, role = :role", _user.Id)
	_, err = db.NamedExec(sqlStr, _user)
	return err
}

// UpdateUser is used to update a record with a id and map[string]interface{} typed key-value parameters.
func UpdateUser(id int64, am map[string]interface{}) error {
	if len(am) == 0 {
		return errors.New("Zero key in the attributes map!")
	}
	am["updated_at"] = time.Now()
	keys := make([]string, len(am))
	i := 0
	for k := range am {
		keys[i] = k
		i++
	}
	sqlFmt := `UPDATE users SET %s WHERE id = %v`
	setKeysArr := []string{}
	for _, v := range keys {
		s := fmt.Sprintf(" %s = :%s", v, v)
		setKeysArr = append(setKeysArr, s)
	}
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(setKeysArr, ", "), id)
	_, err := db.NamedExec(sqlStr, am)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Update is a method used to update a User record with the map[string]interface{} typed key-value parameters.
func (_user *User) Update(am map[string]interface{}) error {
	if _user.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateUser(_user.Id, am)
	return err
}

func (_user *User) UpdateAttributes(am map[string]interface{}) error {
	if _user.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateUser(_user.Id, am)
	return err
}

// UpdateColumns method is supposed to be used to update User records as corresponding update_columns in Rails
func (_user *User) UpdateColumns(am map[string]interface{}) error {
	if _user.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateUser(_user.Id, am)
	return err
}

// UpdateUsersBySql is used to update User records by a SQL clause
// that use '?' binding syntax
func UpdateUsersBySql(sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return 0, errors.New("A blank SQL clause")
	}
	sql = strings.Replace(strings.ToLower(sql), "set", "set updated_at = ?, ", 1)
	args = append([]interface{}{time.Now()}, args...)
	stmt, err := db.Preparex(db.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
