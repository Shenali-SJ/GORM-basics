package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// User define default values
type User struct {
	gorm.Model
	Name string
	Age int `gorm:"default:18"`
	Birthday time.Time
}

type CreditCard struct {
	CardNum int
}

type Customer struct {
	gorm.Model
	Name string
	CreditCard CreditCard
}

type resultStr struct {
	Date time.Time
	Total int
}

type APIUser struct {
	ID uint
	Name string
}

type Result struct {
	ID int
	Name string
	Age int
}

func main() {
	//open database connection
	db, err := gorm.Open(sqlite.Open("testDb.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	user1 := User{
		Name:     "Anna",
		Age:      30,
		Birthday: time.Now(),
	}

	user2 := User{
		Model:    gorm.Model{},
		Name:     "Mindy",
		Age:      29,
		Birthday: time.Now(),
	}

	user3 := User{
		Model:    gorm.Model{},
		Name:     "Barry",
		Age:      25,
		Birthday: time.Now(),
	}

	user4 := User{
		Model:    gorm.Model{},
		Name:     "Rachel",
		Age:      24,
		Birthday: time.Now(),
	}

	user5 := User{
		Model:    gorm.Model{},
		Name:     "Monica",
		Age:      23,
		Birthday: time.Now(),
	}

	user6 := User{
		Model:    gorm.Model{},
		Name:     "Phoebe",
		Age:      25,
		Birthday: time.Now(),
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	// 1. Create database and insert a record
	result := db.Create(&user1)

	fmt.Println("Primary key of users : ", user1.ID)
	fmt.Println("Error? ", result)
	fmt.Println("Inserted records count : ", result.RowsAffected)
	fmt.Println()

	// 2. INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775")
	db.Select("Name", "Age", "CreatedAt").Create(&user2)

	// fields passed into omit will be omitted
	db.Omit("Age", "CreatedAt").Create(&user3)

	// 3. inserting large number of records

	// 3.1 pass a slice to create function
	users := []User{user4, user5, user6}
	db.Create(&users)

	//get the pks of records of the slice
	for _, user := range users {
		fmt.Println(user.ID)
	}

	// 3.2 use createInBatches(slice, batchSize)

	//create hooks
	// can define hooks to be implemented for,
	//	BeforeSave
	//	BeforeCreate
	//	AfterSave
	//	AfterCreate

	// 4. Create from Map
	db.Model(&User{}).Create([]map[string]interface{} {
		{"Name": "Chandler", "Age": 24},
		{"Name": "Ross", "Age": 24},
	})

	// 5. Associations

	//creditCard1 := CreditCard{23441}
	//
	//// this will insert into both users and credit cards tables
	//// foreign key should be specified, else error
	//cust1 := Customer{
	//	Model:      gorm.Model{},
	//	Name:       "",
	//	CreditCard: creditCard1,
	//}
	//
	//db.Create(&cust1)

	//skip having association
	//db.Omit("CreditCard").Create(&cust1)

	// checking functionality of default tag
	user7 := User{
		Model:    gorm.Model{},
		Name:     "Peter Pan",
		Age:      10,
		Birthday: time.Now(),
	}

	// omitting age
	// default age of 18 specified in the struct will be used
	db.Omit("Age").Create(&user7)

	// when having virtual/ generated value in database - might need to update permission
	// to skip a default value definition use ----> `gorm:"default:-"` in struct

	//upsert/ On Conflict

	user8 := User{
		Model:    gorm.Model{},
		Name:     "Dilan",
		Age:      32,
		Birthday: time.Now(),
	}

	//do nothing on conflict
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user8)

	//update all the columns except pr to new value on conflict
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&users)


	//----------------------query---------------------
	var user User

	// 1. retrieve a single object

	// ordered by primary key
	db.First(&user)
	fmt.Println(user.ID, " : ", user.Name)

	fmt.Println()

	// get one record, no specified order
	db.Take(&user)
	fmt.Println(user.ID, " : ", user.Name)

	fmt.Println()

	// get last record, ordered by pk
	db.Last(&user)
	fmt.Println(user.ID, " : ", user.Name)

	// get last record, ordered by pk
	result1 := db.First(&user)
	//count of records affected
	fmt.Println("Rows affected : ", result1.RowsAffected)
	fmt.Println("Error ? ", result1.Error) // error or nil
	//check what type of error it is
	errors.Is(result1.Error, gorm.ErrRecordNotFound)

	// works bc model is specified
	result2 := map[string]interface{}{}
	db.Model(&User{}).First(&result2)

	// this way doesn't work with First
	result3 := map[string]interface{}{}
	db.Table("users").Take(&result3)

	// if pk is not defined, results will be ordered by the first field

	// 2. Retrieving objects using pk

	//using inline condition
	//if ID is an int, all is fine. Else esp attention need to be given to avoid SQL injection
	var userU2 User
	db.First(&userU2, "3")  // as a string
	fmt.Println("ID 3 : ", userU2.Name)

	var userU3 User
	db.First(&userU3, 2) // as an int
	fmt.Println("ID 2 : ", userU3.Name)

	var usersU1 []User
	db.Find(&usersU1, []int{4, 5, 6})
	for _, u := range usersU1 {
		fmt.Println(u.Name)
	}

	fmt.Println()

	// 3. retrieving all the objects
	db.Find(&users)
	for _, u := range users {
		fmt.Println("ID : ", u.ID, ", Name : ", u.Name)
	}

	fmt.Println()

	// 4. string conditions
	var userWhere User
	db.Where("Name = ?", "Rachel").Find(&userWhere)
	fmt.Println("Where - ", userWhere.Name, userWhere.ID)

	fmt.Println()

	//all matching records
	//var usersWhere []User
	//db.Where("Name <> ? ", "Phoebe").Find(&usersWhere)
	//for _, u := range usersWhere {
	//	fmt.Println(u.Name)
	//}

	// where can be used with IN, LIKE, AND OR operators, with a specific field, BETWEEN
	// where can also be used with struct and map, slice of oks

	//struct
	//db.Where(&User{"Mindy", "20"}).First(&user)

	//to include zero values in a query condition, a map should be used

	// 5. Inline condition
	var usersInline []User
	db.Find(&usersInline, "Name = ? ", "Monica")
	for _, u := range usersInline {
		fmt.Println("ID : ", u.ID, " Name : ", u.Name,  "Age : ", u.Age)
	}

	fmt.Println()

	//Plain SQL, pk(non integer type), struct, map can be used with inline conditions
	var user23Age User
	db.Find(&user23Age, User {Age: 23})
	fmt.Println("Age 23 : ", user23Age.Name)

	fmt.Println()

	// 6.NOT
	//works similar to where
	//plain SQL, mp ,struct, slice can be used
	var usersNot []User
	db.Not(map[string]interface{}{"Name": []string{"Monica", "Chandler"}}).Find(&usersNot)
	for _, u := range usersNot {
		fmt.Println("Not : ", u.Name)
	}

	fmt.Println()

	// 7. OR
	//plain SQL map, struct can be used
	var usersOr []User
	db.Where("Name = 'Rachel'").Or(User{Age: 23}).Find(&usersOr)
	for _, u := range usersOr {
		fmt.Println("OR : ", u.Name)
	}

	fmt.Println()

	// 8. SELECT - to select specific fields
	//else will return all the fields
	var usersSelect []User
	//query ca also be passed as a string slice
	// string[]{"Name", "Age"}
	db.Select("Name", "Age").Find(&usersSelect)
	for _, u := range usersSelect {
		//note tha ID will be zero value bc we are not selecting that field in thr initial query
		fmt.Println("SELECT : ", u.Name, u.Age, u.ID)
	}

	fmt.Println()

	// 9. Limit and Offset
	// LIMIT- max number of records to retrieve
	// OFFSET - number of records to skip before starting to return records
	var usersLO []User
	db.Limit(3).Offset(5).Find(&usersLO)
	for _, u := range usersLO {
		//note tha ID will be zero value bc we are not selecting that field in thr initial query
		fmt.Println("LIMIT and OFFER : ", u.ID, u.Name)
	}

	fmt.Println()

	// 10. GROUP BY and HAVING

	// GROUP BY
	var resultGroup []User
	// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name`
	db.Model(&User{}).Select("Name, sum(Age) as total").Where("Name LIKE ?", "Ph%").Group("Name").Find(&resultGroup)
	for _, u := range resultGroup {
		//note tha ID will be zero value bc we are not selecting that field in thr initial query
		fmt.Println("GROUP BY : ", u.ID, u.Name)
	}

	fmt.Println()

	//HAVING
	var resultHaving []User
	// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"
	db.Model(&User{}).Select("Name, sum(Age) as total").Group("Name").Having("Name = ?", "Phoebe").Find(&resultHaving)
	for _, u := range resultHaving {
		//note tha ID will be zero value bc we are not selecting that field in thr initial query
		fmt.Println("HAVING : ", u.ID, u.Name)
	}

	fmt.Println()

	// 11. DISTINCT
	var resultDistinct []User
	db.Distinct("Name", "Age").Order("Name, age desc").Find(&resultDistinct)
	for _, u := range resultDistinct {
		//note tha ID will be zero value bc we are not selecting that field in thr initial query
		fmt.Println("DISTINCT : ", u.Age, u.Name)
	}

	fmt.Println()

	// 12. Smart select fields

	// specify a struct fpr API usage which can select specific fields automatically
	// check APIUser struct
	db.Model(&User{}).Limit(10).Find(&APIUser{})

	// 13. Locking
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
	// SELECT * FROM `users` FOR UPDATE

	// 14, Sub query
	//db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
	//// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

	// 15. Group conditions
	//can use group conditions to deal with complex SQL queries
	//db.Where(db.Where("pizza = ?", "pepperoni").Where(db.Where("size = ?", "small").Or("size = ?", "medium")),
	//).Or(db.Where("pizza = ?", "hawaiian").Where("size = ?", "xlarge"),
	//).Find(&Pizza{}).Statement

	// SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")

	// 16. Named arguments
	//db.Where("name1 = @name OR name2 = @name", sql.Named("name", "Ross")).Find(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

	// 17. Find to Map
	//var resultMap map[string]interface{}
	//db.Model(&User{}).First(&resultMap, "id = ?", 1)

	// 18. FirstOrInt
	// get first matched record or initialize with given conditions
	//only works with map and struct
	db.Where(User{Name: "Monica"}).FirstOrInit(&user5)
	// user -> User{ID: 111, Name: "Monica", Age: 23}

	//db.FirstOrInit(&user, map[string]interface{}{"name": "jinzhu"})
	//// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// 19. FirstOrCreate
	//Get first matched record or create a new one with given conditions
	//only works with map and struct


	// 20. FindInBatches - query and process records in batches
	// 21. Query Hooks
		// AfterFind hook
	// 22. Pluck - query single column from db and scan into a slice

	//more than one column - Select with Scan / select with find

	fmt.Println()

	// 23. scopes
	// specify commonly used queries as method calls
	var usersScopes []User
	db.Scopes(ageGreaterThan27).Find(&usersScopes)

	for _, u := range usersScopes {
		//note tha ID will be zero value bc we are not selecting that field in thr initial query
		fmt.Println("SCOPES : ", u.Age, u.Name)
	}

	fmt.Println()

	// 24. Count - get matched record count
	var count int64
	db.Model(&User{}).Where("Name = ?", "Ross").Count(&count)
	fmt.Println("No of records with name 'Ross' : ", count)


	//----------------------------update------------------------------------
	//db.Save(&user) - save all the fields

	// 1. updating single column - update
	// should have condition else error
	db.Model(&user2).Update("Name", "Emily")

	// 2. Updating multiple columns - updates
	// use struct - only update non zero values
	// or map[string]interface{}
	db.Model(&user1).Updates(User{Name: "Mike", Age: 25})

	//use SELECT and OMIT when u want to update selected fields or ignore some fields

	// 3. Update Hooks
	//BeforeSave, BeforeUpdate, AfterSave, AfterUpdate

	// 4. Batch Update
	//If we haven’t specified a record having primary key value with Model, GORM will perform a batch updates

	// 5. Block global update
	//If you perform a batch update without any conditions, GORM WON’T run it and will return ErrMissingWhereClause error by default
	//to fix that use some condionts or use raw SQL or enable AllowGlobalUpdate mode

	// 6. Update with SQL expressions
	db.Model(&user2).UpdateColumn("Age", gorm.Expr("Age - ?", 2))
	// this works with Update, Updates, UpdateColumn and Where

	// 7. Update from sub query
	//db.Table("users as u").Where("name = ?", "jinzhu").Update("company_name", db.Table("companies as c").Select("name").Where("c.id = u.company_id"))

	//If you want to skip Hooks methods and don’t track the update time when updating, you can use UpdateColumn, UpdateColumns

	// 8. Changed
	// can check whether fields has changed with Changed method
	// it should be used in Before Update Hook
	// The Changed method only works with methods Update, Updates

	//-----------------------------------Delete-----------------------------------
	//When deleting a record, the deleted value needs to have primary key or it will trigger a Batch Delete

	// 1.
	db.Delete(&user2)

	// 2.
	db.Where("Name = ?", "Barry").Delete(&user3)

	// 3. delete with pks
	db.Delete(&User{}, 7)
	db.Delete(&User{}, []int{3,4})

	// 4. delete hooks - BeforeDelete, AfterDelete

	// 5. Batch delete
	// The specified value has no primary value, GORM will perform a batch delete, it will delete all matched records

	// 6. Block global delete
	// work same as block global update

	// 7. Soft delete
	//If your model includes a gorm.DeletedAt field (which is included in gorm.Model), it will get soft delete ability automatically!
	//When calling Delete, the record WON’T be removed from the database, but GORM will set the DeletedAt‘s value to the current time, and the data is not findable with normal Query methods anymore.

	// even if you are not using gorm.Model in struct still you can have a field Deleted with type gorm.DeletedAt

	// 8. Find soft deleted records
	var softDeleted []User
	db.Unscoped().Where("Age = 27").Find(&softDeleted)
	for _, u := range softDeleted {
		fmt.Println(u.Name, u.ID, u.Age)
	}

	// 9. Delete permanently
	db.Unscoped().Delete(&softDeleted)

	fmt.Println()

	//-------------------------Raw sql--------------------------------------

	// 1. with Scan
	var resultRaw Result
	db.Raw("SELECT id, Name, Age FROM users WHERE name = ?", "Monica").Scan(&resultRaw)

	// 2. With exec
	//db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")

	// 3. DryRun mode
	//Generate SQL without executing, can be used to prepare or test generated SQL

	// 4. Row and Rows

	//gorm API or raw sql can be used
	row := db.Table("users").Where("Name = ?", "Dilan").Select("Name", "Age").Row()
	row.Scan()

	//rows - example
	//rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows()
	//defer rows.Close()
	//for rows.Next() {
	//	rows.Scan(&name, &age, &email)
	//
	//	// do something
	//}

	// ScanRows can be used to scan a row into a struct


}


func ageGreaterThan27(db *gorm.DB) *gorm.DB {
	return db.Where("Age > ?", 27)
}
