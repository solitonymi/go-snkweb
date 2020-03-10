package snkweb

import (
	"os"
	"testing"
	"time"
)

// TestResource test resource access
func TestSnkWeb(t *testing.T) {
	url := os.Getenv("SNK_URL")
	if url == "" {
		t.Fatal("No Url")
	}
	uid := os.Getenv("SNK_UID")
	if uid == "" {
		t.Fatal("No UID")
	}
	passwd := os.Getenv("SNK_PASSWD")
	if passwd == "" {
		t.Fatal("No Password")
	}
	s := &WebAPI{}
	if err := s.Login(url, uid, passwd); err != nil {
		t.Fatalf("Login error=%#v", err)
	}
	defer s.Close()
	// Create Resource
	r, err := s.CreateResource("test.txt", "test", false)
	if err != nil {
		t.Fatalf("CreateResource error=%#v", err)
	}
	r2, err := s.UploadResource("test.txt", r.GUID, []byte("test"))
	if err != nil {
		t.Fatalf("UploadResource error=%#v", err)
	}
	if r2.GUID != r.GUID {
		t.Errorf("Upload Resource GUID Missmatch %#v != %#v", r.GUID, r2.GUID)
	}
	list, err := s.GetResources()
	if err != nil {
		t.Fatalf("GetResources error=%#v", err)
	}
	for _, re := range list {
		if re.GUID == r.GUID {
			if r.Name != re.Name {
				t.Errorf("Resource Name %#v != %#v", r.Name, re.Name)
			}
		}
	}
	err = s.DeleteResource(r.GUID)
	if err != nil {
		t.Fatalf("DeleteResource error=%#v", err)
	}
	if err := s.Logout(); err != nil {
		t.Fatalf("Logout error =%#v", err)
	}
	t.Log("Test Done!")
}

// TestSearch : Search Test via Websocket
func TestSearch(t *testing.T) {
	url := os.Getenv("SNK_URL")
	if url == "" {
		t.Fatal("No Url")
	}
	uid := os.Getenv("SNK_UID")
	if uid == "" {
		t.Fatal("No UID")
	}
	passwd := os.Getenv("SNK_PASSWD")
	if passwd == "" {
		t.Fatal("No Password")
	}
	s := &WebAPI{}
	if err := s.Login(url, uid, passwd); err != nil {
		t.Fatalf("Login error=%#v", err)
	}
	defer s.Close()
	err := s.ConnectWebsocket()
	if err != nil {
		t.Fatalf("ConnectWebsocket error=%#v", err)
	}
	rx, err := s.SendWebsocketCommand(`{"type":"parse","data":{"SearchString":"tag=syslog"}}`, false)
	if err != nil {
		t.Fatalf("SendWebsocketCommand error=%#v", err)
	}
	t.Log(rx)
	et := time.Now()
	st := et.Add(-60 * time.Second)
	outsub, err := s.StartSearch(st, et, "tag=syslog length")
	if err != nil {
		t.Fatalf("StartSearch error=%#v", err)
	}
	done := false
	count := 0
	for !done {
		done, count, err = s.GetSearchStats(outsub)
		if err != nil {
			t.Fatalf("GetSearchStats error=%#v", err)
		}
		time.Sleep(1 * time.Second)
	}
	t.Logf("Search Count=%d", count)
	if count > 0 {
		resluts, err := s.GetSearchResult(outsub, 0, count)
		if err != nil {
			t.Fatalf("GetSearchResult error=%#v", err)
		}
		if len(resluts) != count {
			t.Fatalf("GetSearchResult count missmatch! len=%d", len(resluts))
		}
		t.Log(string(resluts[0].Data))
		t.Log(resluts[0].Enums)
	}
	if err := s.Logout(); err != nil {
		t.Fatalf("Logout error =%#v", err)
	}
	t.Log("Test Done!")
}
