// Package models includes the functions on the model Post.
package models

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

// set flags to output more detailed log
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Post struct {
	Id        int64     `json:"id,omitempty" db:"id" valid:"-"`
	Title     string    `json:"title,omitempty" db:"title" valid:"required,length(10|50)"`
	Content   string    `json:"content,omitempty" db:"content" valid:"required,length(20|4294967295)"`
	UserId    int64     `json:"user_id,omitempty" db:"user_id" valid:"-"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
	User      User      `json:"user,omitempty" db:"user" valid:"-"`
}

// DataStruct for the pagination
type PostPage struct {
	WhereString string
	WhereParams []interface{}
	Order       map[string]string
	FirstId     int64
	LastId      int64
	PageNum     int
	PerPage     int
	TotalPages  int
	TotalItems  int64
	orderStr    string
}

// Current get the current page of PostPage object for pagination
func (_p *PostPage) Current() ([]Post, error) {
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("current")

	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)

	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	posts, err := FindPostsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(posts) != 0 {
		_p.FirstId, _p.LastId = posts[0].Id, posts[len(posts)-1].Id
	}
	return posts, nil
}

// Previous get the previous page of PostPage object for pagination
func (_p *PostPage) Previous() ([]Post, error) {
	if _p.PageNum == 0 {
		return nil, errors.New("This's the first page, no previous page yet")
	}
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("previous")

	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)

	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	posts, err := FindPostsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(posts) != 0 {
		_p.FirstId, _p.LastId = posts[0].Id, posts[len(posts)-1].Id
	}
	_p.PageNum -= 1
	return posts, nil
}

// Next get the next page of PostPage object for pagination
func (_p *PostPage) Next() ([]Post, error) {
	if _p.PageNum == _p.TotalPages-1 {
		return nil, errors.New("This's the last page, no next page yet")
	}
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("next")

	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)

	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	posts, err := FindPostsWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(posts) != 0 {
		_p.FirstId, _p.LastId = posts[0].Id, posts[len(posts)-1].Id
	}
	_p.PageNum += 1
	return posts, nil
}

// GetPage is a helper function for the PostPage object to return a corresponding page due to
// the parameter passed in, one of "previous, current or next"
func (_p *PostPage) GetPage(direction string) (ps []Post, err error) {
	switch direction {
	case "previous":
		ps, _ = _p.Previous()
	case "next":
		ps, _ = _p.Next()
	case "current":
		ps, _ = _p.Current()
	default:
		return nil, errors.New("Error: wrong dircetion! None of previous, current or next!")
	}
	return
}

// buildOrder is for PostPage object to build SQL ORDER clause
func (_p *PostPage) buildOrder() {
	tempList := []string{}
	for k, v := range _p.Order {
		tempList = append(tempList, fmt.Sprintf("%v %v", k, v))
	}
	_p.orderStr = " ORDER BY " + strings.Join(tempList, ", ")
}

// buildIdRestrict is for PostPage object to build a SQL clause for ID restriction,
// implementing a simple keyset style pagination
func (_p *PostPage) buildIdRestrict(direction string) (idStr string, idParams []interface{}) {
	switch direction {
	case "previous":
		if strings.ToLower(_p.Order["id"]) == "desc" {
			idStr += "id > ? "
			idParams = append(idParams, _p.FirstId)
		} else {
			idStr += "id < ? "
			idParams = append(idParams, _p.FirstId)
		}
	case "current":
		// trick to make Where function work
		if _p.PageNum == 0 && _p.FirstId == 0 && _p.LastId == 0 {
			idStr += "id > ? "
			idParams = append(idParams, 0)
		} else {
			if strings.ToLower(_p.Order["id"]) == "desc" {
				idStr += "id <= ? AND id >= ? "
				idParams = append(idParams, _p.FirstId, _p.LastId)
			} else {
				idStr += "id >= ? AND id <= ? "
				idParams = append(idParams, _p.FirstId, _p.LastId)
			}
		}
	case "next":
		if strings.ToLower(_p.Order["id"]) == "desc" {
			idStr += "id < ? "
			idParams = append(idParams, _p.LastId)
		} else {
			idStr += "id > ? "
			idParams = append(idParams, _p.LastId)
		}
	}
	if _p.WhereString != "" {
		idStr = " AND " + idStr
	}
	return
}

// buildPageCount calculate the TotalItems/TotalPages for the PostPage object
func (_p *PostPage) buildPageCount() error {
	count, err := PostCountWhere(_p.WhereString, _p.WhereParams...)
	if err != nil {
		return err
	}
	_p.TotalItems = count
	if _p.PerPage == 0 {
		_p.PerPage = 10
	}
	_p.TotalPages = int(math.Ceil(float64(_p.TotalItems) / float64(_p.PerPage)))
	return nil
}

// FindPost find a single post by an ID
func FindPost(id int64) (*Post, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID: it can't be zero")
	}
	_post := Post{}
	err := DB.Get(&_post, DB.Rebind(`SELECT * FROM posts WHERE id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_post, nil
}

