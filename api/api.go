package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Pupil type
type Pupil struct {
	Name              string
	ItemType          string
	SchoolClassID     string
	SchoolClassName   string
	SchoolID          string
	HasUnreadMessages bool
	ID                string
}

type FDItem struct {
	Name                 string
	FullPath             string
	CreatorName          string
	ItemType             string
	ItemSubType          string
	TeachersAccessType   string
	ParentsAccessType    string
	NumberOfParticipants int
	HasPreview           bool
	Size                 int
	ParentItemID         *string
	SchoolClassID        string
	PupilID              string
	AccessType           string
	ID                   string
}

// Login returns auth token
func Login(user string, password string) (string, error) {

	url := "https://api.schoolfox.com/api/Users/login"

	payload := strings.NewReader("{username: \"" + user + "\", password: \"" + password + "\", applicationType: \"SF\"}")

	req, _ := http.NewRequest("POST", url, payload)

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
	req.Header.Add("Content-Length", "79")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	//fmt.Println(res)
	//fmt.Println(string(body))
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result["token"].(string), nil
}

// Inventory loads the inventory
func Inventory(authToken string) ([]Pupil, error) {

	url := "https://api.schoolfox.com/api/Common/Inventory"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("ZUMO-API-VERSION", "2.0.0")
	req.Header.Add("DNT", "1")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	req.Header.Add("X-ZUMO-AUTH", authToken)
	req.Header.Add("Origin", "https://my.schoolfox.app")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://my.schoolfox.app/")
	req.Header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("Host", "api.schoolfox.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var inventory []Pupil
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &inventory)
	//fmt.Printf("inventory : %+v", inventory)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

// LoadFDItems
func LoadFDItems(authToken string, parentItemID string, pupil Pupil) ([]FDItem, error) {
	baseURL := "https://api.schoolfox.com/tables/FoxDriveItems"

	if parentItemID != "null" {
		parentItemID = "%27" + parentItemID + "%27"
	}

	query := fmt.Sprintf("$count=true&$orderby=ItemType,+Name&$filter=SchoolClassId+eq+%%27%v%%27+and+ParentItemId+eq+%v+and+(PupilId+eq+null+or+PupilId+eq+%%27%v%%27)+and+Deleted+eq+false", pupil.SchoolClassID, parentItemID, pupil.ID)

	url := baseURL + "?" + query

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("ZUMO-API-VERSION", "2.0.0")
	req.Header.Add("DNT", "1")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	req.Header.Add("X-ZUMO-AUTH", authToken)
	req.Header.Add("Origin", "https://my.schoolfox.app")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://my.schoolfox.app/")
	req.Header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	type Response struct {
		Count   int
		Results []FDItem
	}
	var fdResult Response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fdResult)
	if err != nil {
		fmt.Println(res)
		fmt.Println(string(body))
		return nil, err
	}
	return fdResult.Results, nil
}

func LoadFDItem(authToken string, itemID string, pupil Pupil) (*FDItem, error) {
	url := fmt.Sprintf("https://api.schoolfox.com/api/FoxDriveItems/%v/Item/%v?pupilId=%v", pupil.SchoolClassID, itemID, pupil.ID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("ZUMO-API-VERSION", "2.0.0")
	req.Header.Add("DNT", "1")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	req.Header.Add("X-ZUMO-AUTH", authToken)
	req.Header.Add("Origin", "https://my.schoolfox.app")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://my.schoolfox.app/")
	req.Header.Add("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var fdResult FDItem
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fdResult)
	if err != nil {
		fmt.Println(res)
		fmt.Println(string(body))
		return nil, err
	}

	return &fdResult, nil
}
