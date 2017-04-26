package ayxfetch

import (
	"encoding/json"
	"fmt"
	"goryx/ayxauth"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

// FetchConnections pulls back data from the serverdataconnections
// and returns them in a ConnectionsList struct.
func FetchConnections(ayxSigner *ayxauth.AyxSigner, connectionType string) *ConnectionsList {
	signer := *ayxSigner
	var endpoint string
	switch connectionType {
	case `server`:
		endpoint = `serverdataconnections/`
	case `system`:
		endpoint = `systemdataconnections/`
	}

	httpMethod := `GET`

	requestURL := signer.BuildRequest(endpoint, httpMethod)
	// Send the request
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var connections ConnectionsList

	err = json.NewDecoder(resp.Body).Decode(&connections)

	if err != nil {
		log.Println("ERROR:", err)
	}
	fmt.Println(connections)
	fmt.Println(resp.Status)
	return &connections
}

// WriteConnections method takes the JSON returned from
// Connections APIs and writes to CSV
func WriteConnections(connectionsList ConnectionsList, records *[]*ConnectionRecord, outputPath string) {
	if len(connectionsList) > 0 {
		for _, conn := range connectionsList {
			*records = append(*records, &ConnectionRecord{
				ConnectionID:      conn.ConnectionID,
				ConnectionName:    conn.ConnectionName,
				ConnectionString:  conn.ConnectionString,
				ConnectionType:    conn.ConnectionType,
				SubscriptionCount: conn.SubscriptionCount,
				UserCount:         conn.UserCount})
		}
		csvContent, csvErr := gocsv.MarshalString(records)
		if csvErr != nil {
			panic(csvErr)
		}
		fmt.Println(csvContent)

		outFile, fileErr := os.Create(outputPath)
		if fileErr != nil {
			panic(fileErr)
		}

		writeErr := gocsv.MarshalFile(records, outFile)
		if writeErr != nil {
			panic(writeErr)
		}
	}
}

// ConnectionsList is a struct that holds the list
// of data connections returned by the server and/or
// system data connections API.
type ConnectionsList []struct {
	ConnectionID      string `json:"connectionId"`
	ConnectionName    string `json:"connectionName"`
	ConnectionString  string `json:"connectionString"`
	ConnectionType    string `json:"connectionType"`
	SubscriptionCount int    `json:"subscriptionCount"`
	UserCount         int    `json:"userCount"`
}

// ConnectionRecord is a struct used to hold records
// for writing the connections info to CSV
type ConnectionRecord struct {
	ConnectionID      string `csv:"ConnectionID"`
	ConnectionName    string `csv:"ConnectionName"`
	ConnectionString  string `csv:"ConnectionString"`
	ConnectionType    string `csv:"ConnectionType"`
	SubscriptionCount int    `csv:"SubscriptionCount"`
	UserCount         int    `csv:"UserCount"`
}

// FetchWorkflows hits the workflows endpoint of
// the Gallery API and returns a list of the workflows
// present on the server.
func FetchWorkflows(ayxSigner *ayxauth.AyxSigner) *WorkflowList {
	signer := *ayxSigner
	var endpoint string
	if strings.Contains(signer.GalleryURL, `/admin/v1`) {
		endpoint = `workflows/`
	} else {
		endpoint = `workflows/subscriptions/`
	}
	httpMethod := `GET`

	fmt.Println(endpoint)

	requestURL := signer.BuildRequest(endpoint, httpMethod)

	fmt.Println(requestURL)

	// Send the request
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
	fmt.Println(resp.Body)

	var workflows WorkflowList
	//stringResponse := bytes.TrimPrefix(htmlData, []byte("\xef\xbb\xbf"))
	//resp.Body = bytes.TrimPrefix(resp.Body, []byte("\xef\xbb\xbf"))

	err = json.NewDecoder(resp.Body).Decode(&workflows)

	if err != nil {
		log.Println("ERROR:", err)
	}
	fmt.Println(workflows)
	fmt.Println(resp.Status)
	return &workflows
}

// WriteWorkflows function will iterate through the workflows
// response and write each record out to a slice
// of pointers to records (note that it's a pointer which will
// allow for modifying the slice in place.)
func WriteWorkflows(workflowList WorkflowList, records *[]*WorkflowRecord, outputPath string) {
	if len(workflowList) > 0 {
		for _, workflow := range workflowList {
			var collectionNames []string
			for _, collection := range workflow.Collections {
				collectionNames = append(collectionNames, collection.CollectionName)
			}
			collectionList := strings.Join(collectionNames, ",")
			*records = append(*records, &WorkflowRecord{
				FileName:       workflow.FileName,
				ID:             workflow.ID,
				Author:         workflow.MetaInfo.Author,
				Description:    workflow.MetaInfo.Description,
				Name:           workflow.MetaInfo.Name,
				PackageType:    workflow.PackageType,
				Public:         workflow.Public,
				RunCount:       workflow.RunCount,
				RunDisabled:    workflow.RunDisabled,
				SubscriptionID: workflow.SubscriptionID,
				UploadDate:     workflow.UploadDate,
				Collections:    collectionList,
				Version:        strconv.Itoa(workflow.PublishedVersionNumber)})
		}

		csvContent, csvErr := gocsv.MarshalString(records)
		if csvErr != nil {
			panic(csvErr)
		}
		fmt.Println(csvContent)

		outFile, fileErr := os.Create(outputPath)
		if fileErr != nil {
			panic(fileErr)
		}

		writeErr := gocsv.MarshalFile(records, outFile)
		if writeErr != nil {
			panic(writeErr)
		}
	}
}

//WorkflowList is a struct that holds the list of workflows
//returned by the Alteryx Server getWorkflowsList API
type WorkflowList []struct {
	FileName  string `json:"fileName"`
	ID        string `json:"id"`
	IsChained bool   `json:"isChained"`
	MetaInfo  struct {
		Author               string `json:"author"`
		Copyright            string `json:"copyright"`
		Description          string `json:"description"`
		Name                 string `json:"name"`
		NoOutputFilesMessage string `json:"noOutputFilesMessage"`
		OutputMessage        string `json:"outputMessage"`
		URL                  string `json:"url"`
		URLText              string `json:"urlText"`
	} `json:"metaInfo"`
	PackageType    int         `json:"packageType"`
	Public         bool        `json:"public"`
	RunCount       int         `json:"runCount"`
	RunDisabled    bool        `json:"runDisabled"`
	SubscriptionID string      `json:"subscriptionId"`
	UploadDate     string      `json:"uploadDate"`
	Version        interface{} `json:"version"`
	Collections    []struct {
		CollectionID   string `json:"collectionId"`
		CollectionName string `json:"collectionName"`
	} `json:"collections"`
	LastRunDate            interface{} `json:"lastRunDate"`
	PublishedVersionID     string      `json:"publishedVersionId"`
	PublishedVersionNumber int         `json:"publishedVersionNumber"`
	PublishedVersionOwner  struct {
		Active         bool        `json:"active"`
		Email          string      `json:"email"`
		FirstName      string      `json:"firstName"`
		ID             string      `json:"id"`
		LastName       string      `json:"lastName"`
		SID            interface{} `json:"sId"`
		SubscriptionID string      `json:"subscriptionId"`
	} `json:"publishedVersionOwner"`
	SubscriptionName string `json:"subscriptionName"`
}

// WorkflowRecord struct is used for each row
// written out to the CSV file.
type WorkflowRecord struct {
	FileName       string `csv:"FileName"`
	ID             string `csv:"ID"`
	Author         string `csv:"Author"`
	Description    string `csv:"Description"`
	Name           string `csv:"Name"`
	PackageType    int    `csv:"PackageType"`
	Public         bool   `csv:"Public"`
	RunCount       int    `csv:"runCount"`
	RunDisabled    bool   `csv:"runDisabled"`
	SubscriptionID string `csv:"subscriptionId"`
	UploadDate     string `csv:"uploadDate"`
	Version        string `csv:"version"`
	Collections    string `csv:"collections"`
}
