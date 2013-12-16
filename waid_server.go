package main

import (
	"database/sql"
	"github.com/codegangsta/martini"
	//"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rharriso/waid/entry"
	"os/user"
)

var (
	dbMap *gorp.DbMap
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(dbConnect())

	// index route
	m.Get("/", func(r render.Render) {
		r.JSON(200, entry.All(dbMap))
	})

	// initialize server
	m.Run()
}

/*
  dbConnect ->
    connect to databse and create tables maybe
*/
func dbConnect() martini.Handler {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("sqlite3", usr.HomeDir+"/.waid.db")
	if err != nil {
		panic(err)
	}

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	// add entries table
	dbMap.AddTableWithName(entry.Entry{}, "entries").SetKeys(true, "Id")
	dbMap.CreateTablesIfNotExists()

	return func(c martini.Context) {

	}
}
