package sabnzbd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toomanybyt3s/sab_monitarr/internal/sabnzbd"
)

func mockServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"status": "ok",
			"queue": {
				"status": "Downloading",
				"speed": "2.5 MB/s",
				"sizeleft": "500 MB",
				"timeleft": "00:03:20",
				"slots": [
					{
						"filename": "test_file.mkv",
						"status": "Downloading",
						"sizeleft": "500 MB",
						"percentage": "75",
						"timeleft": "00:03:20"
					}
				]
			}
		}`))
	}))
}

func TestFetchStatus(t *testing.T) {
	srv := mockServer(t)
	defer srv.Close()

	status, err := sabnzbd.FetchStatus(srv.URL, "test-api-key", false)
	if err != nil {
		t.Fatalf("FetchStatus returned error: %v", err)
	}
	if status.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", status.Status)
	}
	if status.Queue.Status != "Downloading" {
		t.Errorf("expected queue status 'Downloading', got '%s'", status.Queue.Status)
	}
	if status.Queue.Speed != "2.5 MB/s" {
		t.Errorf("expected speed '2.5 MB/s', got '%s'", status.Queue.Speed)
	}
	if len(status.Queue.Slots) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(status.Queue.Slots))
	}
	if status.Queue.Slots[0].Filename != "test_file.mkv" {
		t.Errorf("expected filename 'test_file.mkv', got '%s'", status.Queue.Slots[0].Filename)
	}
	if status.Queue.Slots[0].Percentage != "75" {
		t.Errorf("expected percentage '75', got '%s'", status.Queue.Slots[0].Percentage)
	}
}

func TestFetchStatusAPIKeyQueryParam(t *testing.T) {
	var receivedKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.URL.Query().Get("apikey")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","queue":{"status":"Idle","speed":"0","sizeleft":"0","timeleft":"0","slots":[]}}`))
	}))
	defer srv.Close()

	if _, err := sabnzbd.FetchStatus(srv.URL, "secret-key", false); err != nil {
		t.Fatalf("FetchStatus returned error: %v", err)
	}
	if receivedKey != "secret-key" {
		t.Errorf("expected apikey query param 'secret-key', got '%s'", receivedKey)
	}
}

func TestFetchStatusNonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	if _, err := sabnzbd.FetchStatus(srv.URL, "bad-key", false); err == nil {
		t.Error("expected error for non-OK response, got nil")
	}
}
