package matching

import (
	"fmt"
	"log"

	"github.com/mpraski/clusters"
	"github.com/opa-social/backend/internal/database"
	"github.com/opa-social/backend/internal/firebase"
)

// Match is a goroutine that takes questionnaire response data from Firebase, learns
// on it using the OPTICS clustering algorithm, and stores it in Redis.
func Match(id string, fb *firebase.Controller, db *database.Database) {
	responses, err := fb.GetAllResponses(id)
	if err != nil {
		log.Printf("Could not get responses for %s because \"%s\"", id, err)
		return
	}

	matcher, err := clusters.OPTICS(0, 0, 0, 4, clusters.EuclideanDistance)
	if err != nil {
		log.Printf("Could not initialize clustering because \"%s\"", err)
		return
	}

	// Learn using OPTICS on dataset.
	matcher.Learn(responses.GetRawResponses())

	err = db.ImportClustersForEvent(responses, matcher.Guesses(), len(matcher.Sizes()))
	if err != nil {
		log.Printf("Could not store matches because \"%s\"", err)
		return
	}
}

// GetMatches returns the list of matches for the current user based on previously
// stored group matches in Redis.
func GetMatches(uid, eid string, fb *firebase.Controller, db *database.Database) (firebase.MatchList, error) {
	uids, err := db.GetMatchesForUID(uid, eid, 5)
	if err != nil {
		return firebase.MatchList{}, fmt.Errorf("Could not contact database for matches")
	}

	color, err := db.GetColorForGroup(uid, eid)
	if err != nil {
		// Fall back to white.
		color = "FFFFFF"
	}

	matches, err := fb.GetUsersFromUIDs(uids, color)
	if err != nil {
		return firebase.MatchList{}, fmt.Errorf("Could not contact firebase for user data")
	}

	return matches, nil
}
