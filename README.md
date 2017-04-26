# goRyx

# goRyx

goRyx is a command line utility for interacting with the Alteryx Gallery API.

It is very much a work in progress.  

goRyx makes use of ayxfetch, ayxdl, and ayxauth packages which are within this repository.
These packages are also a work in progress and haven't been published outside of
this repository.  

```sh
goryx.exe <alteryx_consumerKey> <alteryx_consumerSecret> <galleryUrl> <endpoint> <file_output_path>
```
	-	alteryx_consumerKey: The key provided by activating Gallery API

	-	alteryx_consumerSecret: The secret provided by activating Gallery API

	-	galleryUrl: This is the URL of the API endpoint.
								For admin use, use
								http://{hostname}/gallery/api/admin/v1/

								For nonAdmin use, use
								http://{hostname}/gallery/api/v1/
								****** NOTE ***********
								NonAdmin users will not be able to connections, and will
								only be able to download workflows that are present in their
								respective studio.

	-	endpoint:  This specifies which function the method the user wants to
							 call on the API.

							 Possible values for "endpoint" parameter:

						workflows :  Fetches list of all workflows and outputs to a CSV
						connections : Fetches list of all data connections
						              (both server and system)
						download :  This will download all of the Alteryx workflows that have
									been deployed to the Gallery.  
									Note that the file_output_path
									parameter needs to be a directory for this function.
