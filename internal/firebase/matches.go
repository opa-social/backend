package firebase

import (
	"context"
	"fmt"
	"log"

	"github.com/opa-social/backend/internal/colors"
)

// MatchedUser is a struct containing information for a matching user.
type MatchedUser struct {
	// UID is the UID of the matched user.
	UID string `json:"uid"`
	// Name is the display name for the user.
	Name string `json:"name"`
	// Company is the company name that the user works for.
	Company string `json:"company"`
	// Phone is the phone number for the current user.
	Phone string `json:"phone"`
}

// MatchList is a struct containing the fields for the list of matches for
// the current user and the color associated with their group.
type MatchList struct {
	// Color is the color for the given group of matches.
	Color string `json:"color"`
	// Matches is the list of matches for the current user.
	Matches []*MatchedUser `json:"matches"`
}

// GetUserSelection gets a number of ordered users from the database.
func (c *Controller) GetUserSelection(size int, eventID string) (MatchList, error) {
	var matches = []*MatchedUser{}

	query, err := c.Database.NewRef(fmt.Sprintf("/events/%s/users", eventID)).
		OrderByKey().
		LimitToFirst(size).
		GetOrdered(context.Background())
	if err != nil {
		log.Println("Could not get ref to /users")
		return MatchList{}, fmt.Errorf("Could not access /users in database")
	}

	for _, m := range query {
		match := &MatchedUser{}

		err = c.Database.NewRef("/users").Child(m.Key()).Get(context.Background(), match)
		if err != nil {
			log.Println("Could not unmarshal match.")
			continue // Skip this iteration if unmarshaller didn't work.
		}

		match.UID = m.Key()
		matches = append(matches, match)
	}

	return MatchList{
		Color:   colors.GenerateRandomColors(1)[0],
		Matches: matches,
	}, nil
}
