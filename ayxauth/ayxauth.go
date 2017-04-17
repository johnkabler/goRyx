package ayxauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// AyxSigner struct holds params required for building
// requests
type AyxSigner struct {
	ConsumerKey    string
	ConsumerSecret string
	GalleryURL     string
}

// BuildRequest actually generates a request when passed
// ayxSigner struct
func (ayxSigner *AyxSigner) BuildRequest(endpoint string,
	httpMethod string) string {
	signer := *ayxSigner
	consumerKey := signer.ConsumerKey
	consumerSecret := signer.ConsumerSecret
	galleryURL := signer.GalleryURL

	//Build out a clean map for payload params.
	paramsMap := make(map[string]string)

	paramsMap["oauth_consumer_key"] = consumerKey
	paramsMap["oauth_nonce"] = generateNonce()
	paramsMap["oauth_signature_method"] = "HMAC-SHA1"
	paramsMap["oauth_timestamp"] = generateTimestamp()
	paramsMap["oauth_version"] = "1.0"

	//Build request URL
	requestURL := galleryURL + endpoint
	encodedURL := httpMethod + "&" + PercentEncode(requestURL)

	// Concatenate the payload params with key/val
	// in aplphabetical order
	paramsList := make([]string, 0, len(paramsMap))
	for paramKey := range paramsMap {
		paramsList = append(paramsList, paramKey)
	}
	sort.Strings(paramsList)

	var concatParams string
	for pIndex, param := range paramsList {
		if pIndex < (len(paramsMap) - 1) {
			concatParams = concatParams + param + "=" + paramsMap[param] + "&"
		} else {
			concatParams = concatParams + param + "=" + paramsMap[param]
		}
	}
	// Encode the concatenated parameter string
	encodedParams := PercentEncode(concatParams)

	// Create signature base string
	signatureBaseString := encodedURL + "&" + encodedParams

	// Create oauth 1.0 signature
	consumerSecret = PercentEncode(consumerSecret) + "&"
	mac := hmac.New(sha1.New, []byte(consumerSecret))
	mac.Write([]byte(signatureBaseString))
	signatureBytes := mac.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(signatureBytes)
	fmt.Println(signature)
	requestURL = requestURL + "?" + concatParams + "&oauth_signature=" + PercentEncode(signature)
	return (requestURL)
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
