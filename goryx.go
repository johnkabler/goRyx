package main

import (
	"fmt"
	"goryx/ayxauth"
	"goryx/ayxfetch"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf(`usage: %s <alteryx_consumerKey> <alteryx_clientKey> <galleryUrl>
			<endpoint> <file_output_path>

			possible values for "endpoint" parameter:

			workflows :  Fetches list of all workflows.
			connections : Fetches list of all data connections (both server and system)`,
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	var consumerKey = os.Args[1]
	var consumerSecret = os.Args[2]
	var galleryURL = os.Args[3]
	var outputPath = os.Args[5]
	var endpoint = os.Args[4]

	signer := ayxauth.AyxSigner{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		GalleryURL:     galleryURL}

	switch endpoint {
	case "workflows":
		fmt.Println("workflows")
		workflowList := ayxfetch.FetchWorkflows(&signer)
		fmt.Println(workflowList)
		output := []*ayxfetch.WorkflowRecord{}
		ayxfetch.WriteWorkflows(*workflowList, &output, outputPath)

	case "connections":
		fmt.Println("connections")
		connectionsList := ayxfetch.FetchConnections(&signer, "server")
		systemConnectionsList := ayxfetch.FetchConnections(&signer, "system")
		for _, conn := range *systemConnectionsList {
			*connectionsList = append(*connectionsList, conn)
		}
		fmt.Println(connectionsList)
		output := []*ayxfetch.ConnectionRecord{}
		ayxfetch.WriteConnections(*connectionsList, &output, outputPath)
	}
}
