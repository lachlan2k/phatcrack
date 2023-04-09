package main

import (
	"log"
	"os"

	"github.com/lachlan2k/phatcrack/api/internal/dbnew"
	"github.com/lachlan2k/phatcrack/api/internal/webserver"
)

func main() {
	if os.Getenv("DB_DSN") == "" {
		log.Fatal("DB_DSN was not specified")
	}

	err := dbnew.Connect(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	// dbnew.RegisterUser("admin", "changeme", "admin")
	// dbnew.RegisterUser("bobby", "changeme", "admin")
	// dbnew.RegisterUser("lachlan", "changeme", "admin")
	// db := dbnew.GetInstance()
	// for i := 0; i < 10000; i++ {
	// dbnew.AddJobStdline("c66e79b7-cee8-4dad-aeb3-8fea50a4b5ac", dbnew.JobStdLineStreamStdout, "foo")
	// }

	// // Push to target_hashes
	// //  update jobs set target_hashes = array_append(target_hashes, 'foo')  where id = 'dfe6c30e-703d-47eb-9f75-ee50b9744cc7';

	// // db.Model
	// j := &dbnew.Job{
	// 	RuntimeData: dbnew.JobRuntimeData{
	// 		Status:        "foo",
	// 		OutputLines:   datatypes.NewJSONSlice([]dbnew.JobRuntimeOutputLine{}),
	// 		StatusUpdates: datatypes.NewJSONSlice([]hashcattypes.HashcatStatus{}),
	// 	},
	// 	HashlistVersion: 2,
	// 	HashcatParams: datatypes.JSONType[hashcattypes.HashcatParams]{
	// 		Data: hashcattypes.HashcatParams{
	// 			AttackMode: 123,
	// 		},
	// 	},
	// 	TargetHashes:    []string{"abc", "def", "ghi"},
	// 	AssignedAgentID: nil,
	// }

	// db.Create(j)
	// db.Save(j)

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if os.Getenv("HC_PATH") == "" {
		log.Printf("HC_PATH was not specified, some API endpoints may not work if hashcat is not in PATH\n")
	}

	err = webserver.Listen(port)
	if err != nil {
		log.Fatalf("couldn't run server: %v", err)
	}
}
