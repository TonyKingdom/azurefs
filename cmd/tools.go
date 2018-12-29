// Copyright © 2018 TonyKindom
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
	"os"
	"strings"
)

func uploadFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	fileInfo, _ := os.Stat(fileName)
	if fileInfo.IsDir() {
		panic("Only files suported!")
	}
	credential, err := azblob.NewSharedKeyCredential(account, key)
	if err != nil {
		fmt.Println("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.chinacloudapi.cn/%s", account, container))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background() // This example uses a never-expiring context
	index := strings.LastIndex(fileName, "/")
	blobName := fileName
	if index != -1 {
		blobName = fileName[index+1:]
	}
	blobUrl := containerURL.NewBlockBlobURL(blobName)
	fmt.Printf("Uploading the file with blob name: %s\n", blobName)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobUrl, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		panic(err)
	}
}

func listFile() {
	credential := azblob.NewAnonymousCredential()
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.chinacloudapi.cn/%s", account, container))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background() // This example uses a never-expiring context
	// List the container that we have created above
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			panic(err)
		}

		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print(blobInfo.Name + "\n")
		}
	}
}

func downloadFile(downloadFile, downLoadDir string) {
	ctx := context.Background() // This example uses a never-expiring context
	credential := azblob.NewAnonymousCredential()
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.chinacloudapi.cn/%s", account, container+"/"+downloadFile))
	blobUrl := azblob.NewBlobURL(*URL, p)
	if _, err := os.Stat(downLoadDir); os.IsNotExist(err) {
		os.MkdirAll(downLoadDir, 755)
	}
	file, _ := os.Create(downLoadDir + "/" + downloadFile)
	azblob.DownloadBlobToFile(ctx, blobUrl, 0, azblob.CountToEnd, file, azblob.DownloadFromBlobOptions{})
}
