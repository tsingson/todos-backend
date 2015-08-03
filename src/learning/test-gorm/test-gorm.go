package main

 import (
         "fmt"
         _ "github.com/go-sql-driver/mysql"
         "github.com/jinzhu/gorm"
         "time"
                "encoding/json"
 )

 // Activities table SQL :
 // id          bigint(20)  AUTO_INCREMENT
 // username    varchar(50)
 // created_on  timestamp
 // action      char(1)
 // description varchar(300)
 // visibility  char(1)

/*
  create table MyClass( id   bigint(20)   not null primary key auto_increment,  username    varchar(50),created_on  timestamp,action      char(1),description varchar(300),visibility  char(1));
    > sex int(4) not null default '0',
    > degree double(16,2));

  */

 // NOTE : In the struct, CreatedOn will be translated into created_on in sql level
 type MyClass struct {
         Id          int       `sql:"AUTO_INCREMENT"`
         Username    string    `sql:"varchar(50);unique"`
         CreatedOn   time.Time `sql:"timestamp"`
         Action      string    `sql:"type:char(1)" `
         Description string    `sql:"type:varchar(300)" `
         Visibility  string    `sql:"type:char(1)" `
 }

 func (c MyClass) TableName() string {
    return "MyClass"
}


 func main() {

         dbConn, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")

         if err != nil {
                 fmt.Println(err)
         }

         dbConn.DB()
         dbConn.DB().Ping()
         dbConn.DB().SetMaxIdleConns(10)
         dbConn.DB().SetMaxOpenConns(100)

        // dbConn.CreateTable( &MyClass{})
dbConn.AutoMigrate(&MyClass{})


         activity := MyClass{
                 Username:    "testuser",
                 CreatedOn:   time.Now().UTC(),
                 Description: "Testing",
                 Visibility:  "S",
         }

         dbConn.Create(&activity)

         // use a clean Activities struct for update purpose
         act := MyClass{}

         // Get the last record
         dbConn.Last(&act)

         fmt.Println("Before update: ", act)

         fmt.Println("Created on : ", act.CreatedOn)

         currentId := act.Id

         fmt.Println("Current ID : ", currentId)

         // update the last activity struct
         act.Username = "test test test"
         act.Description = "This is a test test test description"
         act.Visibility = "A"

         //dbConn.Save(&act) // update current id

         // best practice is to include the id for the update statement
         dbConn.Where("id = ?", currentId).Save(&act)

         // Get the last record - again
         dbConn.Last(&act)

         act.Username = "tsingson"
         act.Description = "tsignson description "

         dbConn.Save( &act)
actJson, _ := json.Marshal( &act )

var dat MyClass
if err := json.Unmarshal(actJson, &dat); err != nil {
        panic(err)
    }
    fmt.Println("Unmarshal ", dat)


         fmt.Println("After update : ", string(actJson))

 }
