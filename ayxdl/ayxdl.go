package ayxdl

import (
	"fmt"
	"goryx/ayxauth"
	"goryx/ayxfetch"
	"os"

	"github.com/cavaliercoder/grab"
	"github.com/gocarina/gocsv"
)

// DownloadRecord struct will hold values for each file downloaded
// to ensure we can write a record of all files retrieved to a csv.
type DownloadRecord struct {
	FilePath string `csv:"FilePath"`
}

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
func DownloadAllWorkflows(ayxSigner *ayxauth.AyxSigner, outputDir string) *[]*DownloadRecord {
	signer := *ayxSigner

	workflowList := ayxfetch.FetchWorkflows(&signer)

	records := []*DownloadRecord{}

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

		records = append(records, &DownloadRecord{
			FilePath: resp.Filename})

		fmt.Printf("Successfully downloaded to %s\n", resp.Filename)

	}
	return &records
}

// WriteDownloadedFiles will create a CSV which lists out the paths
// to the files that have been downloaded.  This makes for easier integration
// with other tools in Alteryx.
func WriteDownloadedFiles(records *[]*DownloadRecord) {
	outputFile := `.\results.csv`
	csvContent, csvErr := gocsv.MarshalString(records)
	if csvErr != nil {
		panic(csvErr)
	}
	fmt.Println(csvContent)

	outFile, fileErr := os.Create(outputFile)
	if fileErr != nil {
		panic(fileErr)
	}

	writeErr := gocsv.MarshalFile(records, outFile)
	if writeErr != nil {
		panic(writeErr)
	}
}
