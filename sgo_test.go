package sgo_test

import (
	"testing"

	_ "github.com/proullon/ramsql/driver"
	"github.com/srajelli/sgo"
)

type User struct {
	Name  string `sql:"name"`
	Email string `sql:"email"`
	Plan  string `sql:"plan"`
}

var session sgo.QueryBuilder

func TestLoadDB(t *testing.T) {
	batch := []string{
		`CREATE TABLE users (id BIGSERIAL PRIMARY KEY, email VARCHAR(255), name VARCHAR(255), plan VARCHAR(255));`,
	}

	db, err := sgo.Open("ramsql", "TestLoadDB")
	session = db

	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}

	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

}

func TestInsert(t *testing.T) {
	user := User{}
	user.Name = "John Due"
	user.Email = "john.due@gmail.com"
	user.Plan = "basic"

	err := session.Table("users").Insert(&user)
	if err != nil {
		t.Fatalf("sgo.Insert: Error: %s\n", err)
	}
}

func TestGet(t *testing.T) {
	user := User{}

	err := session.Table("users").Where("plan = 'basic'").Get(&user)
	if err != nil {
		t.Fatalf("sgo.Get: Error: %s\n", err)
	}
}

func TestUpdate(t *testing.T) {
	user := User{}
	user.Plan = "pro"

	err := session.Table("users").Where("plan = 'basic'").Update(&user)
	if err != nil {
		t.Fatalf("sgo.Update: Error: %s\n", err)
	}
}
func TestDelete(t *testing.T) {
	user := User{}
	user.Plan = "pro"

	err := session.Table("users").Where("plan = 'pro'").Delete(&user)
	if err != nil {
		t.Fatalf("sgo.Delete: Error: %s\n", err)
	}
	session.Close()
}
