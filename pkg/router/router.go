package router

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/customErrors"
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/fetchLog"
	"github.com/Aman-Codes/e2e-dashboard-backend/pkg/sendRequest"
	"github.com/gin-gonic/gin"
	"github.com/litmuschaos/litmus-go/pkg/log"
)

type LogsInput struct {
	PipelineId string `json:"pipelineId" binding:"required"`
	JobName    string `json:"jobName" binding:"required"`
	StepNumber string `json:"stepNumber" binding:"required"`
}

func handleError(c *gin.Context, err error) {
	log.Errorf("exiting with error %s", err.Error())
	c.JSON(400, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

func Router() {
	router := gin.Default()
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": customErrors.Success(),
		})
	})
	router.GET("/runs/logs/:id", fetchLog.FetchLogApi)
	router.POST("/logs", func(c *gin.Context) {
		var logsInput LogsInput
		c.BindJSON(&logsInput)
		log.Infof("received parameters for post request are, pipelineId: %s, jobName: %s, stepNumber: %s", logsInput.PipelineId, logsInput.JobName, logsInput.StepNumber)
		fullURLFile := "https://api.github.com/repos/litmuschaos/litmus-e2e/actions/runs/" + logsInput.PipelineId + "/logs"
		err := fetchLog.FetchLog(fullURLFile)
		if err != nil {
			handleError(c, err)
			return
		}
		dir := "./output/" + logsInput.JobName
		log.Infof("start reading dir %s", dir)
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Errorf("unable to read dir, err %v", err)
			handleError(c, err)
			return
		}
		for _, f := range files {
			log.Infof("file name is %s", f.Name())
			if strings.HasPrefix(f.Name(), logsInput.StepNumber+"_") {
				log.Infof("the required file name is %s", f.Name())
				fileName := "./output/" + logsInput.JobName + "/" + f.Name()
				log.Infof("start reading file %s", fileName)
				buf, err := os.ReadFile(fileName)
				if err != nil {
					log.Errorf("unable to read file, err %v", err)
					handleError(c, err)
					return
				}
				c.Data(200, "application/json; charset=utf-8", buf)
				return
			}
		}
		handleError(c, err)
	})
	router.GET("/readfile/:name/:number", func(c *gin.Context) {
		name := filepath.Clean(c.Param("name"))
		number := filepath.Clean(c.Param("number"))
		dir := "./output/" + name
		log.Infof("start reading dir %s", dir)
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Errorf("unable to read dir, err %v", err)
			handleError(c, err)
			return
		}
		for _, f := range files {
			log.Infof("file name is %s", f.Name())
			if strings.HasPrefix(f.Name(), number+"_") {
				log.Infof("the required file name is %s", f.Name())
				fileName := "./output/" + name + "/" + f.Name()
				log.Infof("start reading file %s", fileName)
				buf, err := os.ReadFile(fileName)
				if err != nil {
					log.Errorf("unable to read file, err %v", err)
					handleError(c, err)
					return
				}
				c.Data(200, "application/json; charset=utf-8", buf)
				return
			}
		}
		handleError(c, err)
	})
	router.GET("/repos/:orgName/litmus-e2e/actions/workflows", func(c *gin.Context) {
		orgName := c.Param("orgName")
		sendRequest.SendGetRequestWrapper(c, sendRequest.BaseGitHubUrl+"/repos/"+orgName+"/litmus-e2e/actions/workflows")
	})
	router.GET("/repos/:orgName/litmus-e2e/actions/runs/:pipelineId/jobs", func(c *gin.Context) {
		orgName := c.Param("orgName")
		pipelineId := c.Param("pipelineId")
		sendRequest.SendGetRequestWrapper(c, sendRequest.BaseGitHubUrl+"/repos/"+orgName+"/litmus-e2e/actions/runs/"+pipelineId+"/jobs")
	})
	router.GET("/repos/:orgName/litmus-e2e/actions/runs", func(c *gin.Context) {
		orgName := c.Param("orgName")
		sendRequest.SendGetRequestWrapper(c, sendRequest.BaseGitHubUrl+"/repos/"+orgName+"/litmus-e2e/actions/runs")
	})
	router.GET("/repos/:orgName/litmus-e2e/actions/workflows/:workflowName/runs", func(c *gin.Context) {
		orgName := c.Param("orgName")
		workflowName := c.Param("workflowName")
		sendRequest.SendGetRequestWrapper(c, sendRequest.BaseGitHubUrl+"/repos/"+orgName+"/litmus-e2e/actions/workflows/"+workflowName+"/runs")
	})
	router.GET("/repos/:orgName/litmus-go/commits", func(c *gin.Context) {
		orgName := c.Param("orgName")
		sendRequest.SendGetRequestWrapper(c, sendRequest.BaseGitHubUrl+"/repos/"+orgName+"/litmus-go/commits")
	})
	router.Run(":8080")
}
