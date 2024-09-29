module github.com/lachlan2k/phatcrack/agent

go 1.23

replace github.com/lachlan2k/phatcrack/common v0.0.0 => ../common

replace gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.6.0

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/lachlan2k/phatcrack/common v0.0.0
	github.com/nxadm/tail v1.4.11
	golang.org/x/sys v0.25.0
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
