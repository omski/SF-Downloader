package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/inancgumus/screen"
	"github.com/omski/SF-Downloader/api"
	"github.com/omski/SF-Downloader/client"
)

func main() {
	screen.Clear()
	sfClient := new(client.SFClient)
	// Login
	for {
		user, _ := promptForString("SF user")
		password, _ := promptForString("SF password")
		err := sfClient.Login(user, password)
		if err != nil {
			println("login failed > " + err.Error())
			continue
		}
		break
	}
	// Load inventory
	err := sfClient.LoadInventory()
	if err != nil {
		println(err.Error())
	}
	// select from inventory
	var items []api.FDItem
	for {
		screen.Clear()
		shownPupils := 0

		for i := 0; i < len(sfClient.Pupils); i++ {
			if sfClient.Pupils[i].ItemType == "School" {
				continue
			}
			fmt.Printf("[%v] = %v, %v \n", shownPupils, sfClient.Pupils[i].Name, sfClient.Pupils[i].SchoolClassName)
			shownPupils++
		}
		pupilsIndex, err := promptForIntInRange(fmt.Sprintf("select pupil [%v-%v]", 0, shownPupils-1), 0, shownPupils-1)
		if err != nil {
			continue
		}
		sfClient.SelectedPupil = &sfClient.Pupils[pupilsIndex]
		fmt.Printf("selected pupil %v \n", sfClient.SelectedPupil.Name)
		break
	}
	// Load FD root items
	println("loading FD root folder...")
	items, err = sfClient.LoadFDItems(nil)
	if err != nil {
		println("failed to load FD root > " + err.Error())
	}
	// select folders
	for {
		screen.Clear()
		hasParent := 0
		if sfClient.SelectedFolder != nil {
			fmt.Println("[0] ..")
			hasParent = 1
		}
		for i := 0; i < len(items); i++ {
			fmt.Printf("[%v] \"%v\" Type:%v/%v Size:%v Access:%v\n", i+hasParent, items[i].FullPath, items[i].ItemType, items[i].ItemSubType, items[i].Size, items[i].AccessType)
		}
		c := len(items)
		if c > 0 {
			if hasParent == 0 {
				c--
			}
		}

		folderIndex, err := promptForIntInRange(fmt.Sprintf("select folder [%v-%v]", 0, c), 0, c)
		if err != nil {
			continue
		}
		fmt.Printf("selected index %v \n", folderIndex)
		if folderIndex == 0 && hasParent == 1 && sfClient.SelectedFolder.ParentItemID == nil {
			sfClient.SelectedFolder = nil
		} else if folderIndex == 0 && hasParent == 1 && sfClient.SelectedFolder.ParentItemID != nil {
			parent, err := sfClient.LoadFDItem(*sfClient.SelectedFolder.ParentItemID)
			if err != nil {
				println(err.Error())
			} else {
				sfClient.SelectedFolder = parent
			}
		} else if folderIndex > 0 && hasParent == 1 {
			sfClient.SelectedFolder = &items[folderIndex-1]
		} else {
			sfClient.SelectedFolder = &items[folderIndex]
		}
		subItems, err := sfClient.LoadFDItems(sfClient.SelectedFolder)
		if err != nil {
			println(err.Error())
		} else {
			items = subItems
		}
	}
}

func promptForIntInRange(prompt string, lowerbound int, upperbound int) (int, error) {
	for {
		out, err := promptForInt(prompt)
		if err != nil {
			println(err.Error())
			return out, err
		}
		if out < lowerbound || out > upperbound {
			println("value out of range")
			return out, err
		}
		return out, nil
	}
}

func promptForInt(prompt string) (int, error) {
	var in string
	var out int
	in, err := promptForString(prompt)
	if err != nil {
		return out, err
	}
	out, err = strconv.Atoi(in)
	return out, err
}

func promptForString(prompt string) (string, error) {
	println(prompt)
	in := bufio.NewReader(os.Stdin)
	input, err := in.ReadString('\n')
	if err != nil {
		err = errors.New("input could not be read")
	}
	return strings.TrimSpace(input), err
}
