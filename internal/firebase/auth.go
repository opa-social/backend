package firebase

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
)

// AuthMiddleware creates an HTTP authentication middleware handler for authenticating
// HTTP requests against Firebase authentication.
func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Read JWT from request body.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Malformed JWT", http.StatusBadRequest)
			log.Printf("Got malformed JWT from %s", r.RemoteAddr)

			return
		}

		// Validate token against Firebase.
		_, err = c.client.VerifyIDToken(context.Background(), string(body))
		if err != nil {
			http.Error(w, "Invalid token!", http.StatusUnauthorized)
			log.Printf("Failed to authenticate user from %s", r.RemoteAddr)

			return
		}

		next.ServeHTTP(w, r)
	})
}
