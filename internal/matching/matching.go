package matching

import (
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
