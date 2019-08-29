// Package snkweb is Soliton NK Web API for GO Lang
package snkweb

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// ErrNoClient is no client error
	ErrNoClient = fmt.Errorf("No Client")
)

// WebAPI is struct of Soliton NK Web API.
type WebAPI struct {
	// Json Web Token for Soliton NK
	jwt string
	// URL of Soliton NK
	url string
	// HTTP connection to Soliton  NK
	client *http.Client
	// Web Socket connection to Soliton NK
	ws *websocket.Conn
}

// Send GET request to Soliton NK and recive responce
func (s *WebAPI) getReq(act string) ([]byte, error) {
	req, err := http.NewRequest(
		"GET",
		s.url+act, nil,
	)
	if err != nil {
		return nil, err
	}
	if s.client == nil {
		return nil, ErrNoClient
	}
	req.Header.Set("Authorization", "Bearer "+s.jwt)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Status Code = %v", resp.StatusCode)
	}
	return body, nil
}

// Login is login to Soliton NK
func (s *WebAPI) Login(url, user, pass string) error {
	loginParam := map[string]interface{}{
		"User": user,
		"Pass": pass,
	}
	loginStr, _ := json.Marshal(loginParam)
	req, err := http.NewRequest(
		"POST",
		url+"/api/login",
		bytes.NewBuffer(loginStr),
	)
	if err != nil {
		return err
	}
	if s.client == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		s.client = &http.Client{Transport: tr}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Login HTTP Status Code = %v", resp.StatusCode)
	}
	r := map[string]interface{}{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}
	var v interface{}
	var ok bool
	if v, ok = r["LoginStatus"]; !ok {
		return fmt.Errorf("No LoginStatus")
	}
	if !v.(bool) {
		return fmt.Errorf("No LoginStatus")
	}
	if v, ok = r["JWT"]; !ok {
		return fmt.Errorf("No JWT")
	}
	s.jwt = v.(string)
	s.url = url
	return nil
}

// Logout  logou from Soliton NK
func (s *WebAPI) Logout() error {
	if s.client == nil {
		return ErrNoClient
	}
	req, err := http.NewRequest(
		"PUT",
		s.url+"/api/logout", nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.jwt)
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Status Code = %v", resp.StatusCode)
	}
	return nil
}

// ResorceEnt is entory of resorce on Soliton NK.
type ResorceEnt struct {
	// GUID of Resource
	GUID string `json:"GUID"`
	// Resource name
	Name string `json:"ResourceName"`
	// Desctiption of resouce
	Descr string `json:"Description"`
	// Size of resouce
	Size int64 `json:"Size"`
	// Hash of resouce
	Hash string `json:"Hash"`
	// VersionNumber of resorce
	VersionNumber int `json:"VersionNumber"`
}

