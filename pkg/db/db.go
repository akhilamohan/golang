package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	//	"github.com/golang/protobuf/jsonpb"
	//	"k8s.io/klog"
)

var (
	dbhost     = ""
	dbport     = ""
	dbuser     = ""
	dbpassword = ""
	dbname     = ""
)

func init() {
	dbhost = os.Getenv("DBHOST")
	dbport = os.Getenv("DBPORT")
	dbuser = os.Getenv("DBUSER")
	dbpassword = os.Getenv("DBPASSWORD")
	dbname = os.Getenv("DBNAME")
}

const (
	dbDriver = "mysql"
)

func ConnectToDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbuser, dbpassword, dbhost, dbport, dbname)
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("DB open failed: %v", err)
	}
	return db, nil
}

func CloseConnection(db *sql.DB) {
	db.Close()
}
