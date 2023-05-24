package zutils

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/demdxx/gocast"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var msql *MSQL

type MSQL struct {
	DB *sql.DB
}

// NewDB
func NewDB(sqlType, user, passwd, host, port, database, filePath string) (*MSQL, error) {
	switch sqlType {
	case "mysql":
		dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, passwd, host, port, database, "utf8")
		//打开连接失败
		db, err := sql.Open("mysql", dbDSN)
		if err != nil {
			return nil, err
		}

		// 最大连接数
		db.SetMaxOpenConns(100)
		// 闲置连接数
		db.SetMaxIdleConns(20)
		// 最大连接周期
		db.SetConnMaxLifetime(100 * time.Second)

		if err = db.Ping(); nil != err {
			return nil, err
		}
		msql = &MSQL{DB: db}
		return msql, nil
	case "mssql":
		p, _ := strconv.Atoi(port)
		dbDSN := fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s", host, p, user, passwd, database)
		db, err := sql.Open("mssql", dbDSN)
		if err != nil {
			return nil, err
		}
		msql = &MSQL{DB: db}
		return msql, nil
	case "sqlite3":
		db, err := sql.Open("sqlite3", filePath)
		if err != nil {
			return nil, err
		}
		msql = &MSQL{DB: db}
		return msql, nil
	default:
		return nil, errors.New("not supported SQL type: " + sqlType)
	}
}

// GetDB
func GetDB() *MSQL {
	return msql
}

// QueryBy 查询操作
func (m *MSQL) QueryBy(sqlStr string, args ...interface{}) (*[]map[string]interface{}, error) {
	var maps = make([]map[string]interface{}, 0)
	// `SELECT * FROM user WHERE mobile=?`
	stmt, err := m.DB.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	//遍历每一行
	colNames, _ := rows.Columns()
	var cols = make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		cols[i] = new(interface{})
	}

	for rows.Next() {
		err := rows.Scan(cols...)
		//checkErr(err)
		if err != nil {
			return nil, err
		}
		var rowMap = make(map[string]interface{})
		for i := 0; i < len(colNames); i++ {
			rowMap[colNames[i]] = convertRow(*(cols[i].(*interface{})))
		}
		maps = append(maps, rowMap)
	}
	//fmt.Println(maps)
	return &maps, nil //返回指针
}

//ModifyBy 修改数据操作
func (m *MSQL) ModifyBy(sqlStr string, args ...interface{}) (int64, error) {
	// `INSERT user (uname, age, mobile) VALUES (?, ?, ?)`
	// "update user set mobile=? where id=?"
	// "DELETE FROM user where id=?"
	stmt, err := m.DB.Prepare(sqlStr) // Exec、Prepare均可实现增删改查
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	//判断执行结果
	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

// convertRow 行数据转换
func convertRow(row interface{}) interface{} {
	switch row.(type) {
	case int:
		return gocast.ToInt(row)
	case int32:
		return gocast.ToFloat32(row)
	case int64:
		return gocast.ToFloat64(row)
	case float32:
		return gocast.ToFloat32(row)
	case float64:
		return gocast.ToFloat64(row)
	case string:
		return gocast.ToString(row)
	case []byte:
		return gocast.ToString(row)
	case bool:
		return gocast.ToBool(row)
	}
	return row
}
