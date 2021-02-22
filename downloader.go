package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	command := selectFolder(sfClient, items)

	if command == "x" {
		println("bye bye...")
		os.Exit(0)
	}
	if command == "s" {
		commands := make(map[string]string)
		commands["1"] = "download items of selected folder only"
		commands["2"] = "download items of selected folder only and delete files after successful download"
		commands["3"] = "download items of selected folder and all subfolders"
		commands["4"] = "download items of selected folder and all subfolders and delete files after successful download"
		var command string
		for {
			c, err := selectCommand("select command: ", commands)
			if err != nil {
				println(err.Error())
				continue
			}
			command = c
			break
		}
		switch command {
		case "1":
			downloadItems(sfClient, sfClient.SelectedFolder, false, false)
		case "2":
			downloadItems(sfClient, sfClient.SelectedFolder, true, false)
		case "3":
			downloadItems(sfClient, sfClient.SelectedFolder, false, true)
		case "4":
			downloadItems(sfClient, sfClient.SelectedFolder, true, true)
		}
	}
}

func downloadItems(sfClient *client.SFClient, item *api.FDItem, deleteAfterDownload bool, recursive bool) error {
	items, err := sfClient.LoadFDItems(item)
	if err != nil {
		println("failed to load contents of selected folder >" + err.Error())
		return err
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	downloadRoot := filepath.Join(filepath.Clean(dir), "FD_downloads")
	err = makePath(downloadRoot)
	if err != nil {
		println("failed to create download root path > " + err.Error())
		return err
	}
	for _, v := range items {
		filePathName := filepath.Join(downloadRoot, v.FullPath)
		err := makePath(filepath.Dir(filePathName))
		if err != nil {
			println("failed to create path > " + err.Error())
		}
		if strings.EqualFold(v.ItemType, "file") {
			written, err := sfClient.DownloadFDItem(v, filePathName)
			if err != nil {
				fmt.Printf("failed to download [%v] to [%v] > %v \n", v.Name, filePathName, err.Error())
				continue
			}
			fmt.Printf("downloaded %v bytes to %v\n", written, filePathName)
			if deleteAfterDownload {
				err := sfClient.DeleteFDItem(v)
				if err != nil {
					fmt.Printf("failed to delete [%v] > %v \n", v.Name, err.Error())
				}
			}
		} else if recursive {
			err := makePath(filePathName)
			if err != nil {
				println("failed to create path > " + err.Error())
				continue
			}
			fmt.Printf("created directory %v\n", filePathName)
			// recurse into subdir
			downloadItems(sfClient, &v, deleteAfterDownload, recursive)
		}
	}
	return nil
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

func selectFolder(sfClient *client.SFClient, items []api.FDItem) string {
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
		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("[%v] command: %v\n", k, commands[k])
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
	if sfClient.SelectedFolder == nil {
		return selectedCommand
	}
	return selectedCommand
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

func selectCommand(prompt string, commands map[string]string) (string, error) {
	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("[%v] %v\n", k, commands[k])
	}
	out, err := promptForString(prompt)
	if err != nil {
		return "nil", err
	}
	if _, found := commands[out]; found {
		return out, nil
	}
	return "nil", errors.New("invalid command")
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
