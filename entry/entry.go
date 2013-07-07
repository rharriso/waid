package entry

type Entry struct {
	Id    int64  `db:"id"`
	Start int64  `db:"start_time"`
	End   int64  `db:"end_time"`
	Msg   string `db:"message"`
}
