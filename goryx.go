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
		fmt.Printf("usage: %s <alteryx_consumerKey> <alteryx_clientKey> <galleryUrl>\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	var consumerKey = os.Args[1]
	var consumerSecret = os.Args[2]
	var galleryURL = os.Args[3]

	signer := ayxauth.AyxSigner{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		GalleryURL:     galleryURL}

	ayxfetch.FetchWorkflows(&signer)

	// requestURL := signer.BuildRequest(endpoint, httpMethod)
	//
	// resp, err := http.Get(requestURL)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer resp.Body.Close()
	//
	// fmt.Println(resp)
	// fmt.Println(resp.Body)
	//
	// var workflows WorkflowList
	//
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	//
	// err = json.NewDecoder(resp.Body).Decode(&workflows)
	// //err = json.Unmarshal(stringResponse, &workflows)
	// if err != nil {
	// 	log.Println("ERROR:", err)
	// 	return
	// }
	// fmt.Println(workflows)
	// fmt.Println(resp.Status)

}
