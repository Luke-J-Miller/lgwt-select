package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("compares the speed of servers, returns the url of the fastet", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		defer slowServer.Close()
		fastServer := makeDelayedServer(0 * time.Millisecond)
		defer slowServer.Close()

		want := fastServer.URL
		got, err := Racer(slowServer.URL, fastServer.URL)

		if err != nil {
			t.Errorf("didn't want an error but got %v", err)
		}

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("returns an error if a server doesn't repond within 10s", func(t *testing.T) {
		ServerA := makeDelayedServer(11 * time.Millisecond)
		defer ServerA.Close()
		ServerB := makeDelayedServer(12 * time.Millisecond)
		defer ServerB.Close()

		_, err := ConfigurableRacer(ServerA.URL, ServerB.URL, (10 * time.Millisecond))

		if err == nil {
			t.Error("Expected an error but didn't get one")
		}
	})

}
func makeDelayedServer(delay time.Duration) (delayedServer *httptest.Server) {
	delayedServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
	return
}
