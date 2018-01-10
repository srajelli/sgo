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

type QueryBuilder builder.Builder

func Open(driver string, dataSourceName string) (QueryBuilder, error) {
	//dataSourceName := username + ":" + password + "@/" + db
	var query = QueryBuilder{}
	conn, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return query, err
	}
	session = conn
	return query, err
}

func (b QueryBuilder) Close() error {
	err := session.Close()
	return err
}

func (b QueryBuilder) Exec(query string) (sql.Result, error) {
	resp, err := session.Exec(query)
	return resp, err
}

func (b QueryBuilder) Table(name string) QueryBuilder {
	return builder.Set(b, "TableName", name).(QueryBuilder)
}

func (b QueryBuilder) Where(cond string) QueryBuilder {
	return builder.Set(b, "Selector", cond).(QueryBuilder)
}

func (b QueryBuilder) And(cond string) QueryBuilder {
	andCond, _ := builder.Get(b, "Selector")
	andCond = andCond.(string) + " AND " + cond

	return builder.Set(b, "Selector", andCond.(string)).(QueryBuilder)
}

func (b QueryBuilder) Or(cond string) QueryBuilder {
	orCond, _ := builder.Get(b, "Selector")
	orCond = orCond.(string) + " OR " + cond

	return builder.Set(b, "Selector", orCond.(string)).(QueryBuilder)
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

func (b QueryBuilder) Get(t interface{}) error {
	tablename, _ := builder.Get(b, "TableName")
	selectors, _ := builder.Get(b, "Selector")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(cols(t), ", "), tablename.(string), selectors.(string))
	rows, err := session.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := sqlstruct.Scan(t, rows)
		if err != nil {
			return err
		}
	}
	return err
}

func (b QueryBuilder) Insert(t interface{}) error {
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
		return err
	}
	return err
}

func (b QueryBuilder) Update(t interface{}) error {
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
		return err
	}

	return err
}

func (b QueryBuilder) Delete(t interface{}) error {
	tablename, _ := builder.Get(b, "TableName")
	selectors, _ := builder.Get(b, "Selector")

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tablename.(string), selectors.(string))
	_, err := session.Exec(query)
	if err != nil {
		return err
	}

	return err
}

var ChainBuilder = builder.Register(QueryBuilder{}, Chain{}).(QueryBuilder)
