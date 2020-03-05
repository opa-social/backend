package firebase

import (
	"context"
	"log"
	"net/http"
	"strings"
)

// AuthMiddleware creates an HTTP authentication middleware handler for authenticating
// HTTP requests against Firebase authentication.
func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")

		if len(authHeader) < 2 {
			http.Error(w, "No token provided", http.StatusBadRequest)
			log.Printf("Client at %s sent empty token or invalid header", r.RemoteAddr)

			return
		}

		// Get the token from the header. JWT is sent as "Authorization: Bearer <token>"
		// so "<token>" is always the second value in the array.
		// Validate token against Firebase.
		user, err := c.Client.VerifyIDToken(context.Background(), authHeader[1])
		if err != nil {
			http.Error(w, "Invalid token!", http.StatusUnauthorized)
			log.Printf("Failed to authenticate user from %s with error \"%s\"", r.RemoteAddr, err)

			return
		}

		// Add user's UID as a custom header to forwarder response.
		r.Header.Add("X-OPA-UID", user.UID)

		next.ServeHTTP(w, r)
	})
}
