package firebase

import (
	"context"
	"fmt"
	"log"
)

// MatchedUser is a struct containing information for a matching user.
type MatchedUser struct {
	// Name is the display name for the user.
	Name string `json:"name"`
	// Company is the company name that the user works for.
	Company string `json:"company"`
}

// GetUserSelection gets a number of ordered users from the database.
func (c *Controller) GetUserSelection(size int) ([]MatchedUser, error) {
	matches := make([]MatchedUser, size)

	query, err := c.Database.NewRef("/users").OrderByKey().LimitToFirst(size).GetOrdered(context.Background())
	if err != nil {
		log.Println("Could not get ref to /users")
		return nil, fmt.Errorf("Could not access /users in database")
	}

	for _, m := range query {
		match := &MatchedUser{}

		err = m.Unmarshal(match)
		if err != nil {
			log.Println("Could not unmarshal match.")
			continue // Skip this iteration if unmarshaller didn't work.
		}

		matches = append(matches, *match)
	}

	return matches, nil
}
