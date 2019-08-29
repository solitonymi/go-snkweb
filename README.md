# go-snkweb

Soliton NK Web API for GO Lang

[![Godoc Reference](https://godoc.org/github.com/solitonymi/go-snkweb?status.svg)](http://godoc.org/github.com/solitonymi/go-snkweb)
[![Build Status](https://travis-ci.org/solitonymi/go-snkweb.svg?branch=master)](https://travis-ci.org/solitonymi/go-snkweb)
[![Coverage Status](https://coveralls.io/repos/github/solitonymi/go-snkweb/badge.svg?branch=master)](https://coveralls.io/github/solitonymi/go-snkweb?branch=master)
[![Go Report Card](https://goreportcard.com/badge/solitonymi/go-snkweb)](https://goreportcard.com/report/solitonymi/go-snkweb)



## Usage

### Login/Access resource/Logout 

```go
  s := &WebAPI{}
  // Login to Soliton NK
  err := s.Login(url, uid, passwd);
  // Create Resource
  r, err := s.CreateResorce("test.txt", "test", false)
  // Upload Resource
  r, err  = s.UploadResorce(r.GUID, []byte("test"))
  // Get Resource List
  list, err := s.GetResorces()
  for _, re := range list {
    ...
  }
  // Delete Resource
  err = s.DeleteResorce(r.GUID)
  
  // Logout
  err := s.Logout()
  // Closck ALl Socket
  s.Close()
```

### Search via Websocket

```go
  s := &WebAPI{}
  err := s.Login(url, uid, passwd)
  defer s.Close()
  err := s.ConnectWebsocket()
  rx, err := s.SendWebsocketCommand(`{"type":"parse","data":{"SearchString":"tag=syslog"}}`, false)
  et := time.Now()
  st := et.Add(-60 * time.Second)
  outsub, err := s.StartSearch(st, et, "tag=syslog length")
  done := false
  count := 0
  for !done {
    done, count, err = s.GetSearchStats(outsub)
    time.Sleep(1 * time.Second)
  }
  if count > 0 {
    resluts, err := s.GetSearchResult(outsub, 0, count)
  }
  err := s.Logout()

```

## Installation

```
$ go get github.com/solitonymi/go-snkweb
```

# License

MIT

# Author

Masayuki Yamai

Soliton Systems K.K 
