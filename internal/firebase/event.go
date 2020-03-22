package firebase

import (
	"context"
	"encoding/json"
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

// EventUserResponses contains the list of user responses for the current event.
type EventUserResponses struct {
	// ID is the ID of the current event.
	ID string
	// Results is a list of all responses from all users in the event.
	Results []struct {
		// UID is the ID of the user.
		UID string
		// Responses is list of all responses to the current event.
		Responses []int
	}
}

// GetRawResponses returns a 2D slice containing the ordered responses from each user
// in integer form. This is to extract the relevant data for training.
func (e EventUserResponses) GetRawResponses() [][]float64 {
	rawResponses := make([][]float64, 0, len(e.Results))

	for _, r := range e.Results {
		row := make([]float64, 0, len(r.Responses))
		for _, v := range r.Responses {
			row = append(row, float64(v))
		}

		rawResponses = append(rawResponses, row)
	}

	return rawResponses
}

// UnmarshalJSON is the custom implementation of json.Unmarshal for the EventUserResponses
// type.
func (e *EventUserResponses) UnmarshalJSON(data []byte) error {
	type alias map[string]map[string]int
	aux := &alias{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	for k, v := range *aux {
		responses := make([]int, 0, len(v))
		for _, i := range v {
			responses = append(responses, i)
		}

		e.Results = append(e.Results, struct {
			UID       string
			Responses []int
		}{k, responses})
	}

	return nil
}

// GetAllResponses returns the stored the responses for the given event.
func (c *Controller) GetAllResponses(id string) (EventUserResponses, error) {
	responses := &EventUserResponses{
		ID: id,
	}

	err := c.Database.NewRef(fmt.Sprintf("/events/%s/users", id)).Get(context.Background(), responses)
	if err != nil {
		log.Printf("Could not deserialize list of responses for event \"%s\" because \"%s\"", id, err)
		return *responses, err
	}

	return *responses, nil
}
