package ayxdl

import (
	"fmt"
	"goryx/ayxauth"
	"goryx/ayxfetch"
	"os"

	"github.com/cavaliercoder/grab"
)

// DownloadWorkflow method creates a signed request for the given appID,
// downloads the file from the Alteryx Gallery API, and writes it to the
// outputPath location.
func DownloadWorkflow(ayxSigner *ayxauth.AyxSigner, appID string, outputPath string) {
	signer := *ayxSigner
	endpoint := appID + `/package/`
	httpMethod := `GET`

	requestURL := signer.BuildRequest(endpoint, httpMethod)

	resp, err := grab.Get(outputPath, requestURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", requestURL, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully downloaded to %s\n", resp.Filename)
}

// DownloadAllWorkflows method will iterate through all of the returned
// AppIDs and download all to the specified output directory
func DownloadAllWorkflows(ayxSigner *ayxauth.AyxSigner, outputDir string) {
	signer := *ayxSigner

	workflowList := ayxfetch.FetchWorkflows(&signer)

	for _, workflow := range *workflowList {
		endpoint := workflow.ID + `/package/`
		httpMethod := `GET`

		requestURL := signer.BuildRequest(endpoint, httpMethod)
		filePath := outputDir + workflow.ID + `.zip`

		resp, err := grab.Get(filePath, requestURL)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", requestURL, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully downloaded to %s\n", resp.Filename)

	}
}
