package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UserRank struct {
	FirstName      string `json:"first_name"`
	Country        string `json:"country"`
	ProfilePicture string `json:"profile_picture"`
	TotalPoints    int    `json:"total_points"`
	Rank           int    `json:"rank"`
}

	

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/user_activities")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return db, nil
}

func GetUserRankingHandler(w http.ResponseWriter, r *http.Request) {
	db, err := ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT
			u.first_name,
			u.country,
			u.profile_picture,
			SUM(a.points) AS total_points,
			RANK() OVER (ORDER BY SUM(a.points) DESC) 
		FROM
			users u
		JOIN
			activity_logs al ON u.id = al.user_id
		JOIN
			activities a ON al.activity_id = a.id
			GROUP BY
			u.id
		ORDER BY
			total_points DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	rankings := []UserRank{}
	for rows.Next() {
		var ranking UserRank
		err := rows.Scan(&ranking.FirstName, &ranking.Country, &ranking.ProfilePicture, &ranking.TotalPoints, &ranking.Rank)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rankings = append(rankings, ranking)
	}

	jsonResult, err := json.Marshal(rankings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func main() {
	fmt.Println("hey")
	http.HandleFunc("/users/ranking", GetUserRankingHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
