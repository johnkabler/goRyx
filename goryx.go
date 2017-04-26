package main

import (
	"fmt"
	"goryx/ayxauth"
	"goryx/ayxdl"
	"goryx/ayxfetch"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf(`usage: %s <alteryx_consumerKey> <alteryx_consumerSecret> <galleryUrl>
			<endpoint> <file_output_path>

			alteryx_consumerKey: The key provided by activating Gallery API

			alteryx_consumerSecret: The secret provided by activating Gallery API

			galleryUrl: This is the URL of the API endpoint.
									For admin use, use
									http://{hostname}/gallery/api/admin/v1/

									For nonAdmin use, use
									http://{hostname}/gallery/api/v1/
									****** NOTE ***********
									NonAdmin users will not be able to connections, and will
									only be able to download workflows that are present in their
									respective studio.

			endpoint:  This specifies which function the method the user wants to
								 call on the API.

								 Possible values for "endpoint" parameter:

							workflows :  Fetches list of all workflows.
							connections : Fetches list of all data connections (both server and system)
							download :  This will download all of the Alteryx workflows that have
													been deployed to the Gallery.  Note that the file_output_path
													parameter needs to be a directory for this function.`,
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
	case "download":
		fmt.Println("download")
		ayxdl.DownloadAllWorkflows(&signer, outputPath)

	}
}
