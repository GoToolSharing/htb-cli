package machines

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Machine struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	OS          string    `json:"os"`
	Difficulty  string    `json:"difficulty"`
	Star        float64   `json:"star"`
	Status      string    `json:"status"`
	ReleaseDate time.Time `json:"release"`
	UserOwns    bool      `json:"user"`
	RootOwns    bool      `json:"root"`
}

// Function to retrieve data and format it like the API
func GetMachinesFromCache(db *sql.DB, status string) ([]interface{}, error) {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT id, name, os, difficulty, star, status, release_date, user, root
		FROM machines WHERE status='%s'
	`, status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []interface{}

	for rows.Next() {
		var (
			id          int
			name        string
			os          string
			difficulty  string
			star        float64
			status      string
			releaseDate string
			userOwns    bool
			rootOwns    bool
		)

		err := rows.Scan(&id, &name, &os, &difficulty, &star, &status, &releaseDate, &userOwns, &rootOwns)
		if err != nil {
			return nil, err
		}

		difficultyKey := "difficultyText"

		if status == "Scheduled" {
			difficultyKey = "difficulty_text"
		}

		machineMap := map[string]interface{}{
			"id":                 id,
			"name":               name,
			"os":                 os,
			difficultyKey:        difficulty,
			"star":               star,
			"active":             status,
			"release":            releaseDate,
			"authUserInUserOwns": userOwns,
			"authUserInRootOwns": rootOwns,
		}

		machines = append(machines, machineMap)
	}

	return machines, nil
}

// Function to insert machines retrieved from the API into the database
func InsertMachines(db *sql.DB, data interface{}, title string) error {
	machines, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("the data is in the wrong format")
	}

	for _, item := range machines {
		machineMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		id := int(machineMap["id"].(float64))
		name := machineMap["name"].(string)
		os := machineMap["os"].(string)
		releaseDateStr := machineMap["release"].(string)
		status := title
		var difficulty string
		var star float64
		var userOwns bool
		var rootOwns bool

		// Scheduled machines
		if val, ok := machineMap["difficulty_text"].(string); ok {
			difficulty = val
			star = 0
			userOwns = false
			rootOwns = false
		} else {
			difficulty = machineMap["difficultyText"].(string)
			star = machineMap["star"].(float64)
			userOwns = machineMap["authUserInUserOwns"].(bool)
			rootOwns = machineMap["authUserInRootOwns"].(bool)
		}

		releaseDate, err := time.Parse("2006-01-02T15:04:05.000000Z", releaseDateStr)
		if err != nil {
			return fmt.Errorf("date parsing error for %s: %v", name, err)
		}

		_, err = db.Exec(`
			INSERT INTO machines (id, name, os, difficulty, star, status, release_date, user, root)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, id, name, os, difficulty, star, status, releaseDate, userOwns, rootOwns)
		if err != nil {
			return fmt.Errorf("error during %s insertion: %v", name, err)
		}
	}

	return nil
}
