package olympicwinners

import (
	"log"
	"net/http"
)

type olympicWinners struct{}


// FindAll func
func(c *olympicWinners) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("asdasdadasdasdasdasd")
	// var rows []Model
	// var err error

	// qryStr := `SELECT * FROM olympic_winners LIMIT 1`
	// err = db.Select(&rows, qryStr)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// result, err := json.Marshal(rows)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// log.Println(string(result))
	// w.Write(result)
}


var OlympicWinners = &olympicWinners{}