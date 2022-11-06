module github.com/lachlan2k/phatcrack/agent

go 1.18

replace github.com/lachlan2k/phatcrack/common v0.0.0 => ../common
replace gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.6.0

require (
	github.com/gorilla/websocket v1.5.0
	github.com/lachlan2k/phatcrack/common v0.0.0
)

require (
	github.com/hpcloud/tail v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
	gopkg.in/fsnotify.v1 v1.6.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
