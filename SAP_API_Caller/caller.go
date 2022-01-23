package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-profit-center-reads/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library/logger"
	"golang.org/x/xerrors"
)

type SAPAPICaller struct {
	baseURL string
	apiKey  string
	log     *logger.Logger
}

func NewSAPAPICaller(baseUrl string, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL: baseUrl,
		apiKey:  GetApiKey(),
		log:     l,
	}
}

func (c *SAPAPICaller) AsyncGetProfitCenter(controllingArea, profitCenter, language, profitCenterName string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "Header":
			func() {
				c.Header(controllingArea, profitCenter)
				wg.Done()
			}()
		case "CompanyCodeAssignment":
			func() {
				c.CompanyCodeAssignment(controllingArea, profitCenter)
				wg.Done()
			}()
		case "ProfitCenterName":
			func() {
				c.ProfitCenterName(language, profitCenterName)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) Header(controllingArea, profitCenter string) {
	headerData, err := c.callProfitCenterSrvAPIRequirementHeader("A_ProfitCenter", controllingArea, profitCenter)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(headerData)

	companyCodeAssignmentData, err := c.callToCompanyCodeAssignment(headerData[0].ToCompanyCodeAssignment)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(companyCodeAssignmentData)

	textData, err := c.callToText(headerData[0].ToText)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(textData)

}

func (c *SAPAPICaller) callProfitCenterSrvAPIRequirementHeader(api, controllingArea, profitCenter string) ([]sap_api_output_formatter.Header, error) {
	url := strings.Join([]string{c.baseURL, "API_PROFITCENTER_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithHeader(req, controllingArea, profitCenter)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToHeader(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) callToCompanyCodeAssignment(url string) ([]sap_api_output_formatter.ToCompanyCodeAssignment, error) {
	req, _ := http.NewRequest("GET", url, nil)
	c.setHeaderAPIKeyAccept(req)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToToCompanyCodeAssignment(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) callToText(url string) ([]sap_api_output_formatter.ToText, error) {
	req, _ := http.NewRequest("GET", url, nil)
	c.setHeaderAPIKeyAccept(req)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToToText(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) CompanyCodeAssignment(controllingArea, profitCenter string) {
	data, err := c.callProfitCenterSrvAPIRequirementCompanyCodeAssignment("A_PrftCtrCompanyCodeAssignment", controllingArea, profitCenter)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)

}

func (c *SAPAPICaller) callProfitCenterSrvAPIRequirementCompanyCodeAssignment(api, controllingArea, profitCenter string) ([]sap_api_output_formatter.CompanyCodeAssignment, error) {
	url := strings.Join([]string{c.baseURL, "API_PROFITCENTER_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithCompanyCodeAssignment(req, controllingArea, profitCenter)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToCompanyCodeAssignment(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) ProfitCenterName(language, profitCenterName string) {
	data, err := c.callProfitCenterSrvAPIRequirementProfitCenterName("A_ProfitCenterText", language, profitCenterName)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)

}

func (c *SAPAPICaller) callProfitCenterSrvAPIRequirementProfitCenterName(api, language, profitCenterName string) ([]sap_api_output_formatter.Text, error) {
	url := strings.Join([]string{c.baseURL, "API_PROFITCENTER_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithProfitCenterName(req, language, profitCenterName)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToText(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithHeader(req *http.Request, controllingArea, profitCenter string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("ControllingArea eq '%s' and ProfitCenter eq '%s'", controllingArea, profitCenter))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithCompanyCodeAssignment(req *http.Request, controllingArea, profitCenter string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("ControllingArea eq '%s' and ProfitCenter eq '%s'", controllingArea, profitCenter))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithProfitCenterName(req *http.Request, language, profitCenterName string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Language eq '%s' and substringof('%s', ProfitCenterName)", language, profitCenterName))
	req.URL.RawQuery = params.Encode()
}
