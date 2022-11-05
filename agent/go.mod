module github.com/lachlan2k/phatcrack/agent

go 1.18

replace github.com/lachlan2k/phatcrack/common v0.0.0 => ../common

require (
	github.com/gorilla/websocket v1.5.0
	github.com/lachlan2k/phatcrack/common v0.0.0
)
