package firebase

import (
	"context"
	"fmt"
	"log"
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomID() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

// EventRequest represents the content of a valid event creation
// request.
type EventRequest struct {
	// Expires is the expiry time. Not required, but defaults to 10 days from now.
	Expires uint `json:"expires,omitempty"`
	// Questionnaire is the ID of the questionnaire for the event.
	Questionnaire string `json:"questionnaire"`
}

// EventResponse contains the response values for an event creation.
type EventResponse struct {
	// EventID is the code and identifier for an event.
	EventID string `json:"eventId"`
}

type event struct {
	Manager       string   `json:"manager"`
	Expires       uint     `json:"expires"`
	Questionnaire string   `json:"questionnaire"`
	Users         []string `json:"-"` // Omit this so we don't end up loading tons of data.
}

func (c *Controller) createEventID() string {
	var newID string
	eventFiller := &event{}
	ref := c.Database.NewRef("/events")

	// Break out of loop if there is a nil error. Means that ref was not found and
	// proposed ID does not exist.
	for {
		newID = generateRandomID()
		if err := ref.Get(context.Background(), eventFiller); err == nil {
			break
		}
	}

	return newID
}

func (c *Controller) commitEvent(id string, e *event) error {
	ref := c.Database.NewRef(fmt.Sprintf("/events/%s", id))
	return ref.Set(context.Background(), e)
}

// CreateNewEvent creates a new event and returns the JSON request.
func (c *Controller) CreateNewEvent(uid string, request *EventRequest) (EventResponse, error) {
	newID := c.createEventID()
	newEvent := &event{
		Manager:       uid,
		Expires:       request.Expires,
		Questionnaire: request.Questionnaire,
	}

	err := c.commitEvent(newID, newEvent)
	if err != nil {
		log.Printf("Could not create event because \"%s\"", err)
		return EventResponse{}, fmt.Errorf("Could not create event")
	}

	return EventResponse{
		EventID: newID,
	}, nil
}
