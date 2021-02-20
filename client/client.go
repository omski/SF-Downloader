package client

import (
	"errors"

	"github.com/omski/SF-Downloader/api"
)

// SFClient
type SFClient struct {
	AuthToken      *string
	SelectedPupil  *api.Pupil
	SelectedFolder *api.FDItem

	Pupils []api.Pupil
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
	sf.Pupils = inventory
	return nil
}

// LoadFDItems
func (sf *SFClient) LoadFDItems(parent *api.FDItem) ([]api.FDItem, error) {
	if sf.AuthToken == nil {
		return nil, errors.New("you must login first")
	}
	if sf.SelectedPupil == nil {
		return nil, errors.New("you must login first")
	}
	parentItemID := "null"
	if parent != nil {
		parentItemID = parent.ID
	}
	items, err := api.LoadFDItems(*sf.AuthToken, parentItemID, *sf.SelectedPupil)
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
	if sf.SelectedPupil == nil {
		return nil, errors.New("you must select something form your inventory first")
	}
	item, err := api.LoadFDItem(*sf.AuthToken, itemID, *sf.SelectedPupil)
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
	if sf.SelectedPupil == nil {
		return written, errors.New("you must select something form your inventory first")
	}
	written, err := api.DownloadFDItem(*sf.AuthToken, *item.ParentItemID, item.ID, filePathName)
	return written, err
}
