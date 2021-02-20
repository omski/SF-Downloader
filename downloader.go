package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	login(sfClient)
	// Load inventory
	loadInventory(sfClient)
	// select from inventory / pupils
	var items []api.FDItem
	selectPupil(sfClient)
	// Load FD root items
	items = loadFDroot(sfClient)
	// select folder

	item, command := selectFolder(sfClient, items)

	if command == "x" {
		println("bye bye...")
		os.Exit(0)
	}
	if command == "s" {
		items, err := sfClient.LoadFDItems(&item)
		if err != nil {
			println("failed to load contents of selected folder >" + err.Error())

		}
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		downloadRoot := filepath.Join(filepath.Clean(dir), "FD_downloads")
		err = makePath(downloadRoot)
		if err != nil {
			println("failed to create path > " + err.Error())
		}
		for _, v := range items {
			filePathName := filepath.Join(downloadRoot, v.FullPath)
			err = makePath(filepath.Dir(filePathName))
			if err != nil {
				println("failed to create path > " + err.Error())
			}
			written, err := sfClient.DownloadFDItem(v, filePathName)
			if err != nil {
				fmt.Printf("failed to download [%v] to [%v] > %v \n", v.Name, filePathName, err.Error())
				continue
			}
			fmt.Printf("downloaded %v bytes to %v\n", written, filePathName)
		}
	}
}

func makePath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func selectFolder(sfClient *client.SFClient, items []api.FDItem) (api.FDItem, string) {
	commands := make(map[string]string)
	commands["s"] = "select current folder"
	commands["x"] = "exit"

	var selectedCommand string

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
		for i, v := range commands {
			fmt.Printf("[%v] command: %v\n", i, v)
		}
		folderIndex, commandIndex, err := promptForIntInRangeOrCommand(fmt.Sprintf("select folder [%v-%v] or select a command", 0, c), 0, c, commands)
		if err != nil {
			continue
		}

		if commandIndex != "nil" {
			fmt.Printf("selected command: %v \n", commands[commandIndex])
			selectedCommand = commandIndex
			break
		} else {
			fmt.Printf("selected index: %v \n", folderIndex)
		}

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
	return *sfClient.SelectedFolder, selectedCommand
}

func loadFDroot(sfClient *client.SFClient) []api.FDItem {
	for {
		println("loading FD root folder...")
		items, err := sfClient.LoadFDItems(nil)
		if err != nil {
			println("failed to load FD root > " + err.Error())
			continue
		}
		return items
	}
}

func selectPupil(sfClient *client.SFClient) {
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
			println("something went wrong > " + err.Error())
			continue
		}
		sfClient.SelectedPupil = &sfClient.Pupils[pupilsIndex]
		fmt.Printf("selected pupil: %v \n", sfClient.SelectedPupil.Name)
		break
	}
}

func login(sfClient *client.SFClient) {
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
}

func loadInventory(sfClient *client.SFClient) {
	for {
		err := sfClient.LoadInventory()
		if err != nil {
			println("failed to load inventory > " + err.Error())
			continue
		}
		break
	}
}

func promptForIntInRangeOrCommand(prompt string, lowerbound int, upperbound int, commands map[string]string) (int, string, error) {
	out, err := promptForString(prompt)
	if err != nil {
		return 0, "nil", err
	}
	// is it an int
	outIndex, err := strconv.Atoi(out)
	if err == nil {
		if outIndex < lowerbound || outIndex > upperbound {
			return 0, "nil", errors.New("value out of range")
		}
		return outIndex, "nil", nil
	}

	if _, found := commands[out]; found {
		return 0, out, nil
	}
	return 0, "nil", errors.New("invalid command")
}

func promptForIntInRange(prompt string, lowerbound int, upperbound int) (int, error) {
	out, err := promptForInt(prompt)
	if err != nil {
		return out, err
	}
	if out < lowerbound || out > upperbound {
		return out, errors.New("value out of range")
	}
	return out, nil
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
