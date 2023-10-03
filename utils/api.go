package utils

import (
	"fmt"
	"time"
)

const dateFormat = "2006-01-02"

// ParseAndFormatDate takes an input date string, parses it into a time.Time object, and formats it to the "2006-01-02" format.
func ParseAndFormatDate(input string) (string, error) {
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		return "", fmt.Errorf("error parsing date [%s]: %v", input, err)
	}
	return t.Format(dateFormat), nil
}

// SetStatus determines the status based on user and root flags.
func SetStatus(data map[string]interface{}) string {
	userFlag, userFlagExists := data["authUserInUserOwns"].(bool)
	rootFlag, rootFlagExists := data["authUserInRootOwns"].(bool)

	switch {
	case !userFlagExists && !rootFlagExists:
		return "No flags"
	case userFlag && !rootFlag:
		return "User flag"
	case !userFlag && rootFlag:
		return "Root flag"
	case userFlag && rootFlag:
		return "User & Root"
	default:
		return "No flags"
	}
}

// SetRetiredStatus determines whether an item is retired or not.
func SetRetiredStatus(data map[string]interface{}) string {
	if retired, exists := data["retired"].(bool); exists && !retired {
		return "Yes"
	}
	return "No"
}
