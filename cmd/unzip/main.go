package main

import (
	"os"
	"log"
	"strings"
	"path/filepath"

	"github.com/samkreter/DSA-Workshop/util/extract"
)

var (
	defaultFileURL = "https://ed-public-download.app.cloud.gov/downloads/CollegeScorecard_Raw_Data.zip"
	defaultFilePath = "./CollegeScorecard_Raw_Data.zip"
)

func main(){
	fileURL, ok := os.LookupEnv("FILE_URL")
	if (!ok){
		log.Println("FILE_URL not set. Using default.")
		fileURL = defaultFileURL
	}

	filePath, ok := os.LookupEnv("FILE_PATH")
	if (!ok){
		log.Println("FILE_PATH not set. Using default.")
		filePath = defaultFilePath
	}

	unzipedDest, ok := os.LookupEnv("UNZIPED_DEST")
	if (!ok){
		unzipedDest = strings.TrimSuffix(filePath, filepath.Ext(filePath))
	}

	err := extract.DownloadAndUnzip(fileURL, filePath, unzipedDest)
	if err != nil {
		log.Fatal(err)
	}
}