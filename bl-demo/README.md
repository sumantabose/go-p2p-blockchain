Written by Sumanta Bose, 12 Nov 2018

Current Version : `01-bl-demo.go`

Run : `go run 01-bl-demo.go`

MUX server methods available are:
- http://localhost:port/
- http://localhost:port/info
- http://localhost:port/info/{status}/{member}
- http://localhost:port/add
- http://localhost:port/add/{loop}
- http://localhost:port/move/{serial}
- http://localhost:port/post

FLAGS are:

  -`bldata` *string* : pathname of BL data storage directory (default "bldata")
        
  -`members` *int* : total number of members in the BL supply chain (default 5) // 4 + 1
        
  -`port` *int* : mux server listen port (default 8080)
