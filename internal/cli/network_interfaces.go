package cli

import (
	"errors"
	"net"
	"sort"
	"strings"
)

type InvalidNetworkInterfaceError struct {
	InterfaceNumber int
}

func (e InvalidNetworkInterfaceError) Error() string {
	return "Unknown interface number: " + string(e.InterfaceNumber)
}

type NetworkInterface struct {
	Name, Network string
}

func (n *NetworkInterface) Address() (address string) {
	address = strings.Split(n.Network, "/")[0]

	return
}

type InterfaceList struct {
	Interfaces []NetworkInterface
}

func (a *InterfaceList) Len() int { return len(a.Interfaces) }
func (a *InterfaceList) Swap(i, j int) {
	tmp := a.Interfaces[j]
	a.Interfaces[j] = a.Interfaces[i]
	a.Interfaces[i] = tmp
}
func (a *InterfaceList) Less(i, j int) bool { return a.Interfaces[i].Name < a.Interfaces[j].Name }

func (l *InterfaceList) Init() {
	nics, _ := net.Interfaces()
	l.Interfaces = l.buildInterfaceList(nics)
	sort.Sort(l)
}

func (l *InterfaceList) GetInterfaces() []NetworkInterface {
	return l.Interfaces
}

func (l *InterfaceList) Get(interfaceNumber int) (NetworkInterface, error) {
	if interfaceNumber > len(l.Interfaces) {
		return NetworkInterface{}, InvalidNetworkInterfaceError{interfaceNumber}
	}

	if interfaceNumber < 0 {
		return NetworkInterface{}, InvalidNetworkInterfaceError{interfaceNumber}
	}

	return l.Interfaces[interfaceNumber], nil
}

func (l *InterfaceList) buildInterfaceList(nics []net.Interface) (network_interfaces []NetworkInterface) {
	for _, nic := range nics {
		addrs, _ := nic.Addrs()

		for _, addr := range addrs {
			network_interfaces = append(network_interfaces, NetworkInterface{Name: nic.Name, Network: addr.String()})
		}
	}

	if l.ipInNetworkInterfaceSlice("127.0.0.1", network_interfaces) == false {
		network_interfaces = append(network_interfaces, NetworkInterface{Name: "localhost", Network: "127.0.0.1"})
	}

	network_interfaces = l.filterElement([]string{"fe80"}, network_interfaces)

	return network_interfaces
}

func (l *InterfaceList) ipInNetworkInterfaceSlice(a string, list []NetworkInterface) bool {
	for _, b := range list {
		if b.Address() == a {
			return true
		}
	}
	return false
}

func (l *InterfaceList) filterElement(filters []string, list []NetworkInterface) (result []NetworkInterface) {
	for _, b := range list {
		for _, filter := range filters {
			if !strings.HasPrefix(b.Address(), filter) {
				result = append(result, b)
			}
		}
	}

	return
}

func (l *InterfaceList) Count() int {
	return len(l.Interfaces)
}

type fn func(int, NetworkInterface)

func (l *InterfaceList) Each(f fn) {
	for i, networkInterface := range l.Interfaces {
		f(i, networkInterface)
	}
}

func (l *InterfaceList) IndexOfElement(a string) (int, error) {
	number := 0

	for _, b := range l.Interfaces {
		if b.Address() == a {
			return number, nil
		}
		number += 1
	}

	return 0, errors.New("Element not found")
}