// FirstPost find the first one post by ID ASC order
func FirstPost() (*Post, error) {
	_post := Post{}
	err := DB.Get(&_post, DB.Rebind(`SELECT * FROM posts ORDER BY id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_post, nil
}

// FirstPosts find the first N posts by ID ASC order
func FirstPosts(n uint32) ([]Post, error) {
	_posts := []Post{}
	sql := fmt.Sprintf("SELECT * FROM posts ORDER BY id ASC LIMIT %v", n)
	err := DB.Select(&_posts, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _posts, nil
}

// LastPost find the last one post by ID DESC order
func LastPost() (*Post, error) {
	_post := Post{}
	err := DB.Get(&_post, DB.Rebind(`SELECT * FROM posts ORDER BY id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_post, nil
}

// LastPosts find the last N posts by ID DESC order
func LastPosts(n uint32) ([]Post, error) {
	_posts := []Post{}
	sql := fmt.Sprintf("SELECT * FROM posts ORDER BY id DESC LIMIT %v", n)
	err := DB.Select(&_posts, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _posts, nil
}

// FindPosts find one or more posts by the given ID(s)
func FindPosts(ids ...int64) ([]Post, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_posts := []Post{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := DB.Rebind(fmt.Sprintf(`SELECT * FROM posts WHERE id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := DB.Select(&_posts, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _posts, nil
}

// FindPostBy find a single post by a field name and a value
func FindPostBy(field string, val interface{}) (*Post, error) {
	_post := Post{}
	sqlFmt := `SELECT * FROM posts WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := DB.Get(&_post, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_post, nil
}

// FindPostsBy find all posts by a field name and a value
func FindPostsBy(field string, val interface{}) (_posts []Post, err error) {
	sqlFmt := `SELECT * FROM posts WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = DB.Select(&_posts, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _posts, nil
}

// AllPosts get all the Post records
func AllPosts() (posts []Post, err error) {
	err = DB.Select(&posts, "SELECT * FROM posts")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

// PostCount get the count of all the Post records
func PostCount() (c int64, err error) {
	err = DB.Get(&c, "SELECT count(*) FROM posts")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// PostCountWhere get the count of all the Post records with a where clause
func PostCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM posts"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
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

// PostIncludesWhere get the Post associated models records, it's just the eager_load function
func PostIncludesWhere(assocs []string, sql string, args ...interface{}) (_posts []Post, err error) {
	_posts, err = FindPostsWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _posts, err
	}
	if len(_posts) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_posts))
	for _, v := range _posts {
		ids = append(ids, interface{}(v.Id))
	}
	return _posts, nil
}

// PostIds get all the IDs of Post records
func PostIds() (ids []int64, err error) {
	err = DB.Select(&ids, "SELECT id FROM posts")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// PostIdsWhere get all the IDs of Post records by where restriction
func PostIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := PostIntCol("id", where, args...)
	return ids, err
}

// PostIntCol get some int64 typed column of Post by where restriction
func PostIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM posts"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
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

// PostStrCol get some string typed column of Post by where restriction
func PostStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM posts"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
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

// FindPostsWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindPostsWhere(where string, args ...interface{}) (posts []Post, err error) {
	sql := "SELECT * FROM posts"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&posts, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

// FindPostBySql query use a complete SQL clause
// with placeholders, eg: FindUserBySql("SELECT * FROM users WHERE first_name = ? AND age > ? ORDER BY DESC LIMIT 1", "John", 18)
// will return only One record in the table "users" whose first_name is "John" and age elder than 18
func FindPostBySql(sql string, args ...interface{}) (*Post, error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_post := &Post{}
	err = stmt.Get(_post, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return _post, nil
}

// FindPostsBySql query use a complete SQL clause
// with placeholders, eg: FindUsersBySql("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindPostsBySql(sql string, args ...interface{}) (posts []Post, err error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&posts, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

// CreatePost use a named params to create a single Post record.
// A named params is key-value map like map[string]interface{}{"first_name": "John", "age": 23} .
func CreatePost(am map[string]interface{}) (int64, error) {
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
	sqlFmt := `INSERT INTO posts (%s) VALUES (%s)`
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(keys, ","), ":"+strings.Join(keys, ",:"))
	result, err := DB.NamedExec(sqlStr, am)
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

// Create is a method for Post to create a record
func (_post *Post) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_post)
	if !ok {
		errMsg := "Validate Post struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Post struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_post.CreatedAt = t
	_post.UpdatedAt = t
	sql := `INSERT INTO posts (title,content,user_id,created_at,updated_at) VALUES (:title,:content,:user_id,:created_at,:updated_at)`
	result, err := DB.NamedExec(sql, _post)
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

// CreateUser is a method for a Post object to create an associated User record
func (_post *Post) CreateUser(am map[string]interface{}) error {
	am["post_id"] = _post.Id
	_, err := CreateUser(am)
	return err
}

// Destroy is method used for a Post object to be destroyed.
func (_post *Post) Destroy() error {
	if _post.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyPost(_post.Id)
	return err
}

// DestroyPost will destroy a Post record specified by the id parameter.
func DestroyPost(id int64) error {
	stmt, err := DB.Preparex(DB.Rebind(`DELETE FROM posts WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// DestroyPosts will destroy Post records those specified by the ids parameters.
func DestroyPosts(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM posts WHERE id IN (?%s)`, idsHolder)
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
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

// DestroyPostsWhere delete records by a where clause restriction.
// e.g. DestroyPostsWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyPostsWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM posts WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
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

// Save method is used for a Post object to update an existed record mainly.
// If no id provided a new record will be created. A UPSERT action will be implemented further.
func (_post *Post) Save() error {
	ok, err := govalidator.ValidateStruct(_post)
	if !ok {
		errMsg := "Validate Post struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Post struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _post.Id == 0 {
		_, err = _post.Create()
		return err
	}
	_post.UpdatedAt = time.Now()
	sqlFmt := `UPDATE posts SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "title = :title, content = :content, user_id = :user_id, updated_at = :updated_at", _post.Id)
	_, err = DB.NamedExec(sqlStr, _post)
	return err
}

// UpdatePost is used to update a record with a id and map[string]interface{} typed key-value parameters.
func UpdatePost(id int64, am map[string]interface{}) error {
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
	sqlFmt := `UPDATE posts SET %s WHERE id = %v`
	setKeysArr := []string{}
	for _, v := range keys {
		s := fmt.Sprintf(" %s = :%s", v, v)
		setKeysArr = append(setKeysArr, s)
	}
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(setKeysArr, ", "), id)
	_, err := DB.NamedExec(sqlStr, am)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Update is a method used to update a Post record with the map[string]interface{} typed key-value parameters.
func (_post *Post) Update(am map[string]interface{}) error {
	if _post.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdatePost(_post.Id, am)
	return err
}

// UpdateAttributes method is supposed to be used to update Post records as corresponding update_attributes in Ruby on Rails.
func (_post *Post) UpdateAttributes(am map[string]interface{}) error {
	if _post.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdatePost(_post.Id, am)
	return err
}

// UpdateColumns method is supposed to be used to update Post records as corresponding update_columns in Ruby on Rails.
func (_post *Post) UpdateColumns(am map[string]interface{}) error {
	if _post.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdatePost(_post.Id, am)
	return err
}

// UpdatePostsBySql is used to update Post records by a SQL clause
// using the '?' binding syntax.
func UpdatePostsBySql(sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return 0, errors.New("A blank SQL clause")
	}
	sql = strings.Replace(strings.ToLower(sql), "set", "set updated_at = ?, ", 1)
	args = append([]interface{}{time.Now()}, args...)
	stmt, err := DB.Preparex(DB.Rebind(sql))
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
