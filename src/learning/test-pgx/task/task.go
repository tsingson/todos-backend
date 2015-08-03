package task

import (
	"crypto/rand"
	"errors"
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/jackc/pgx"
	"io"
	"os"
)

type PostgresDB struct {
	Pool       *pgx.ConnPool
	poolConfig pgx.ConnPoolConfig
}

type PostgresTx struct {
	tx *pgx.Tx
}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

// initial all sql
func (pgdb *PostgresDB) afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("getTask", `
    select id,code, description from tasks where id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("listTask", `
    select id,code, description from tasks order by id asc
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("addTask", `
    insert into tasks(code, description) values( $1, $2 )
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("updateTask", `
    update tasks
      set code = $2, description=$3
      where id=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("deleteTask", `
    delete from tasks where id=$1
  `)
	if err != nil {
		return
	}

	/**
	_, err = conn.Prepare("transfer", `select * from transfer('Bob','Mary',14.00)`)
	if err != nil {
		return
	}
	*/

	// There technically is a small race condition in doing an upsert with a CTE
	// where one of two simultaneous requests to the shortened URL would fail
	// with a unique index violation. As the point of this demo is pgx usage and
	// not how to perfectly upsert in PostgreSQL it is deemed acceptable.
	_, err = conn.Prepare("putTask", `
    with upsert as (
      update tasks
      set code = $2, description=$3
      where id=$1
      returning *
    )
    insert into tasks(id, code, description)
    select $1, $2, $3 where not exists(select 1 from upsert)
  `)
	return
}

// connect to postgres DB
func (pgdb *PostgresDB) InitDb(dbhost, dbuser, dbpassword, dbname string) error {

	pgdb.InitConfig(dbhost, dbuser, dbpassword, dbname)
	pgdb.InitConnection()

	return nil
}

// inital PoolConfig of pgx
func (pgdb *PostgresDB) InitConfig(dbhost, dbuser, dbpassword, dbname string) error {

	pgdb.poolConfig = pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbhost,
			User:     dbuser,
			Password: dbpassword,
			Database: dbname,
			Logger:   log.New("module", "pgx"),
		},
		MaxConnections: 5,
		AfterConnect:   pgdb.afterConnect,
	}

	// = connPoolConfig
	return nil
}

// initial ConnPool of pgx
func (pgdb *PostgresDB) InitConnection() error {
	//var pool *pgx.ConnPool
	var err error

	pgdb.Pool, err = pgx.NewConnPool(pgdb.poolConfig)
	if err != nil {
		log.Info("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	log.Info("database connect sueecss")
	return nil
}

/*
// a test function for transection
func (pgdb *PostgresDB) Transfer() error {
	rows, _ := pgdb.Pool.Query("transfer") // limit 4 offset 2")

	for rows.Next() {

		var transfer string
		err := rows.Scan(&transfer)
		if err != nil {
			return err
		}
		fmt.Printf("select * from transfer('Bob','Mary',14.00) return: %s\n", transfer)
	}

	return rows.Err()
}
*/

//
func (pgdb *PostgresDB) ListTasks() error {
	rows, _ := pgdb.Pool.Query("listTask") // limit 4 offset 2")

	for rows.Next() {
		var id int32
		var code string
		var description string
		err := rows.Scan(&id, &code, &description)
		if err != nil {
			return err
		}
		fmt.Printf("%d. %s - %s\n", id, code, description)
	}

	return rows.Err()
}

func (pgdb *PostgresDB) AddTask(description string) error {

	length := len(description)
	fmt.Println("length of description is: ", length)

	if length > 0 {

		tx, err := pgdb.Pool.Begin()
		checkError(err)
		// Rollback is safe to call even if the tx is already closed, so if
		// the tx commits successfully, this is a no-op
		defer tx.Rollback()

		code, _ := newUUID()

		_, err = pgdb.Pool.Exec("addTask", code, description)
		checkError(err)
		err = tx.Commit()
		checkError(err)

	} else {
		fmt.Println(" description is null")
	}

	return nil
}

func (pgdb *PostgresDB) UpdateTask(itemNum int32, code string, description string) error {

	tx, err := pgdb.Pool.Begin()
	checkError(err)
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback()

	_, err = pgdb.Pool.Exec("putTask", itemNum, code, description)

	checkError(err)
	err = tx.Commit()

	return checkError(err)

}

func (pgdb *PostgresDB) RemoveTask(itemNum int32) error {

	commandTag, err := pgdb.Pool.Exec("deleteTask", itemNum)

	if commandTag.RowsAffected() != 1 {
		log.Info("No row found to delete", "error", err)
		return errors.New("No row found to delete")
	}

	return err

}

func checkError(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}

//**************************************************

// Generate a uuid to use as a unique identifier for each Todo
// http://play.golang.org/p/4FkNSiUDMg
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
