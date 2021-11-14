package fetchLog

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/constants"
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/customErrors"
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/deleteFolder"
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/env"
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/unzip"
	"github.com/gin-gonic/gin"
	"github.com/litmuschaos/litmus-go/pkg/log"
)

func FetchLog(fullURLFile string) error {
	log.Info("Start to fetch log")
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Errorf("failed to parse URL, err %v", err)
		return customErrors.InternalServerError()
	}
	log.Infof("The parsed fileURL is %v", fileURL)
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1] + ".zip"

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Errorf("failed to create file %s, err %v", fileName, err)
		return customErrors.InternalServerError()
	}
	log.Infof("Successfully created file %s", fileName)
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	req, err := http.NewRequest(http.MethodGet, fullURLFile, http.NoBody)
	if err != nil {
		log.Errorf("failed to create a new http request, err %v", err)
		return customErrors.InternalServerError()
	}
	req.SetBasicAuth(env.GoDotEnvVariable("GITHUB_USERNAME"), env.GoDotEnvVariable("GITHUB_PAT"))
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("failed to execute http request, err %v", err)
		return customErrors.InternalServerError()
	}
	if resp.StatusCode >= 300 {
		log.Errorf("Request failed with status code %d", resp.StatusCode)
		return customErrors.NonSuccessStatusCode(resp.StatusCode)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Errorf("failed to copy content from response to file, err %v", err)
		return customErrors.InternalServerError()
	}
	defer file.Close()
	deleteFolder.DeleteFolder(constants.OutputFolderPath)
	err = unzip.Unzip(fileName)
	if err != nil {
		log.Errorf("failed to unzip file %s", fileName)
		return customErrors.InternalServerError()
	}
	log.Infof("Downloaded a file %s with size %d", fileName, size)
	log.Info("Successfullly fetched log")
	return nil
}

func FetchLogApi(c *gin.Context) {
	id := c.Param("id")
	fullURLFile := "https://api.github.com/repos/litmuschaos/litmus-e2e/actions/runs/" + id + "/logs"
	err := FetchLog(fullURLFile)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": customErrors.Success(),
		"id":     id,
	})
}
