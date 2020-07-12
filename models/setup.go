package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//DB is goOrm connection pointer ref
var DB *gorm.DB

var ctx = context.Background()

//Rdb is redis connection pointer ref
var Rdb *redis.Client

// ConnectDataBase initilize mysql db conection
func ConnectDataBase() {
	database, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database!")
	}
	database.AutoMigrate(&User{})
	database.LogMode(true)

	DB = database

	ConnectRedis()
	getTableAllData(getTablesNameFromDB())
}

// ConnectRedis initilize redis server conection
func ConnectRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	Rdb = rdb
}

//SetInRedis set value from redis method
func SetInRedis(key, val string) {
	err := Rdb.Set(key, val, 0).Err()
	if err != nil {
		panic(err)
	}
}

//GetInRedis get value from redis method
func GetInRedis(key string) (string, error) {
	val, err := Rdb.Get(key).Result()
	return val, err
}

//GetKeys
func GetKeys(pattern string) ([]string, error) {
	return Rdb.Keys(pattern).Result()
}

//GetAllHMRedis have to explain
func GetAllHMRedis(key string) (map[string]string, error) {
	result, err := Rdb.HGetAll(key).Result()
	return result, err
}

func getTablesNameFromDB() []string {
	rows, _ := DB.
		Table("information_schema.tables").
		Select("table_name").
		Where("table_schema = ?", "test").
		Rows()
	defer rows.Close()
	var tables []string
	var name string
	for rows.Next() {
		rows.Scan(&name)
		tables = append(tables, name)
	}
	return tables
}

func getTableAllData(tablesNames []string) {
	for _, tableName := range tablesNames {
		tableDatas := getDataAsJSON("SELECT * FROM " + tableName)
		for _, tableDatas := range tableDatas {
			redisCaheArray(tableDatas, tableName)
		}
	}
}

func redisCaheArray(tableDatas interface{}, tableName string) {
	tableData := tableDatas.(map[string]interface{})
	uid := strconv.FormatInt(tableData["uid"].(int64), 10)
	uname := tableData["uname"].(string)
	Rdb.HMSet(tableName+":"+uid+":"+uname, tableData)
}

//WriteInDBAndRdb ...(interface{}, error)
func WriteInDBAndRdb(tableDatas interface{}, tableName string) int64 {
	db, err := sql.Open("mysql", "root:@/test")
	checkErr(err)
	defer db.Close()

	dataToInsert := tableDatas.(map[string]interface{})
	createdAt := makeTimestamp()
	stmt, err := db.Prepare("INSERT into " + tableName + "(uname,password,created_at) VALUES(?,?,?)")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(dataToInsert["uname"], dataToInsert["password"], createdAt)
	checkErr(err)

	uid, _ := res.LastInsertId()
	dataToInsert["uid"] = uid
	dataToInsert["created_at"] = createdAt
	go redisCaheArray(dataToInsert, tableName)

	return uid
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getDataAsJSON(sqlString string) []interface{} {
	db, err := sql.Open("mysql", "root:@/test")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(sqlString)
	checkErr(err)

	columnTypes, err := rows.ColumnTypes()
	checkErr(err)

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT":
				scanArgs[i] = new(sql.NullInt64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		checkErr(err)

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}
	return finalRows
}
