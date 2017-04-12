package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

//WorkflowList is a struct that holds the list of workflows
//returned by the Alteryx Server getWorkflowsList API
type WorkflowList []struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscriptionId"`
	Public         bool      `json:"public"`
	RunDisabled    bool      `json:"runDisabled"`
	PackageType    int       `json:"packageType"`
	UploadDate     time.Time `json:"uploadDate"`
	FileName       string    `json:"fileName"`
	MetaInfo       struct {
		Name                 string `json:"name"`
		Description          string `json:"description"`
		Author               string `json:"author"`
		Copyright            string `json:"copyright"`
		URL                  string `json:"url"`
		URLText              string `json:"urlText"`
		OutputMessage        string `json:"outputMessage"`
		NoOutputFilesMessage string `json:"noOutputFilesMessage"`
	} `json:"metaInfo"`
	IsChained bool `json:"isChained"`
	Version   int  `json:"version"`
	RunCount  int  `json:"runCount"`
}

// PercentEncode percent encodes a string according to RFC 3986 2.1.
func PercentEncode(input string) string {
	var buf bytes.Buffer
	for _, b := range []byte(input) {
		// if in unreserved set
		if shouldEscape(b) {
			buf.Write([]byte(fmt.Sprintf("%%%02X", b)))
		} else {
			// do not escape, write byte as-is
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

// shouldEscape returns false if the byte is an unreserved character that
// should not be escaped and true otherwise, according to RFC 3986 2.1.
func shouldEscape(c byte) bool {
	// RFC3986 2.3 unreserved characters
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '.', '_', '~':
		return false
	}
	// all other bytes must be escaped
	return true
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func generateTimestamp() string {
	unixEpochSeconds := time.Now().Unix()
	return strconv.FormatInt(unixEpochSeconds, 10)
}

func main() {
	var endpoint = "workflows/subscription/"
	var httpMethod = "GET"

	if len(os.Args) == 1 {
		fmt.Printf("usage: %s <alteryx_consumerKey> <alteryx_clientKey> <galleryUrl>\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	var consumerKey = os.Args[1]
	var consumerSecret = os.Args[2]
	var galleryURL = os.Args[3]

	paramsMap := make(map[string]string)

	paramsMap["oauth_consumer_key"] = consumerKey
	paramsMap["oauth_nonce"] = generateNonce()
	paramsMap["oauth_signature_method"] = "HMAC-SHA1"
	paramsMap["oauth_timestamp"] = generateTimestamp()
	paramsMap["oauth_version"] = "1.0"
	fmt.Println(paramsMap)
	// Step 1:  Will need to concatenate galleryUrl and
	// endpoint Url together and then encode them (percent)
	var requestURL = galleryURL + endpoint
	fmt.Println(requestURL)
	encodedURL := PercentEncode(requestURL)
	fmt.Println(encodedURL)
	// Step 2: Concatenate the httpMethod with & and then encodedURL
	encodedURL = httpMethod + "&" + encodedURL
	fmt.Println(encodedURL)

	//Step 3: Concatenate the oauth params (in alphebetical order) with key = value
	// and then join them all to each other with &
	paramsList := make([]string, 0, len(paramsMap))
	for paramKey := range paramsMap {
		paramsList = append(paramsList, paramKey)
	}
	sort.Strings(paramsList)

	var concatParams = ""

	for pIndex, param := range paramsList {
		if pIndex < (len(paramsMap) - 1) {
			concatParams = concatParams + param + "=" + paramsMap[param] + "&"
		} else {
			concatParams = concatParams + param + "=" + paramsMap[param]
		}
	}
	fmt.Println(concatParams)
	//Step 4:  % encode this concatenated string of params.
	encodedParams := PercentEncode(concatParams)

	fmt.Println(encodedParams)

	//Step 5: Append the encoded, concatenated params to the encodedURL, separated
	//by &
	signatureBaseString := encodedURL + "&" + encodedParams

	//Step 6: Take the signatureBaseString and consumerSecret and create the
	//signature with HMAC-SHA1
	consumerSecret = PercentEncode(consumerSecret) + "&"
	mac := hmac.New(sha1.New, []byte(consumerSecret))
	mac.Write([]byte(signatureBaseString))
	signatureBytes := mac.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(signatureBytes)
	fmt.Println(signature)
	requestURL = requestURL + "?" + concatParams + "&oauth_signature=" + PercentEncode(signature)

	fmt.Println(requestURL)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
	fmt.Println(resp.Body)

	var workflows WorkflowList

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.NewDecoder(resp.Body).Decode(&workflows)
	//err = json.Unmarshal(stringResponse, &workflows)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(workflows)
	fmt.Println(resp.Status)

}
