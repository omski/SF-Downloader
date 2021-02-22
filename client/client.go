package client

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/omski/SF-Downloader/api"
)

// SFClient
type SFClient struct {
	AuthToken             *string
	SelectedInventoryItem *api.InventoryItem
	SelectedFolder        *api.FDItem

	InventoryItems []api.InventoryItem
}

// Login
func (sf *SFClient) Login(user string, password string) error {
	authToken, err := api.Login(user, password)
	if err != nil {
		return err
	}
	sf.AuthToken = &authToken
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
	sf.InventoryItems = inventory
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
	return items, nil
}

// LoadFDItem
func (sf *SFClient) LoadFDItem(itemID string) (*api.FDItem, error) {
	if sf.AuthToken == nil {
		return nil, errors.New("you must login first")
	}
	if sf.SelectedInventoryItem == nil {
		return nil, errors.New("you must select something form your inventory first")
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
		return written, errors.New("you must select something form your inventory first")
	}
	written, err := api.DownloadFDItem(*sf.AuthToken, *item.ParentItemID, item.ID, filePathName)
	return written, err
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
		log.Print(err)
		return err
	}
	println(s)
	return nil
}
