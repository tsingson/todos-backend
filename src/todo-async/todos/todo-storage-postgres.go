package todos

import (
	"errors"
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/jackc/pgx"

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

// connect to postgres DB
func (pgdb *PostgresDB) InitDb(dbhost, dbuser, dbpassword, dbname string) error {

	pgdb.InitConfig(dbhost, dbuser, dbpassword, dbname)
	pgdb.InitConnection()

	return nil
}

// initial all sql
func (pgdb *PostgresDB) afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("getTask", `
    select id,code, description from tasks where id=$1
  `)
	checkError(err)

	_, err = conn.Prepare("listTask", `
    select id,code, description from tasks order by id asc
  `)
	checkError(err)

	_, err = conn.Prepare("addTask", `
    insert into tasks(code, description) values( $1, $2 )
  `)
	checkError(err)

	_, err = conn.Prepare("updateTask", `
    update tasks
      set code = $2, description=$3
      where id=$1
  `)
	checkError(err)

	_, err = conn.Prepare("deleteTask", `
    delete from tasks where id=$1
  `)
	checkError(err)

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
	checkError(err)
	return
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

	todos := make(map[string]Todo, 0)
	var todo Todo

	rows, _ := pgdb.Pool.Query("listTask") // limit 4 offset 2")

	for rows.Next() {
		var id int32
		var code string
		var description string
		err := rows.Scan(&id, &code, &description)
		checkError(err)
		//fmt.Printf("%d. %s - %s\n", id, code, description)

		newId := fmt.Sprintf("%s", id)
		todo = Todo{Id: newId, Code: code, Description: description}
		todos[todo.Id] = todo

		fmt.Printf("%d - %s - %s \n", todo.Id, todo.Code, todo.Description)

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

	return err

}

func (pgdb *PostgresDB) RemoveTask(itemNum int32) error {

	commandTag, err := pgdb.Pool.Exec("deleteTask", itemNum)

	if commandTag.RowsAffected() != 1 {
		log.Info("No row found to delete", "error", err)
		return errors.New("No row found to delete")
	}

	return err

}

/*
func checkError(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}
*/
