package todo

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func StartConnection() {
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sqlStmt := `
    create table if not exists todo (id integer not null primary key, title text, done bit);
	`
	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func GetTodos() []Todo {
	db, err := sql.Open("sqlite3", "./todo.db")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, title, done from todo")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	list := make([]Todo, 0)
	for rows.Next() {
		var id int
		var title string
		var done bool
		err = rows.Scan(&id, &title, &done)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, Todo{Id: id, Title: title, Done: done})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return list
}

func Insert(todo *Todo) (*Todo, error) {
	db, err := sql.Open("sqlite3", "./todo.db")

	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("insert into todo(title, done) values(?, ?)")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(todo.Title, todo.Done)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &Todo{Title: todo.Title, Done: todo.Done}, nil
}

func DeleteAll() error {
	db, err := sql.Open("sqlite3", "./todo.db")

	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("delete from todo")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func UpdateById(id int64, done bool) error {
	db, err := sql.Open("sqlite3", "./todo.db")

	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update todo set done = ? where id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(done, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

