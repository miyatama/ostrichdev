package main

import (
	"errors"
	"flag"
	"log"
	"fmt"
	"time"
	"miyatama/ostrichdev/ostrich"
	"miyatama/ostrichdev/ostrich/web"
	"os"

	"github.com/hashicorp/logutils"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	outputInfo("start ostrich-development")

	// parsing args
	var (
		behavior = flag.String("behavior", "standalone", "standalone or web")
		repository    = flag.String("repository", "", "repository url.ex)https://github.com/xxx/yyy.git")
		fromBranch    = flag.String("from-branch", "", "committed branch name")
		commitId      = flag.String("commit-id", "", "commit id")
		ostrichBranch = flag.String("ostrich-branch", "", "ostrich repository.")
		logLevel      = flag.String("log-level", "WARN", "log level.DEBUG, INFO, WARN, ERROR")
		port          = flag.Int("port", 8080, "ostrich service web port")
	)

	flag.Parse()
	outputInfo(fmt.Sprintf("\tbehavior: %s", *behavior))
	outputInfo(fmt.Sprintf("\trepository: %s", *repository))
	outputInfo(fmt.Sprintf("\tfromBranch: %s", *fromBranch))
	outputInfo(fmt.Sprintf("\tcommitId: %s", *commitId))
	outputInfo(fmt.Sprintf("\tostrichBranch: %s", *ostrichBranch))
	outputInfo(fmt.Sprintf("\tlogLevel: %s", *logLevel))
	outputInfo(fmt.Sprintf("\tport: %d", *port))

	// setting log level
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(*logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	switch(*behavior){
	case "standalone":
		err := callOstrich(*repository , *fromBranch , *ostrichBranch , *commitId)
		if err != nil {
			outputError(err)
		}
		break
	case "web":
		requests := make(chan web.WebRequest)
		go func() {
			for ;; {
				request := <-requests
				switch(request.Action) {
				case web.WebRequestActionOstrich:
					// wait a 3 times
					for i := 0; i < 3; i++ {
						err := callOstrich(
							request.Info.Repository,
							request.Info.FromBranch,
							request.Info.OstrichBranch,
							request.Info.CommitID)
						if err != nil {
							outputError(err)
							time.Sleep(10 * time.Second)
						} else {
							break
						}
					}
				case web.WebRequestActionDone:
					return
				}
			}
		}()

		rest := gin.Default()

		callOstrichWeb := func (c *gin.Context) {
			body := web.OstrichWebRequest{}
			c.Bind(&body)

			requests <- web.WebRequest{
				Action: web.WebRequestActionOstrich,
				Info: body,
			}
			result := web.OstrichWebResponse{}
			status := http.StatusOK
			c.JSON(status, result)
		}
		rest.POST("/ostrich", callOstrichWeb)
		rest.Run(fmt.Sprintf(":%d", *port))
		break
	}
	os.Exit(0)
}

func callOstrich(repository string, fromBranch string, ostrichBranch string, commitId string) error{
	outputInfo(fmt.Sprintf("\trepository: %s", repository))
	outputInfo(fmt.Sprintf("\tfromBranch: %s", fromBranch))
	outputInfo(fmt.Sprintf("\tcommitId: %s", commitId))
	outputInfo(fmt.Sprintf("\tostrichBranch: %s", ostrichBranch))
	if err := HasArgsError(repository, fromBranch, commitId, ostrichBranch); err != nil {
		return err
	}

	ostrich := ostrich.Ostrich{
		Repository:    repository,
		FromBranch:    fromBranch,
		OstrichBranch: ostrichBranch,
		CommitId:      commitId,
		FileAccessor:  &ostrich.FileAccesser{},
	}

	// call ostrich
	if err := ostrich.Run(); err != nil {
		return err
	}
	return nil
}


func outputError(err error) {
	log.Printf("[ERROR]: %s", err.Error())
	log.Printf("[ERROR]: %#v", err)
}
func outputInfo(message string) {
	log.Printf("[INFO]: %s", message)
}

func HasArgsError(repository, fromBrancch, commitID, ostrichBranch string) error {
	if len(repository) <= 0 {
		return errors.New("repository is must need argus")
	}
	if len(fromBrancch) <= 0 {
		return errors.New("from branch is must need argus")
	}
	if len(commitID) <= 0 {
		return errors.New("commit id is must need argus")
	}
	if len(ostrichBranch) <= 0 {
		return errors.New("ostrich branch is must need argus")
	}
	return nil
}
