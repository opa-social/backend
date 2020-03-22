package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/opa-social/backend/internal/colors"
	"github.com/opa-social/backend/internal/firebase"
)

// Database type is the control interface for the matching datastore.
type Database struct {
	// URI is the Redis instance URI.
	URI    string
	client *redis.Client
}

// New creates a new database type.
func New(uri string, password string) Database {
	if !strings.HasPrefix(uri, "redis://") {
		log.Fatal("Redis URI does not start with proper schema.")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Fatalf("Could not connect to %s: \"%s\"", uri, err)
	}

	log.Print("Connected to database")

	return Database{
		URI:    uri,
		client: client,
	}
}

// CreateEvent creates the metadata mapping for an event.
func (d *Database) CreateEvent(id string, expires time.Time) error {
	newEventResult := d.client.HSetNX(fmt.Sprintf("meta:%s", id), "clusters", 0)
	if err := newEventResult.Err(); err != nil {
		log.Printf("Could not create event \"%s\" because \"%s\"", id, err)
		return err
	}

	// Terminate here if key already exists.
	if !newEventResult.Val() {
		return nil
	}

	expireResult := d.client.ExpireAt(id, expires)
	if err := expireResult.Err(); err != nil {
		log.Printf("Could not set expiry on event \"%s\" because \"%s\"", id, err)
		return err
	}

	return nil
}

// ImportClustersForEvent adds each user in an event to the sorted set corresponding to the event
// and sets their score to their assigned cluster ID.
func (d *Database) ImportClustersForEvent(responses firebase.EventUserResponses, clusters []int, numClusters int) error {
	_, err := d.client.TxPipelined(func(pipe redis.Pipeliner) error {
		for i, r := range responses.Results {
			pipe.ZAdd(responses.ID, &redis.Z{
				Score:  float64(clusters[i]),
				Member: r.UID,
			})
		}

		metaKey := fmt.Sprintf("meta:%s", responses.ID)

		// Set the number of clusters.
		pipe.HSet(metaKey, "clusters", numClusters)

		// Set color metadata information.
		for i, c := range colors.GenerateRandomColors(numClusters) {
			pipe.HSet(metaKey, fmt.Sprintf("color:%d", i+1), c)
		}

		return nil
	})

	return err
}

func (d *Database) getClusterID(uid, eid string) (float64, error) {
	score := d.client.ZScore(eid, uid)
	if err := score.Err(); err != nil {
		return 0.0, err
	}

	return score.Val(), nil
}

// GetMatchesForUID gets a list of `limit` matches for the given uid in the given event ID.
func (d *Database) GetMatchesForUID(uid, eid string, limit int64) ([]string, error) {
	clusterID, err := d.getClusterID(uid, eid)
	if err != nil {
		log.Printf("Could not get cluster for %s:%s because \"%s\"", eid, uid, err)
		return []string{}, err
	}

	result := d.client.ZRangeByScore(eid, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", clusterID),
		Max:    fmt.Sprintf("%f", clusterID),
		Offset: 0,
		Count:  limit,
	})

	if err := result.Err(); err != nil {
		log.Printf("Could not get matches for %s:%s because \"%s\"", eid, uid, err)
		return []string{}, err
	}

	return result.Val(), nil
}

// GetColorForGroup gets the current assigned color for the given event and user.
func (d *Database) GetColorForGroup(uid, eid string) (string, error) {
	clusterID, err := d.getClusterID(uid, eid)
	if err != nil {
		log.Printf("Could not get cluster for %s:%s because \"%s\"", eid, uid, err)
		return "", err
	}

	result := d.client.HGet(fmt.Sprintf("meta:%s", eid), fmt.Sprintf("color:%.0f", clusterID))
	if err := result.Err(); err != nil {
		log.Printf("Could not get color for %s because \"%s\"", eid, err)
		return "", err
	}

	return result.Val(), nil
}
