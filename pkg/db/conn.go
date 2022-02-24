package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// blank import for MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// Driver名
const driverName = "mysql"

// Conn 各repositoryで利用するDB接続(Connection)情報
var Conn *sql.DB

func init() {
	/* ===== データベースへ接続する. ===== */
	// ユーザ
	user := os.Getenv("MYSQL_USER")
	// パスワード
	password := os.Getenv("MYSQL_PASSWORD")
	// 接続先ホスト
	host := os.Getenv("MYSQL_HOST")
	// 接続先ポート
	port := os.Getenv("MYSQL_PORT")
	// 接続先データベース
	database := os.Getenv("MYSQL_DATABASE")

	// 接続情報は以下のように指定する.
	// user:password@tcp(host:port)/database
	var err error
	Conn, err = sql.Open(driverName,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database))
	if err != nil {
		log.Fatal(err)
	}
	if err := Conn.Ping(); err != nil {
		log.Fatalf("can't connect to mysql server. "+
			"MYSQL_USER=%s, "+
			"MYSQL_PASSWORD=%s, "+
			"MYSQL_HOST=%s, "+
			"MYSQL_PORT=%s, "+
			"MYSQL_DATABASE=%s, "+
			"error=%+v",
			user, password, host, port, database, err)
	}
}

// https://medium.com/a-journey-with-go/go-how-does-defer-statement-work-1a9492689b6e#:~:text=defer%20statement%20is%20a%20convenient,reverse%20order%20they%20were%20deferred.
func Transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Begin is failed %v", err)
		return err
	}
	defer func() {
		p := recover()
		switch {
		case err != nil:
			if err := tx.Rollback(); err != nil {
				log.Printf("Rollback is failed %v", err)
			}

		case p != nil:
			if err := tx.Rollback(); err != nil {
				log.Printf("Rollback is failed %v", err)
			}
			panic(p)

		default:
			if err := tx.Commit(); err != nil {
				log.Printf("Commit is failed %v", err)
			}
		}
		// if p := recover(); p != nil {
		// 	if err := tx.Rollback(); err != nil {
		// 		log.Printf("Rollback is failed %v", err)
		// 	}
		// 	panic(p)
		// } else if err != nil {
		// 	if err := tx.Rollback(); err != nil {
		// 		log.Printf("Rollback is failed %v", err)
		// 	}
		// 	return
		// } else {
		// 	if err := tx.Commit(); err != nil {
		// 		log.Printf("Commit is failed %v", err)
		// 	}
		// }
	}()
	err = txFunc(tx)
	return err
}
