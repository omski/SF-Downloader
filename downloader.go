package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"schoolfox/api"
	"strconv"
	"strings"
)

func main() {
	sfClient := new(api.SchoolfoxClient)

	for {
		user, _ := promptForString("Schoolfox user")
		password, _ := promptForString("Schoolfox password")
		err := sfClient.Login(user, password)
		if err != nil {
			println("login failed > " + err.Error())
			continue
		}
		break
	}

	err := sfClient.LoadInventory()
	if err != nil {
		println(err.Error())
	}

	shownPupils := 0
	for i := 0; i < len(sfClient.Pupils); i++ {
		if sfClient.Pupils[i].ItemType == "School" {
			continue
		}
		fmt.Printf("[%v] = %v, %v \n", shownPupils, sfClient.Pupils[i].Name, sfClient.Pupils[i].SchoolClassName)
		shownPupils++
	}
	pupilsIndex := promptForIntInRange(fmt.Sprintf("select pupil [%v-%v]", 0, shownPupils-1), 0, shownPupils-1)
	fmt.Printf("selected index %v \n", pupilsIndex)
	sfClient.SelectedPupil = &sfClient.Pupils[pupilsIndex]

	fmt.Printf("selected %v\n", sfClient.SelectedPupil.Name)

	items, err := sfClient.LoadFoxDriveItems(nil)
	if err != nil {
		println(err.Error())
	}
	for {
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

		folderIndex := promptForIntInRange(fmt.Sprintf("select folder [%v-%v]", 0, c), 0, c)
		fmt.Printf("selected index %v \n", folderIndex)
		if folderIndex == 0 && hasParent == 1 && sfClient.SelectedFolder.ParentItemID == nil {
			sfClient.SelectedFolder = nil
		} else if folderIndex == 0 && hasParent == 1 && sfClient.SelectedFolder.ParentItemID != nil {
			parent, err := sfClient.LoadFoxDriveItem(*sfClient.SelectedFolder.ParentItemID)
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
		subItems, err := sfClient.LoadFoxDriveItems(sfClient.SelectedFolder)
		if err != nil {
			println(err.Error())
		} else {
			items = subItems
		}
	}
}

func promptForIntInRange(prompt string, lowerbound int, upperbound int) int {
	for {
		out, err := promptForInt(prompt)
		if err != nil {
			println(err.Error())
			continue
		}
		if out < lowerbound || out > upperbound {
			println("value out of range")
			continue
		}
		return out
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
