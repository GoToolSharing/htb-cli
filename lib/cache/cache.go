package cache

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
)

// Function to retrieve machines_cache_date and compare with current date
func CheckCacheDate(db *sql.DB) (bool, error) {
	var cacheDateStr string

	err := db.QueryRow("SELECT machines_cache_date FROM config").Scan(&cacheDateStr)
	if err != nil {
		if err == sql.ErrNoRows {
			config.GlobalConfig.Logger.Info("No data in the table (machines_cache_date)")
			_, err := db.Exec("INSERT into config (machines_cache_date) VALUES (CURRENT_TIMESTAMP)")
			if err != nil {
				return false, fmt.Errorf("Error updating date cache: %v", err)
			}
			config.GlobalConfig.Logger.Info("Date cache insered with current date")
			return true, nil
		}
		return false, fmt.Errorf("error retrieving machines_cache_date: %v", err)
	}

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Current date cache (machines): %s", cacheDateStr))

	cacheDate, err := time.Parse(time.RFC3339, cacheDateStr)
	if err != nil {
		return false, fmt.Errorf("error parsing machines_cache_date: %v", err)
	}

	now := time.Now()
	diff := now.Sub(cacheDate)

	if diff.Hours() > 24 {
		return true, nil
	}

	return false, nil
}
