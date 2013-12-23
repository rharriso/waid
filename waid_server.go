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

	// get lasted
	m.Get("/entries/latest", binding.Json(entry.Entry{}), func(params martini.Params, e entry.Entry, r render.Render) {
		r.JSON(200, entry.Latest(dbMap))
	})

	// add route
	m.Post("/entries", binding.Json(entry.Entry{}), func(params martini.Params, e entry.Entry, r render.Render) {
		err := dbMap.Insert(&e)
		if err != nil {
			r.JSON(404, "Unable to update entry.")
			return
		}
		r.JSON(200, e)
	})

	// add route
	m.Delete("/entries", func(r render.Render) {
		err := dbMap.TruncateTables()
		if err != nil {
			r.JSON(404, "Unable to remove all entries.")
			return
		}
		r.JSON(202, nil)
	})

	// replace route
	m.Put("/entries/:id", binding.Json(entry.Entry{}), func(params martini.Params, e entry.Entry, r render.Render) {
		en, err := dbMap.Get(entry.Entry{}, params["id"])

		if err != nil || en == nil {
			r.JSON(404, "Entry not found")
			return
		}
		//replace existing
		_, err = dbMap.Update(&e)
		if err != nil {
			r.JSON(404, "Unable to update entry.")
			return
		}
		r.JSON(200, e)
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
