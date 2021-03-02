package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-retryablehttp"
)

// Timeout sets the timeout of the http client
const Timeout = 30 * time.Second

// Retries sets the maximum retry count
const Retries = 3

// RetryDelay sets the delay between retries
const RetryDelay = 10 * time.Second

// Login returns auth token
func Login(user string, password string) (*string, error) {
	_errFmtString := "Login failed : %v"
	_url := "https://api.schoolfox.com/api/Users/login"

	payload := strings.NewReader("{username: \"" + user + "\", password: \"" + password + "\", applicationType: \"SF\"}")

	req, err := retryablehttp.NewRequest("POST", _url, payload)
	if err != nil {
		return nil, fmt.Errorf("Login failed : %v", err)
	}

	addStdHeaders(req)

	_client := getClient()

	res, err := _client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(_errFmtString, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		token := result["token"].(string)
		return &token, nil
	}
	return nil, fmt.Errorf(_errFmtString, fmt.Errorf("HTTP Response status [%v]", res.Status))
}

func getClient() *retryablehttp.Client {
	_client := retryablehttp.NewClient()
	_client.HTTPClient.Timeout = Timeout
	_client.RetryMax = Retries
	_client.RetryWaitMin = RetryDelay
	_client.RetryWaitMax = RetryDelay
	_client.CheckRetry = retryablehttp.ErrorPropagatedRetryPolicy
	log.SetLevel(log.WarnLevel)
	_client.Logger = log.StandardLogger()
	return _client
}

func addStdHeaders(req *retryablehttp.Request) {
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("ZUMO-API-VERSION", "2.0.0")
	req.Header.Add("DNT", "1")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Origin", "https://my.schoolfox.app")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://my.schoolfox.app/")
	req.Header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("Host", "api.schoolfox.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
}

func addAuthHeader(req *retryablehttp.Request, authToken string) {
	req.Header.Add("X-ZUMO-AUTH", authToken)
}

// Inventory loads the users inventory
func Inventory(authToken string) (*[]InventoryItem, error) {
	_errFmtString := "Failed to load inventory : %v"
	_url := "https://api.schoolfox.com/api/Common/Inventory"

	req, err := retryablehttp.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, fmt.Errorf(_errFmtString, err)
	}

	addStdHeaders(req)
	addAuthHeader(req, authToken)

	_client := getClient()

	res, err := _client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(_errFmtString, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		var inventory []InventoryItem
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		err = json.Unmarshal(body, &inventory)
		//fmt.Printf("inventory : %+v", inventory)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		return &inventory, nil
	}
	return nil, fmt.Errorf(_errFmtString, fmt.Errorf("HTTP Response status [%v]", res.Status))
}

// LoadFDItems loads items in a FD folder
func LoadFDItems(authToken string, parentItemID string, pupil InventoryItem) (*[]FDItem, error) {
	_errFmtString := "Failed to load FoxDrive items : %v"

	baseURL := "https://api.schoolfox.com/tables/FoxDriveItems"
	if !strings.EqualFold(parentItemID, "null") {
		parentItemID = "%27" + parentItemID + "%27"
	}
	query := fmt.Sprintf("$count=true&$orderby=ItemType,+Name&$filter=SchoolClassId+eq+%%27%v%%27+and+ParentItemId+eq+%v", pupil.SchoolClassID, parentItemID)
	if strings.EqualFold(pupil.ItemType, "Pupil") {
		query = query + fmt.Sprintf("+and+(PupilId+eq+null+or+PupilId+eq+%%27%v%%27)", pupil.ID)
	}
	query = query + "+and+Deleted+eq+false"
	url := baseURL + "?" + query

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(_errFmtString, err)
	}

	addStdHeaders(req)
	addAuthHeader(req, authToken)

	_client := getClient()

	res, err := _client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(_errFmtString, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		type Response struct {
			Count   int
			Results []FDItem
		}
		var fdResult Response
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		err = json.Unmarshal(body, &fdResult)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		return &fdResult.Results, nil
	}
	return nil, fmt.Errorf(_errFmtString, fmt.Errorf("HTTP Response status [%v]", res.Status))
}

// LoadFDItem loads a single FDItem
func LoadFDItem(authToken string, itemID string, pupil InventoryItem) (*FDItem, error) {
	_errFmtString := "Failed to load FoxDrive item : %v"

	url := fmt.Sprintf("https://api.schoolfox.com/api/FoxDriveItems/%v/Item/%v?pupilId=%v", pupil.SchoolClassID, itemID, pupil.ID)
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(_errFmtString, err)
	}

	addStdHeaders(req)
	addAuthHeader(req, authToken)

	_client := getClient()

	res, err := _client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(_errFmtString, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		var fdResult FDItem
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		err = json.Unmarshal(body, &fdResult)
		if err != nil {
			return nil, fmt.Errorf(_errFmtString, err)
		}
		return &fdResult, nil
	}
	return nil, fmt.Errorf(_errFmtString, fmt.Errorf("HTTP Response status [%v]", res.Status))
}

// DownloadFDItem downloads a FD file
func DownloadFDItem(authToken string, parentItemID string, downloadItemID string, filePathName string) (int64, error) {
	_errFmtString := "Failed to download FoxDrive item : %v"
	written := int64(-1)
	url := fmt.Sprintf("https://api.schoolfox.com/api/FoxDriveItems/%v/DownloadFile/%v", parentItemID, downloadItemID)

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return written, fmt.Errorf(_errFmtString, err)
	}

	addStdHeaders(req)
	addAuthHeader(req, authToken)

	_client := getClient()

	res, err := _client.Do(req)
	if err != nil {
		return written, fmt.Errorf(_errFmtString, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		out, err := os.Create(filePathName)
		if err != nil {
			return written, fmt.Errorf(_errFmtString, err)
		}
		defer out.Close()
		written, err = io.Copy(out, res.Body)
		if err != nil {
			return written, fmt.Errorf(_errFmtString, err)
		}
		return written, err
	}
	return written, fmt.Errorf(_errFmtString, fmt.Errorf("HTTP Response status [%v]", res.Status))
}

// DeleteFDItem delete an FD item
func DeleteFDItem(authToken string, ItemID string) error {
	_errFmtString := "Failed to delete FoxDrive item : %v"

	url := fmt.Sprintf("https://api.schoolfox.com/tables/FoxDriveItems/%v", ItemID)

	req, err := retryablehttp.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf(_errFmtString, err)
	}

	addStdHeaders(req)
	addAuthHeader(req, authToken)

	_client := getClient()

	res, err := _client.Do(req)
	if err != nil {
		return fmt.Errorf(_errFmtString, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	return fmt.Errorf(_errFmtString, fmt.Errorf("HTTP Response status [%v]", res.Status))
}
