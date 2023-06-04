# GO-MYSQL--Assignment-1


C:\Users\sus>go get -u github.com/go-sql-driver/mysql
go: go.mod file not found in current directory or any parent directory.
        'go get' is no longer supported outside a module.
# how to connnect Go and Mysql
1. go mod init <module-name>
like, go mod init example.com/myproject
2. go get -u github.com/go-sql-driver/mysql
3. go install github.com/go-sql-driver/mysql
4. Define a struct to hold the result data
 Create a struct to hold the result data returned by the API:
  type UserRank struct {
    FirstName      string `json:"first_name"`
    Country        string `json:"country"`
    ProfilePicture string `json:"profile_picture"`
    TotalPoints    int    `json:"total_points"`
    Rank           int    `json:"rank"`
}
5.func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database_name")  
  //here no password used
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
  
  6.
  
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
			SUM(a.point) AS total_points,
			RANK() OVER (ORDER BY SUM(a.point) DESC) AS rank
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

  7./Applications/XAMPP/phpMyAdmin/config.inc.php
  
  $cfg['Servers'][$i]['user'] = 'root'; // MySQL username
$cfg['Servers'][$i]['password'] = ''; // MySQL password
  
  
  
8.go build
9. ./<executable-name>
10. The API server will start running and listening on port 8080. You can now make GET requests to 
  http://localhost:8080/users/ranking to fetch the user rankings from the MySQL database.
