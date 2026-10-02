package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIDaysToNewYear(t *testing.T) {
	srv := httptest.NewServer(routes())
	t.Cleanup(srv.Close)

	client := srv.Client()

	t.Run("valid date returns days", func(t *testing.T) {
		resp, err := client.Get(srv.URL + "/api/days?date=31-12-2023")
		if err != nil {
			t.Fatalf("GET: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, body)
		}

		ct := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			t.Fatalf("Content-Type = %q, want application/json prefix", ct)
		}

		var got Response
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if got.DaysLeft != 1 {
			t.Fatalf("days_left = %d, want 1", got.DaysLeft)
		}
	})

	t.Run("invalid date returns 400", func(t *testing.T) {
		resp, err := client.Get(srv.URL + "/api/days?date=2023-12-31")
		if err != nil {
			t.Fatalf("GET: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}

		var got ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}

		if got.Error == "" {
			t.Fatal("error message is empty")
		}
	})

	t.Run("without date uses current date", func(t *testing.T) {
		resp, err := client.Get(srv.URL + "/api/days")
		if err != nil {
			t.Fatalf("GET: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var got Response
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if got.DaysLeft < 0 || got.DaysLeft > 366 {
			t.Fatalf("days_left = %d, want range 0..366", got.DaysLeft)
		}
	})

	t.Run("wrong method returns 405", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/days", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("POST: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
		}
	})
}

func TestHealthEndpoint(t *testing.T) {
	srv := httptest.NewServer(routes())
	t.Cleanup(srv.Close)

	resp, err := srv.Client().Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var got HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Status != "ok" {
		t.Fatalf("status = %q, want %q", got.Status, "ok")
	}
}
