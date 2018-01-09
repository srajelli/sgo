package sgo

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/kisielk/sqlstruct"
	"github.com/lann/builder"
)

type Chain struct {
	TableName string
	Selector  string
	Columns   string
}

var session *sql.DB

type queryBuilder builder.Builder

func Open(driver string, username string, password string, db string) (queryBuilder, error) {
	dataSourceName := username + ":" + password + "@/" + db
	var query = queryBuilder{}
	conn, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return query, err
	}
	session = conn
	return query, err
}

func (b queryBuilder) Table(name string) queryBuilder {
	return builder.Set(b, "TableName", name).(queryBuilder)
}

func (b queryBuilder) Where(cond string) queryBuilder {
	return builder.Set(b, "Selector", cond).(queryBuilder)
}

func (b queryBuilder) And(cond string) queryBuilder {
	andCond, _ := builder.Get(b, "Selector")
	andCond = andCond.(string) + " AND " + cond

	return builder.Set(b, "Selector", andCond.(string)).(queryBuilder)
}

func (b queryBuilder) Or(cond string) queryBuilder {
	orCond, _ := builder.Get(b, "Selector")
	orCond = orCond.(string) + " OR " + cond

	return builder.Set(b, "Selector", orCond.(string)).(queryBuilder)
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

func (b queryBuilder) Get(t interface{}) (error, Chain) {
	tablename, _ := builder.Get(b, "TableName")
	selectors, _ := builder.Get(b, "Selector")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(cols(t), ", "), tablename.(string), selectors.(string))
	rows, err := session.Query(query)
	if err != nil {
		return err, builder.GetStruct(b).(Chain)
	}
	defer rows.Close()

	for rows.Next() {
		err := sqlstruct.Scan(t, rows)
		if err != nil {
			return err, builder.GetStruct(b).(Chain)
		}
	}
	return err, builder.GetStruct(b).(Chain)
}

func (b queryBuilder) Insert(t interface{}) (error, Chain) {
	tablename, _ := builder.Get(b, "TableName")

	v := reflect.ValueOf(t).Elem()
	values := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		b := v.Field(i).String()
		values = append(values, "\""+b+"\"")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s)  VALUES (%s)", tablename.(string), strings.Join(cols(t), ", "), strings.Join(values, ", "))
	_, err := session.Exec(query)
	if err != nil {
		return err, builder.GetStruct(b).(Chain)
	}
	return err, builder.GetStruct(b).(Chain)
}

func (b queryBuilder) Update(t interface{}) (error, Chain) {
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
	_, err := session.Exec(query)
	if err != nil {
		return err, builder.GetStruct(b).(Chain)
	}

	return err, builder.GetStruct(b).(Chain)
}

func (b queryBuilder) Delete(t interface{}) (error, Chain) {
	tablename, _ := builder.Get(b, "TableName")
	selectors, _ := builder.Get(b, "Selector")

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tablename.(string), selectors.(string))
	_, err := session.Exec(query)
	if err != nil {
		return err, builder.GetStruct(b).(Chain)
	}

	return err, builder.GetStruct(b).(Chain)
}

var ChainBuilder = builder.Register(queryBuilder{}, Chain{}).(queryBuilder)
