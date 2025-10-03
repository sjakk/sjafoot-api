package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"github.com/sjakk/sjafoot/internal/data"
	"testing"
)

type testServer struct {
	*httptest.Server
}

func cleanupTestDatabase(t *testing.T) {
	t.Helper()
	_, err := testApp.models.Users.DB.Exec("TRUNCATE TABLE users, torcedores RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to clean up test database: %s", err)
	}
}

func newTestServer(t *testing.T) *testServer {
	cleanupTestDatabase(t)
	ts := httptest.NewServer(testApp.routes())
	return &testServer{ts}
}

func (ts *testServer) post(t *testing.T, urlPath string, token string, body []byte) (int, http.Header, string) {
	req, err := http.NewRequest("POST", ts.URL+urlPath, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	bodyBytes, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(bodyBytes)
}

func TestHealthcheck(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/v1/healthcheck")
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("got %d; want %d", rs.StatusCode, http.StatusOK)
	}
}

func TestBroadcastAdminOnly(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	t.Run("unauthenticated", func(t *testing.T) {
		statusCode, _, _ := ts.post(t, "/broadcast", "", nil)
		if statusCode != http.StatusUnauthorized {
			t.Errorf("got %d; want %d", statusCode, http.StatusUnauthorized)
		}
	})

	adminUser := &data.User{Name: "Admin", Email: "admin@test.com", Activated: true, Role: "admin"}
	adminUser.Password.Set("password123")
	testApp.models.Users.Insert(adminUser)

	standardUser := &data.User{Name: "User", Email: "user@test.com", Activated: true, Role: "user"}
	standardUser.Password.Set("password123")
	testApp.models.Users.Insert(standardUser)

	t.Run("non-admin user", func(t *testing.T) {
		loginBody := []byte(`{"email": "user@test.com", "password": "password123"}`)
		statusCodeLogin, _, bodyStr := ts.post(t, "/auth/login", "", loginBody)

		if statusCodeLogin != http.StatusOK {
			t.Fatalf("Login for standard user failed, got status %d", statusCodeLogin)
		}

		var tokenStruct struct{ Token string }
		json.Unmarshal([]byte(bodyStr), &tokenStruct)

		broadcastBody := []byte(`{"tipo": "inicio", "time": "Test", "mensagem": "test"}`)
		statusCodeBroadcast, _, _ := ts.post(t, "/broadcast", tokenStruct.Token, broadcastBody)
		if statusCodeBroadcast != http.StatusForbidden {
			t.Errorf("got %d; want %d", statusCodeBroadcast, http.StatusForbidden)
		}
	})

	t.Run("admin user", func(t *testing.T) {
		loginBody := []byte(`{"email": "admin@test.com", "password": "password123"}`)
		statusCodeLogin, _, bodyStr := ts.post(t, "/auth/login", "", loginBody)
		
		if statusCodeLogin != http.StatusOK {
			t.Fatalf("Login for admin user failed, got status %d", statusCodeLogin)
		}

		var tokenStruct struct{ Token string }
		json.Unmarshal([]byte(bodyStr), &tokenStruct)

		broadcastBody := []byte(`{"tipo": "inicio", "time": "Test", "mensagem": "test"}`)
		statusCodeBroadcast, _, _ := ts.post(t, "/broadcast", tokenStruct.Token, broadcastBody)
		if statusCodeBroadcast != http.StatusOK {
			t.Errorf("got %d; want %d", statusCodeBroadcast, http.StatusOK)
		}
	})
}
