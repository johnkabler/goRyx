package ayxfetch

import (
	"encoding/json"
	"fmt"
	"goryx/ayxauth"
	"log"
	"net/http"
	"strings"
)

// FetchWorkflows hits the workflows endpoint of
// the Gallery API and returns a list of the workflows
// present on the server.
func FetchWorkflows(ayxSigner *ayxauth.AyxSigner) {
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
		return
	}
	fmt.Println(workflows)
	fmt.Println(resp.Status)
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
