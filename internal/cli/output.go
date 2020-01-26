package cli

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type ByName []net.Interface

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func DeterminInterfaceToListenOn() string {
	list := InterfaceList{}
	list.Init()
	chosenInterfaceIndex, err := askForInterface(list)

	if err != nil {
		log.Fatal(err)
	}

	foundNetworkInterface, err := list.Get(chosenInterfaceIndex)

	if err != nil {
		log.Fatal(err)
	}

	return foundNetworkInterface.Address()
}

func printInterfaceList(list InterfaceList, defaultValueIndex int) {
	fmt.Println("Available Interfaces")

	marker := " "

	list.Each(func(number int, networkInterface NetworkInterface) {
		number += 1
		if defaultValueIndex == number {
			marker = "*"
		} else {
			marker = " "
		}

		fmt.Printf("[%2d]%s %20s: %s\n", number, marker, networkInterface.Name, networkInterface.Address())
	})
}

func askForInterface(list InterfaceList) (int, error) {
	localhost := "127.0.0.1"
	defaultChoice, err := list.IndexOfElement(localhost)
	// We use normal indizes 1-x on output
	// And array indizes 0-x-1 on array
	defaultChoice += 1

	if err != nil {
		return 0, fmt.Errorf(`Ask for interface "%v": %w`, localhost, err)
	}

	printInterfaceList(list, defaultChoice)
	timeout := 7

	c := make(chan string, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// "ENTER" == use default value
		if input == "" {
			input = strconv.FormatInt(int64(defaultChoice), 10)
		}
		c <- input
	}()
	fmt.Printf("Enter Number [1-%d] (Timeout: %ds, Default: %d, Press Enter for Default): ", list.Count(), timeout, defaultChoice)
	var choice string

	select {
	case choice = <-c:
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Printf("\nTimeout occured. Choosing number %d\n", defaultChoice)
		choice = strconv.FormatInt(int64(defaultChoice), 10)
	}

	var chosenInterface int

	if _, err := strconv.Atoi(choice); err == nil {
		chosenInterface, _ = strconv.Atoi(choice)

		if chosenInterface > list.Count() || chosenInterface < 1 {
			fmt.Printf("[ERROR] Invalid number. Please use a number between 1 and %d\n", list.Count())
			chosenInterface, _ = askForInterface(list)
		}
	} else {
		fmt.Print(err)
		fmt.Printf("[ERROR] Not a number. Please use a number between 1 and %d\n", list.Count())
		chosenInterface, _ = askForInterface(list)
	}

	chosenInterface -= 1

	return chosenInterface, nil
}
