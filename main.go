package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func deleteAll(folderName string) {
	fmt.Println("deleting folder ", folderName)
	err := os.RemoveAll(folderName)
	if err != nil {
		log.Fatal(err)
	}
}

func unzip(fileName string) {
	dst := "output"
	archive, err := zip.OpenReader(fileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

func fetchLog(fullURLFile string) {
	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	fileName += ".zip"

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	req, err := http.NewRequest(http.MethodGet, fullURLFile, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("karan1024x", "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")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	deleteAll("./output")
	unzip(fileName)

	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
}

func main() {

	router := gin.Default()
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	router.GET("/runs/logs/:id", func(c *gin.Context) {
		id := c.Param("id")
		fullURLFile := "https://api.github.com/repos/litmuschaos/litmus-e2e/actions/runs/"
		fullURLFile += id
		fullURLFile += "/logs"
		fetchLog(fullURLFile)
		c.JSON(200, gin.H{
			"status": "OK",
			"id":     id,
		})
	})
	router.Run(":8080")
}