// GetResorces is get resorce list frpm Soliton NK
func (s *WebAPI) GetResorces() ([]ResorceEnt, error) {
	body, err := s.getReq("/api/resources")
	if err != nil {
		return nil, err
	}
	r := []ResorceEnt{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// CreateResorce is get resorce list frpm Soliton NK
func (s *WebAPI) CreateResorce(name, descr string, global bool) (*ResorceEnt, error) {
	if s.client == nil {
		return nil, ErrNoClient
	}
	p := map[string]interface{}{
		"Global":       global,
		"ResourceName": name,
		"Description":  descr,
	}
	sp, _ := json.Marshal(p)
	req, err := http.NewRequest(
		"POST",
		s.url+"/api/resources",
		bytes.NewBuffer(sp),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.jwt)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Login HTTP Status Code = %v", resp.StatusCode)
	}
	r := ResorceEnt{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// DownloadResorce is Download resorce from Soliton NK.
func (s *WebAPI) DownloadResorce(guid string) ([]byte, error) {
	return s.getReq("/api/resources/" + guid + "/raw")
}

// UploadResorce is Upload resorce date to Soliton NK.
func (s *WebAPI) UploadResorce(guid string, data []byte) (*ResorceEnt, error) {
	if s.client == nil {
		return nil, ErrNoClient
	}
	req, err := http.NewRequest(
		"PUT",
		s.url+"/api/resources/"+guid+"/raw", bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.jwt)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Status Code = %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := ResorceEnt{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// DeleteResorce  delete resorce from Soliton NK.
func (s *WebAPI) DeleteResorce(guid string) error {
	if s.client == nil {
		return ErrNoClient
	}
	req, err := http.NewRequest(
		"DELETE",
		s.url+"/api/resources/"+guid, nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.jwt)
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Status Code = %v", resp.StatusCode)
	}
	return nil
}

// ConnectWebsocket is Connect Web Sockecket to Soliton NK.
func (s *WebAPI) ConnectWebsocket() error {
	var err error
	a := strings.Split(s.url, ":")
	if len(a) < 2 {
		return fmt.Errorf("Bad Url")
	}
	url := "wss:" + a[1] + "/api/ws/search"

	dialer := websocket.Dialer{
		Subprotocols:    []string{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	header := http.Header{"Sec-WebSocket-Protocol": []string{s.jwt}}
	s.ws, _, err = dialer.Dial(url, header)
	if err != nil {
		s.ws = nil
		return err
	}
	var cmdList = []string{
		`{"Subs":["PONG","parse","search","attach"]}`,
		`{"type":"PONG","data":{}}`,
	}
	for _, c := range cmdList {
		_, err = s.SendWebsocketCommand(c, false)
		if err != nil {
			s.ws.Close()
			s.ws = nil
			return err
		}
	}
	return nil
}

// SendWebsocketCommand is send WebSock Command and recive Response
func (s *WebAPI) SendWebsocketCommand(tx string, noWait bool) (string, error) {
	err := s.ws.WriteMessage(websocket.TextMessage, []byte(tx))
	if err != nil {
		return "", err
	}
	if noWait {
		return "", nil
	}
	var mt int
	var msg []byte
	if mt, msg, err = s.ws.ReadMessage(); err != nil {
		return "", err
	}
	if mt != websocket.TextMessage {
		return "", fmt.Errorf("MessageType=%v", mt)
	}
	return string(msg), nil
}

// StartSearch : start search via Websocket
func (s *WebAPI) StartSearch(st, et time.Time, search string) (string, error) {
	cmd := map[string]interface{}{
		"type": "search",
		"data": map[string]interface{}{
			"SearchString": search,
			"SearchStart":  st.UTC().Format("2006-01-02T15:04:05.000Z"),
			"SearchEnd":    et.UTC().Format("2006-01-02T15:04:05.000Z"),
			"Background":   false,
		},
	}
	scmd, _ := json.Marshal(cmd)
	rx, err := s.SendWebsocketCommand(string(scmd), false)
	if err != nil {
		return "", err
	}
	r := map[string]interface{}{}
	err = json.Unmarshal([]byte(rx), &r)
	if err != nil {
		return "", err
	}
	if r["data"].(map[string]interface{})["OutputSearchSubproto"] == nil {
		return "", fmt.Errorf("%v", r)
	}
	outsub := r["data"].(map[string]interface{})["OutputSearchSubproto"].(string)
	_, err = s.SendWebsocketCommand(`{"type":"search","data":{"OK":true,"OutputSearchSubproto":"`+outsub+`"}}`, true)
	if err != nil {
		return "", err
	}
	return outsub, nil
}

// GetSearchStats : check search staus via websocket
// Return  done bool , count int, err error
func (s *WebAPI) GetSearchStats(outsub string) (bool, int, error) {
	// REQ_ENTRY_COUNT: 0x3
	rx, err := s.SendWebsocketCommand(`{"type":"`+outsub+`","data":{"ID":3}}`, false)
	if err != nil {
		return false, 0, err
	}
	r := map[string]interface{}{}
	err = json.Unmarshal([]byte(rx), &r)
	if err != nil {
		return false, 0, err
	}
	if r["data"] == nil || r["data"].(map[string]interface{})["Finished"] == nil {
		return false, 0, fmt.Errorf("")
	}
	cnt := 0
	if r["data"].(map[string]interface{})["EntryCount"] != nil {
		c := r["data"].(map[string]interface{})["EntryCount"].(float64)
		cnt = int(c)
	}
	return r["data"].(map[string]interface{})["Finished"].(bool), cnt, nil
}

// SerchResultEnt : Search result entory
type SerchResultEnt struct {
	// Time Stanmp
	Ts time.Time
	// Source IP
	Src string
	// Tag ID
	Tag float64
	// Raw Data
	Data []byte
	// Enum List
	Enums map[string]string
}

// GetSearchResult : Get Search Result
func (s *WebAPI) GetSearchResult(outsub string, start, end int) ([]SerchResultEnt, error) {
	cmd := fmt.Sprintf(`{"type":"`+outsub+`","data":{"ID":16,"EntryRange":{"First":%d,"Last":%d}}}`, start, end)
	rx, err := s.SendWebsocketCommand(cmd, false)
	if err != nil {
		return nil, err
	}
	r := map[string]interface{}{}
	err = json.Unmarshal([]byte(rx), &r)
	if err != nil {
		return nil, err
	}
	if r["data"] == nil || r["data"].(map[string]interface{})["Entries"] == nil {
		return nil, fmt.Errorf("No Entries")
	}
	ents := r["data"].(map[string]interface{})["Entries"].([]interface{})
	ret := []SerchResultEnt{}
	for _, e := range ents {
		re := SerchResultEnt{}
		ee := e.(map[string]interface{})
		if ee["TS"] == nil || ee["Data"] == nil || ee["SRC"] == nil || ee["Tag"] == nil {
			continue
		}
		re.Ts, err = time.Parse(time.RFC3339Nano, ee["TS"].(string))
		if err != nil {
			continue
		}
		re.Data, err = base64.StdEncoding.DecodeString(ee["Data"].(string))
		if err != nil {
			return nil, err
		}
		re.Tag = ee["Tag"].(float64)
		re.Src = ee["SRC"].(string)
		if ee["Enumerated"] != nil {
			enums := ee["Enumerated"].([]interface{})
			re.Enums = make(map[string]string)
			for _, enum := range enums {
				//map[Name:length Value:map[Data:TGYDAAAAAAA= Type:10] ValueStr:222796]
				enmap := enum.(map[string]interface{})
				if enmap["Name"] != nil && enmap["ValueStr"] != nil {
					re.Enums[enmap["Name"].(string)] = enmap["ValueStr"].(string)
				}
			}
		}
		ret = append(ret, re)
	}
	return ret, nil
}

// CloseSearch : Close Search
func (s *WebAPI) CloseSearch(outsub string) error {
	_, err := s.SendWebsocketCommand(`{"type":"`+outsub+`","data":{"ID":1}}`, false)
	return err
}

// CloseWebsocket : Close Websocket
func (s *WebAPI) CloseWebsocket() {
	if s.ws != nil {
		s.ws.Close()
		s.ws = nil
	}
}

// Close is close client connect for Soliton NK
func (s *WebAPI) Close() {
	if s.ws != nil {
		s.ws.Close()
		s.ws = nil
	}
	if s.client != nil {
		s.client.CloseIdleConnections()
		s.client = nil
	}
}
