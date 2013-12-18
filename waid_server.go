package main

import (
	"database/sql"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rharriso/waid/entry"
	"os/user"
	"strconv"
)

var (
	dbMap *gorp.DbMap
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(dbConnect())

	// index route
	m.Get("/entries", func(r render.Render) {
		r.JSON(200, entry.All(dbMap))
	})

	// add route
	m.Post("/entries", binding.Form(entry.Entry{}), func(e entry.Entry, r render.Render) {
		dbMap.Insert(&e)
		r.JSON(200, e)
	})

	// replace route
	m.Put("/entries/:id", binding.Form(entry.Entry{}), func(params martini.Params, e entry.Entry, r render.Render) {
		en, err := dbMap.Get(entry.Entry{}, params["id"])

		if err != nil || en == nil {
			r.JSON(404, "Entry not found")
			return
		}
		//replace existing
		e.Id, _ = strconv.ParseInt(params["id"], 10, 64)
		dbMap.Update(e)
		r.JSON(200, e)
	})

	m.Delete("/entries/:id", func(params martini.Params, r render.Render) {
		e, err := dbMap.Get(entry.Entry{}, params["id"])
		if err != nil || e == nil {
			r.JSON(404, "Entry not found")
			return
		}
		dbMap.Delete(e)
		r.JSON(204, nil)
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
