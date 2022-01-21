package Chassis

import (
	"bufio"
	"flownet/LogicalPort"
	"flownet/Tools"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// ChassisList is a list of all Chassis
type ChassisList struct {
	chassisList []chassis
	chassisDict ChassisDict
}

type chassisPort struct {
	portName string
}

type chassis struct {
	name     string
	hostname string
	encap    string
	ip       string
}

type printableChassis struct {
	data       [][]string
	header     []string
	mergedCell []int
}

//PrintableChassis variable
var PrintableChassis = printableChassis{
	data:       [][]string{{}},
	header:     []string{""},
	mergedCell: []int{0},
}

// ChassisDict is a dict of chassis ID IP and Hostname
type ChassisDict struct {
	ChassisHostNameDict map[string]string
	ChassisIPDict       map[string]string
	ChassisIDDict       map[string]string
}

// New chassis object
func New(ovnPod string) ChassisList {
	ch, _ := listChassis(ovnPod)
	return ch
}

func getChassisIDDict(ovnPod string) map[string]string {
	kubeCtlCmd := exec.Command("sudo", "/usr/bin/kubectl", "exec", ovnPod, "ovn-sbctl", "list", "chassis")
	chassisID := ""
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	tools := Tools.New()
	chassisIDDict := make(map[string]string)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) >= 2 {
			if line[0] == "_uuid" {
				chassisID = tools.RefactorString(line[2])
			} else if line[0] == "name" {
				chassisIDDict[chassisID] = tools.RefactorString(line[2])
			}
		}
	}
	return chassisIDDict
}

func listChassis(ovnPod string) (ChassisList, printableChassis) {
	kubeCtlCmd := exec.Command("sudo", "/usr/bin/kubectl", "exec", ovnPod, "ovn-sbctl", "show")
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	tools := Tools.New()
	chassisList := ChassisList{}
	ch := chassis{}
	data := []string{}
	outputData := [][]string{}
	chassisHostnameDict := make(map[string]string)
	chassisIPDict := make(map[string]string)
	chassisDict := ChassisDict{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) == 2 {
			if line[0] == "Chassis" {
				ch.name = tools.RefactorString(line[1])
			} else if line[0] == "hostname:" {
				ch.hostname = tools.RefactorString(line[1])
				chassisHostnameDict[ch.name] = ch.hostname
			} else if line[0] == "Encap" {
				ch.encap = tools.RefactorString(line[1])
			} else if line[0] == "ip:" {
				ch.ip = tools.RefactorString(line[1])
				chassisIPDict[ch.name] = ch.ip
				chassisList.chassisList = append(chassisList.chassisList, ch)
				data = []string{ch.name, ch.hostname, ch.encap, ch.ip}
				outputData = append(outputData, data)
			}
		}
	}
	PrintableChassis.data = outputData
	PrintableChassis.header = []string{"Chassis Name", "Hostname", "Encap", "IP"}
	PrintableChassis.mergedCell = []int{0}
	chassisDict.ChassisHostNameDict = chassisHostnameDict
	chassisDict.ChassisIPDict = chassisIPDict
	chassisDict.ChassisIDDict = getChassisIDDict(ovnPod)
	chassisList.chassisDict = chassisDict
	return chassisList, PrintableChassis
}

func showChassis(chassisName string, ovnPod string, logicalPortDict LogicalPort.PortDict) {
	kubeCtlCmd := exec.Command("sudo", "/usr/bin/kubectl", "exec", ovnPod, "ovn-sbctl", "show")
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	tools := Tools.New()
	ch := chassis{}
	data := []string{}
	outputData := [][]string{}
	isChassisInfo := false
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) == 2 {
			if line[0] == "Chassis" && tools.RefactorString(line[1]) == chassisName {
				ch.name = chassisName
				isChassisInfo = true
			} else if line[0] == "Chassis" && tools.RefactorString(line[1]) != chassisName {
				isChassisInfo = false
			} else if line[0] == "hostname:" && isChassisInfo {
				ch.hostname = tools.RefactorString(line[1])
			} else if line[0] == "Encap" && isChassisInfo {
				ch.encap = tools.RefactorString(line[1])
			} else if line[0] == "ip:" && isChassisInfo {
				ch.ip = tools.RefactorString(line[1])
				data = []string{ch.name, ch.hostname, ch.encap, ch.ip, "", "", ""}
				outputData = append(outputData, data)
			} else if line[0] == "Port_Binding" && isChassisInfo {
				pbID := tools.RefactorString(line[1])
				PortIP := logicalPortDict.PortIPDict[pbID]
				portMAC := logicalPortDict.PortMACDict[pbID]
				data = []string{ch.name, ch.hostname, ch.encap, ch.ip, pbID, PortIP, portMAC}
				outputData = append(outputData, data)
			}
		}
	}
	header := []string{"Chassis Name", "Hostname", "Encap", "IP", "Port ID", "Port IP", "Port MAC"}
	tools.ShowInTable(outputData, header, []int{0, 1, 2, 3})
}

func printListChassis(ovnPod string) {
	tools := Tools.New()
	//listChassis(ovnPod)
	tools.ShowInTable(PrintableChassis.data, PrintableChassis.header, PrintableChassis.mergedCell)
}

// GetChassisDict returns a dict of Chassis IP and Hostname
func (ch *ChassisList) GetChassisDict() ChassisDict {
	return ch.chassisDict
}

// ListChassisDetail executes all ls ch commands
func (ch *ChassisList) ListChassisDetail(ovnPod string, inputParams []string, jsonOutput bool) {
	if len(inputParams) == 1 {
		printListChassis(ovnPod)
	} else {
		fmt.Println("Print Help")
	}
}

// ShowChassisDetail executes all sh ch commands
func (ch *ChassisList) ShowChassisDetail(ovnPod string, inputParams []string, logicalPortDict LogicalPort.PortDict, jsonOutput bool) {
	tools := Tools.New()
	if len(inputParams) == 2 {
		showChassis(inputParams[1], ovnPod, logicalPortDict)
	} else {
		fmt.Println("Command not complete!")
		tools.PrintHelp()
	}
}
