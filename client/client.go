package client

import (
	"errors"
	"sf-downloader/api"
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
		return nil, errors.New("you must login first")
	}
	item, err := api.LoadFDItem(*sf.AuthToken, itemID, *sf.SelectedPupil)
	if err != nil {
		return nil, err
	}
	return item, nil
}
