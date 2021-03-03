package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/inancgumus/screen"
	"github.com/omski/SF-Downloader/api"
	"github.com/omski/SF-Downloader/client"
)

func main() {
	screen.Clear()
	sfClient, err := client.RestoreState()
	var restored bool
	if err != nil {
		sfClient = new(client.SFClient)
	} else {
		println("Previous saved state successfully restored.")
		restored = true
	}
	err = sfClient.LoadInventory()
	if restored && err != nil {
		println("Auth token expired...")
	}
	// Login
	if err != nil || sfClient.AuthToken == nil {
		login(sfClient)
		if restored {
			sfClient.SaveState()
		}
	}

	// Load inventory
	if sfClient.SelectedInventoryItem == nil {
		loadInventory(sfClient)
		// select from inventory / pupils
		selectFromInventory(sfClient)
	}

	var items []api.FDItem
	var command string
	if sfClient.SelectedCommand == nil {
		// Load FD root items
		items = loadFDroot(sfClient)
		// select folder
		command = selectFolder(sfClient, items)

		if strings.EqualFold(command, "x") {
			println("bye bye...")
			os.Exit(0)
		}
	} else {
		command = "s"
	}
	if strings.EqualFold(command, "s") {

		endGame := "nil"

		if sfClient.SelectedCommand == nil {
			commands := make(map[string]string)
			commands["1"] = "Download items of selected folder only"
			commands["2"] = "Download items of selected folder only and delete files after successful download"
			commands["3"] = "Download items of selected folder and all subfolders"
			commands["4"] = "Download items of selected folder and all subfolders and delete files after successful download"

			for {
				c, err := selectCommand("Select command: ", commands)
				if err != nil {
					println(err.Error())
					continue
				}
				sfClient.SelectedCommand = &c
				break
			}
		} else {
			command = *sfClient.SelectedCommand
		}
		for {
			screen.Clear()
			var err error
			println("Start download...")
			switch *sfClient.SelectedCommand {
			case "1":
				err = downloadItems(sfClient, sfClient.SelectedFolder, false, false)
			case "2":
				err = downloadItems(sfClient, sfClient.SelectedFolder, true, false)
			case "3":
				err = downloadItems(sfClient, sfClient.SelectedFolder, false, true)
			case "4":
				err = downloadItems(sfClient, sfClient.SelectedFolder, true, true)
			}
			println("Download finished")
			if err != nil {
				println("Something went seriously wrong > " + err.Error())
			} else {
				if strings.EqualFold(endGame, "nil") {
					for {
						commands := make(map[string]string)
						commands["1"] = "Save current settings and exit"
						commands["2"] = "Delete saved state and exit"
						commands["3"] = "Restart command every 15 minutes"
						commands["4"] = "Exit"
						for {
							c, err := selectCommand("Select command: ", commands)
							if err != nil {
								println(err.Error())
								continue
							}
							endGame = c
							break
						}
						break
					}
				}
				switch endGame {
				case "1":
					err = sfClient.SaveState()
					println("bye bye...")
					os.Exit(0)
				case "2":
					err = client.DeleteStateFile()
					if err != nil {
						println("Failed to delete state file > " + err.Error())
					}
					println("bye bye...")
					os.Exit(0)
				case "3":
					time.Sleep(15 * time.Minute)
				case "4":
					println("bye bye...")
					os.Exit(0)
				}
			}
		}

	}
}

