package gzhu

import (
	"log"

	"github.com/Lofanmi/v3nayou/utils"

	"github.com/joho/godotenv"
	"github.com/parnurzeal/gorequest"
)

func subViewState(value string) string {
	s := utils.StrCut(value, `="__VIEWSTATE" value="`, `"`)
	if s == "" {
		s = utils.StrCut(value, `="__VIEWSTATE" id="__VIEWSTATE" value="`, `"`)
	}
	return s
}

func subViewStateGenerator(value string) string {
	return utils.StrCut(value, `="__VIEWSTATEGENERATOR" value="`, `"`)
}

func subEventValidation(value string) string {
	return utils.StrCut(value, `="__EVENTVALIDATION" value="`, `"`)
}

//--------------------------------------------------------------------------------
// testing
//--------------------------------------------------------------------------------

var r *gorequest.SuperAgent

func getRequest() *gorequest.SuperAgent {
	err := godotenv.Load(`../../.env`)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// os.Setenv("GOREQUEST_DEBUG", "1")
	if r == nil {
		r = gorequest.New().SetDoNotClearSuperAgent(true)
	}
	return r
}
