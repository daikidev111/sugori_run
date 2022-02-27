package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"22dojo-online/pkg/interfaces/database"

	_ "github.com/go-sql-driver/mysql"
)

// Driver名
const driverName = "mysql"

type SQLHandler struct {
	Conn *sql.DB // Conn 各repositoryで利用するDB接続(Connection)情報
}

func NewSQLHandler() *SQLHandler {
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
	db := os.Getenv("MYSQL_DATABASE")

	// 接続情報は以下のように指定する.
	// user:password@tcp(host:port)/database
	var err error
	conn, err := sql.Open(driverName,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, db))
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatalf("can't connect to mysql server. "+
			"MYSQL_USER=%s, "+
			"MYSQL_PASSWORD=%s, "+
			"MYSQL_HOST=%s, "+
			"MYSQL_PORT=%s, "+
			"MYSQL_DATABASE=%s, "+
			"error=%+v",
			user, password, host, port, db, err)
	}
	sqlHandler := new(SQLHandler)
	sqlHandler.Conn = conn
	return sqlHandler
}

// // TODO: Make it abstract ...?
// func Transact(ctx context.Context, db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		log.Printf("Begin is failed %v", err)
// 		return err
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			if err := tx.Rollback(); err != nil {
// 				log.Printf("Rollback is failed %v", err)
// 			}
// 			panic(p)
// 		} else if err != nil {
// 			if err := tx.Rollback(); err != nil {
// 				log.Printf("Rollback is failed %v", err)
// 			}
// 		} else {
// 			if err := tx.Commit(); err != nil {
// 				log.Printf("Commit is failed: %v", err)
// 			}
// 		}
// 	}()
// 	err = txFunc(tx)
// 	return err
// }

func (handler *SQLHandler) Execute(statement string, args ...interface{}) (database.Result, error) {
	res := SQLResult{}
	result, err := handler.Conn.Exec(statement, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, nil
}

func (handler *SQLHandler) Query(statement string, args ...interface{}) (database.Rows, error) {
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		return new(SQLRows), err
	}
	row := new(SQLRows)
	row.Rows = rows
	return row, nil
}

func (handler *SQLHandler) QueryRow(statement string, args ...interface{}) database.Row {
	queryRow := handler.Conn.QueryRow(statement, args...)
	row := new(SQLRow)
	row.Row = queryRow
	return row
}

type SQLResult struct {
	Result sql.Result
}

func (r SQLResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SQLResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

type SQLRows struct {
	Rows *sql.Rows
}

type SQLRow struct {
	Row *sql.Row
}

func (r SQLRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SQLRows) Next() bool {
	return r.Rows.Next()
}

func (r SQLRows) Close() error {
	return r.Rows.Close()
}

func (r SQLRow) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}

func (r SQLRow) Err() error {
	return r.Row.Err()
}