func downloadItems(sfClient *client.SFClient, item *api.FDItem, deleteAfterDownload bool, recursive bool) error {
	items, err := sfClient.LoadFDItems(item)
	if err != nil {
		println("Failed to load contents of selected folder > " + err.Error())
		return err
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	downloadRoot := filepath.Join(filepath.Clean(dir), client.DownloadRoot, sfClient.SelectedInventoryItem.Name)
	err = makePath(downloadRoot)
	if err != nil {
		println("Failed to create download path > " + err.Error())
		return err
	}

	for _, v := range items {
		filePathName := filepath.Join(downloadRoot, v.FullPath)
		err := makePath(filepath.Dir(filePathName))
		if err != nil {
			println("Failed to create path > " + err.Error())
		}
		if strings.EqualFold(v.ItemType, "file") {
			written, err := sfClient.DownloadFDItem(v, filePathName)
			if err != nil {
				fmt.Printf("Failed to download [%v] to [%v] > %v \n", v.Name, filePathName, err.Error())
				continue
			}
			if written == -1 {
				fmt.Printf("File %v already exist\n", filePathName)
			} else {
				fmt.Printf("Downloaded %v bytes to %v\n", written, filePathName)
			}
			if deleteAfterDownload && (item != nil && strings.EqualFold(item.ItemSubType, "SubmissionFolder")) && !strings.EqualFold(v.AccessType, "ReadOnly") {
				err := sfClient.DeleteFDItem(v)
				if err != nil {
					fmt.Printf("Failed to delete [%v] > %v \n", v.Name, err.Error())
				} else {
					fmt.Printf("Deleted FoxDrive item %v\n", v.Name)
				}
			}
		} else if recursive {
			err := makePath(filePathName)
			if err != nil {
				println("Failed to create path > " + err.Error())
				continue
			}
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
		fmt.Printf("Created directory %v\n", path)
	}
	return nil
}

func selectFolder(sfClient *client.SFClient, items []api.FDItem) string {
	commands := make(map[string]string)
	commands["s"] = "Select current folder"
	commands["x"] = "Exit"

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
		folderIndex, commandIndex, err := promptForIntInRangeOrCommand(fmt.Sprintf("Select folder [%v-%v] or select a command", 0, c), 0, c, commands)
		if err != nil {
			continue
		}

		if !strings.EqualFold(commandIndex, "nil") {
			fmt.Printf("Selected command: %v \n", commands[commandIndex])
			selectedCommand = commandIndex
			break
		} else {
			fmt.Printf("Selected index: %v \n", folderIndex)
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
		println("Loading FD root folder...")
		items, err := sfClient.LoadFDItems(nil)
		if err != nil {
			println("Failed to load FD root > " + err.Error())
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

func selectFromInventory(sfClient *client.SFClient) {
	for {
		screen.Clear()
		shownPupils := 0
		for i := 0; i < len(sfClient.InventoryItems); i++ {
			if strings.EqualFold(sfClient.InventoryItems[i].ItemType, "School") {
				continue
			}
			fmt.Printf("[%v] = %v, %v \n", shownPupils, sfClient.InventoryItems[i].Name, sfClient.InventoryItems[i].SchoolClassName)
			shownPupils++
		}
		pupilsIndex, err := promptForIntInRange(fmt.Sprintf("Select pupil or class [%v-%v]", 0, shownPupils-1), 0, shownPupils-1)
		if err != nil {
			println("Something went wrong > " + err.Error())
			continue
		}
		sfClient.SelectedInventoryItem = &sfClient.InventoryItems[pupilsIndex]
		fmt.Printf("Selected: %v \n", sfClient.SelectedInventoryItem.Name)
		break
	}
}

func login(sfClient *client.SFClient) {
	for {
		user, _ := promptForString("SF user")
		password, _ := promptForString("SF password")
		err := sfClient.Login(user, password)
		if err != nil {
			println(err.Error())
			continue
		}
		break
	}
}

func loadInventory(sfClient *client.SFClient) {
	for {
		err := sfClient.LoadInventory()
		if err != nil {
			println(err.Error())
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
			return 0, "nil", errors.New("Value out of range")
		}
		return outIndex, "nil", nil
	}

	if _, found := commands[out]; found {
		return 0, out, nil
	}
	return 0, "nil", errors.New("Invalid command")
}

func promptForIntInRange(prompt string, lowerbound int, upperbound int) (int, error) {
	out, err := promptForInt(prompt)
	if err != nil {
		return out, err
	}
	if out < lowerbound || out > upperbound {
		return out, errors.New("Value out of range")
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
		err = errors.New("Input could not be read")
	}
	return strings.TrimSpace(input), err
}
