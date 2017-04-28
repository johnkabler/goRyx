package ayxdl

import (
	"archive/zip"
	"fmt"
	"goryx/ayxauth"
	"goryx/ayxfetch"
	"io"
	"os"
	"path/filepath"

	"github.com/cavaliercoder/grab"
	"github.com/gocarina/gocsv"
)

// DownloadRecord struct will hold values for each file downloaded
// to ensure we can write a record of all files retrieved to a csv.
type DownloadRecord struct {
	AppID    string `csv:"AppID"`
	FilePath string `csv:"FilePath"`
	FileName string `csv:"FileName"`
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

		unzip(resp.Filename, outputDir+workflow.ID)

		fileList := listFiles(outputDir + workflow.ID)

		for _, fileName := range fileList {
			records = append(records, &DownloadRecord{
				AppID:    workflow.ID,
				FilePath: outputDir + workflow.ID,
				FileName: fileName})
		}

		fmt.Printf("Successfully downloaded to %s\n", resp.Filename)

	}
	return &records
}

// WriteDownloadedFiles will create a CSV which lists out the paths
// to the files that have been downloaded.  This makes for easier integration
// with other tools in Alteryx.
func WriteDownloadedFiles(records *[]*DownloadRecord, outputPath string) {
	outputFile := outputPath + `downloadresults.csv`
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

// Helper functions
func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func listFiles(directory string) []string {
	dirname := directory + string(filepath.Separator)

	d, err := os.Open(dirname)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Reading " + dirname)

	return files
}
