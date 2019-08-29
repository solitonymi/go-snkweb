# go-snkweb

Soliton NK Web API for GO Lang

[![Godoc Reference](https://godoc.org/github.com/mattn/go-colorable?status.svg)](http://godoc.org/github.com/mattn/go-colorable)
[![Build Status](https://travis-ci.org/mattn/go-colorable.svg?branch=master)](https://travis-ci.org/mattn/go-colorable)
[![Coverage Status](https://coveralls.io/repos/github/mattn/go-colorable/badge.svg?branch=master)](https://coveralls.io/github/mattn/go-colorable?branch=master)
[![Go Report Card](https://goreportcard.com/badge/mattn/go-colorable)](https://goreportcard.com/report/mattn/go-colorable)



## Usage

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
  
  // Connect WebSocket
  err = s.ConnectWebsocket()
  // Send WebScokcer Command
  rx, err := s.WebSockCommand(` {"type":"parse","data":{"SearchString":"tag=syslog"}}`, false)
  // Logout
  err := s.Logout()
  // Closck ALl Socket
  s.Close()
```

You can compile above code on non-windows OSs.

## Installation

```
$ go get github.com/solitonymi/snkweb
```

# License

MIT

# Author

Masayuki Yamai
Soliton Systems K.K 
