package gym

type Workout struct {
	RowNum int     `db:"row_number"`
	Date   string  `db:"workout_date"`
	Weight float64 `db:"weight"`
	ID     int     `db:"id"`
}

type PullUps struct {
	RowNum   int    `db:"row_number"`
	Date     string `db:"date"`
	Quantity int    `db:"quantity"`
	ID       int    `db:"id"`
}
