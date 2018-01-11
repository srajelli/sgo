# sgo 

sgo provides easy and fun way to interact with sql database

[![Build Status](https://travis-ci.org/srajelli/sgo.svg?branch=master)](https://travis-ci.org/srajelli/sgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/srajelli/sgo)](https://goreportcard.com/report/github.com/srajelli/sgo)

### Installation
```
go get github.com/srajelli/sgo
```
### Usage
Just define tag names to your fields. Tag names is equal to column names of the table
```go
type User struct {
  Name        string `sql:"name"`
  Email       string `sql:"email"`
  AccessToken string `sql:"accesstoken"`
}

```
##### Select
```go
user := User{}
db.Table("user").
  Where("plan = 'basic'").
  And("name = 'John'").
  Or("email = 'john.due@gmail.com'").
  Get(&user)
  
fmt.Println(user)
```

##### Insert
```go
user := User{}
user.Name = "John Due"
user.Email = "john.due@gmail.com"
user.Plan = "basic"

db.Table("users").
  Insert(&user)
```
##### Update
```go
user := User{}
user.Plan = "pro"

db.Table("users").Table("users").
  Where("email = 'john.due@gmail.com'").
  Update(&user)
```

##### Delete
```go
db.Table("users").Table("users").
  Where("email = 'john.due@gmail.com'").
  Delete()
```
