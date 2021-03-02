package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/omski/SF-Downloader/api"
)

// StateFileName name of the state file
const StateFileName = ".fd-downloader_state"

// DownloadRoot name of the download root dir
const DownloadRoot = "FD_downloads"

// SFClient
type SFClient struct {
	AuthToken             *string
	SelectedInventoryItem *api.InventoryItem
	SelectedFolder        *api.FDItem
	SelectedCommand       *string

	InventoryItems []api.InventoryItem
}

// Login
func (sf *SFClient) Login(user string, password string) error {
	authToken, err := api.Login(user, password)
	if err != nil {
		return err
	}
	sf.AuthToken = authToken
	return nil
}

// loads the inventory
func (sf *SFClient) LoadInventory() error {
	if sf.AuthToken == nil {
		return errors.New("you must login first")
	}
	inventory, err := api.Inventory(*sf.AuthToken)
	if err != nil {
		return err
	}
	sf.InventoryItems = *inventory
	return nil
}

// LoadFDItems
func (sf *SFClient) LoadFDItems(parent *api.FDItem) ([]api.FDItem, error) {
	if sf.AuthToken == nil {
		return nil, errors.New("you must login first")
	}
	if sf.SelectedInventoryItem == nil {
		return nil, errors.New("you must select something from your inventory first")
	}
	parentItemID := "null"
	if parent != nil {
		parentItemID = parent.ID
	}
	items, err := api.LoadFDItems(*sf.AuthToken, parentItemID, *sf.SelectedInventoryItem)
	if err != nil {
		return nil, err
	}
	return *items, nil
}

// LoadFDItem
func (sf *SFClient) LoadFDItem(itemID string) (*api.FDItem, error) {
	if sf.AuthToken == nil {
		return nil, errors.New("you must login first")
	}
	if sf.SelectedInventoryItem == nil {
		return nil, errors.New("you must select something from your inventory first")
	}
	item, err := api.LoadFDItem(*sf.AuthToken, itemID, *sf.SelectedInventoryItem)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DownloadFDItem downloads a FD item
func (sf *SFClient) DownloadFDItem(item api.FDItem, filePathName string) (int64, error) {
	written := int64(-1)
	if sf.AuthToken == nil {
		return written, errors.New("you must login first")
	}
	if sf.SelectedInventoryItem == nil {
		return written, errors.New("you must select something from your inventory first")
	}
	_, err := os.Stat(filePathName)
	if os.IsNotExist(err) {
		written, err = api.DownloadFDItem(*sf.AuthToken, *item.ParentItemID, item.ID, filePathName)
		return written, err
	}
	return written, nil
}

// DeleteFDItem deletes a FD item
func (sf *SFClient) DeleteFDItem(item api.FDItem) error {
	if sf.AuthToken == nil {
		return errors.New("you must login first")
	}
	err := api.DeleteFDItem(*sf.AuthToken, item.ID)
	return err
}

// SaveState saves the current state of the client
func (sf *SFClient) SaveState() error {
	s, err := json.Marshal(sf)
	if err != nil {
		log.Println(err)
		return err
	}
	println(s)
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	stateFileName := filepath.Join(filepath.Clean(dir), StateFileName)
	out, err := os.Create(stateFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = out.Write(s)
	return err
}

// RestoreState restores a previously saved state
func RestoreState() (*SFClient, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	stateFileName := filepath.Join(filepath.Clean(dir), StateFileName)
	_, err = os.Stat(stateFileName)
	if os.IsNotExist(err) {
		return nil, err
	}
	file, err := os.Open(stateFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	jsonContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var sfClient SFClient
	err = json.Unmarshal(jsonContent, &sfClient)
	if err != nil {
		return nil, err
	}
	return &sfClient, err
}

// DeleteStateFile deletes the state file
func DeleteStateFile() error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	stateFileName := filepath.Join(filepath.Clean(dir), StateFileName)

	_, err = os.Stat(stateFileName)
	if os.IsNotExist(err) {
		return err
	}
	err = os.Remove(stateFileName)
	return err
}
