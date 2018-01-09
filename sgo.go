package sgo

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/kisielk/sqlstruct"
	"github.com/lann/builder"
)

type Muppet struct {
	TableName string
	Selector  string
	Columns   string
}

var session *sql.DB

type muppetBuilder builder.Builder

func Open(driver string, username string, password string, db string) muppetBuilder {
	dataSourceName := username + ":" + password + "@/" + db
	conn, err := sql.Open(driver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	session = conn
	var query = muppetBuilder{}
	return query
}

func (b muppetBuilder) Table(name string) muppetBuilder {
	return builder.Set(b, "TableName", name).(muppetBuilder)
}

func (b muppetBuilder) Where(cond string) muppetBuilder {
	return builder.Set(b, "Selector", cond).(muppetBuilder)
}

func (b muppetBuilder) And(cond string) muppetBuilder {
	andCond, _ := builder.Get(b, "Selector")
	andCond = andCond.(string) + " AND " + cond

	return builder.Set(b, "Selector", andCond.(string)).(muppetBuilder)
}

func (b muppetBuilder) Or(cond string) muppetBuilder {
	orCond, _ := builder.Get(b, "Selector")
	orCond = orCond.(string) + " OR " + cond

	return builder.Set(b, "Selector", orCond.(string)).(muppetBuilder)
}

func cols(s interface{}) []string {
	v := reflect.ValueOf(s).Elem()
	fields := v.Type()

	cols := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		cols = append(cols, fields.Field(i).Tag.Get("sql"))
	}
	return cols
}

func (b muppetBuilder) Get(t interface{}) Muppet {
	tablename, _ := builder.Get(b, "TableName")
	selectors, _ := builder.Get(b, "Selector")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(cols(t), ", "), tablename.(string), selectors.(string))
	rows, err := session.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = sqlstruct.Scan(t, rows)
		if err != nil {
			log.Fatal(err)
		}
	}
	return builder.GetStruct(b).(Muppet)
}

func (b muppetBuilder) Insert(t interface{}) Muppet {
	tablename, _ := builder.Get(b, "TableName")

	v := reflect.ValueOf(t).Elem()
	values := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		b := v.Field(i).String()
		values = append(values, "\""+b+"\"")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s)  VALUES (%s)", tablename.(string), strings.Join(cols(t), ", "), strings.Join(values, ", "))
	stmnt, err := session.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stmnt.RowsAffected())
	return builder.GetStruct(b).(Muppet)
}

func (b muppetBuilder) Update(t interface{}) Muppet {
	tablename, _ := builder.Get(b, "TableName")
	selectors, _ := builder.Get(b, "Selector")

	v := reflect.ValueOf(t).Elem()
	values := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		b := v.Field(i).String()
		values = append(values, "\""+b+"\"")
	}

	updateQuery := make([]string, 0)
	for i, ele := range cols(t) {
		updateQuery = append(updateQuery, ele+" = "+values[i])
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tablename.(string), strings.Join(updateQuery, ", "), selectors.(string))
	stmnt, err := session.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stmnt.RowsAffected())

	return builder.GetStruct(b).(Muppet)
}

func (b muppetBuilder) Delete(t interface{}) Muppet {
	tablename, _ := builder.Get(b, "TableName")
	selectors, _ := builder.Get(b, "Selector")

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tablename.(string), selectors.(string))
	stmnt, err := session.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stmnt.RowsAffected())

	return builder.GetStruct(b).(Muppet)
}

var MuppetBuilder = builder.Register(muppetBuilder{}, Muppet{}).(muppetBuilder)
