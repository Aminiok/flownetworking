package LogicalSwitch

import (
	"bufio"
	"flownet/Tools"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type logicalSwitchPort struct {
	portName string
	portID   string
	portMac  string
	portIP   string
	portType string
	portTag  string
}

type logicalSwitch struct {
	switchName      string
	switchPortsName []string
}

func New() logicalSwitch {
	ls := logicalSwitch{}
	return ls
}

func generatePortHostDict(ovnPod string) []logicalSwitch {
	tools := Tools.New()
	kubeCtlCmd := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch")
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	logicalSwitches := []logicalSwitch{}
	ls := logicalSwitch{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) >= 3 {
			if line[0] == "name" {
				ls.switchName = tools.RefactorString(line[2])
			} else if line[0] == "ports" {
				ls.switchPortsName = tools.RefactorStringList(line[2:len(line)])
			}
		} else {
			logicalSwitches = append(logicalSwitches, ls)
		}
	}
	logicalSwitches = append(logicalSwitches, ls)
	return logicalSwitches
}

func getLogicalSwitches(ovnPod string) []logicalSwitch {
	tools := Tools.New()
	kubeCtlCmd := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch")
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	logicalSwitches := []logicalSwitch{}
	ls := logicalSwitch{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) >= 3 {
			if line[0] == "name" {
				ls.switchName = tools.RefactorString(line[2])
			} else if line[0] == "ports" {
				ls.switchPortsName = tools.RefactorStringList(line[2:len(line)])
			}
		} else {
			logicalSwitches = append(logicalSwitches, ls)
		}
	}
	logicalSwitches = append(logicalSwitches, ls)
	return logicalSwitches
}

func getAllLogicalSwitchPorts(ovnPod string) []logicalSwitchPort {
	tools := Tools.New()
	logicalSwitchPortExternalRegex := regexp.MustCompile(".*\"neutron:cidrs\"=\"([0-9./]*)")
	lspList := []logicalSwitchPort{}
	kubeCtlCmd := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch_Port")
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	lsp := logicalSwitchPort{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) >= 3 {
			if line[0] == "name" {
				lsp.portName = tools.RefactorString(line[2])
			} else if line[0] == "_uuid" {
				lsp.portID = tools.RefactorString(line[2])
			} else if line[0] == "addresses" && line[2] != "[router]" && line[2] != "[unknown]" {
				lsp.portMac = tools.RefactorStringList(line[2:len(line)])[0]
			} else if line[0] == "addresses" {
				lsp.portMac = ""
			} else if line[0] == "type" {
				lsp.portType = tools.RefactorString(line[2])
			} else if line[0] == "external_ids" {
				split := strings.Split(logicalSwitchPortExternalRegex.FindString(line[2]), "=")
				if len(split) == 2 {
					lsp.portIP = tools.RefactorString(split[1])
				}
			} else if line[0] == "tag" {
				lsp.portTag = tools.RefactorString(line[2])
			}
		} else if len(line) == 0 {
			lspList = append(lspList, lsp)
			lsp := logicalSwitchPort{}
			if lsp.portName != "" {
				log.Fatal("Error in creation of struct")
			}
		}
	}
	lspList = append(lspList, lsp)
	return lspList
}

func getLogicalSwitchPort(ovnPod string, portName string) logicalSwitchPort {
	tools := Tools.New()
	logicalSwitchPortExternalRegex := regexp.MustCompile(".*\"neutron:cidrs\"=\"([0-9./]*)")
	lsp := logicalSwitchPort{}
	kubeCtlCmd := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch_Port", portName)
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) >= 3 {
			if line[0] == "name" {
				lsp.portName = tools.RefactorString(line[2])
			} else if line[0] == "addresses" && line[2] != "[router]" && line[2] != "[unknown]" {
				lsp.portMac = tools.RefactorStringList(line[2:len(line)])[0]
			} else if line[0] == "type" {
				lsp.portType = tools.RefactorString(line[2])
			} else if line[0] == "external_ids" {
				split := strings.Split(logicalSwitchPortExternalRegex.FindString(line[2]), "=")
				if len(split) == 2 {
					lsp.portIP = tools.RefactorString(split[1])
				}
			} else if line[0] == "tag" {
				lsp.portTag = tools.RefactorString(line[2])
			}
		}
	}
	return lsp
}

func getLogicalSwitchDetail(ovnPod string) {
	tools := Tools.New()
	outputData := [][]string{}
	logicalSwitches := getLogicalSwitches(ovnPod)
	logicalSwitchPortList := getAllLogicalSwitchPorts(ovnPod)
	for _, lSwitch := range logicalSwitches {
		data := []string{}
		for _, switchPort := range logicalSwitchPortList {
			if tools.Contains(lSwitch.switchPortsName, switchPort.portID) {
				data = []string{lSwitch.switchName, switchPort.portIP, switchPort.portMac, switchPort.portType, switchPort.portTag}
				outputData = append(outputData, data)
			}
		}
	}
	header := []string{"Logical Switch Name", "Port IP", "Port MAC", "Port Type", "Port Tag"}
	tools.ShowInTable(outputData, header, []int{0})
}

func (ls logicalSwitch) ListLogicalSwitchDetail(ovnPod string, inputParams []string) {
	if len(inputParams) == 1 {
		getLogicalSwitchDetail(ovnPod)
	} /*else if (len(inputParams) == 2 && inputParams[1] == "routes") || (len(inputParams) == 2 && inputParams[1] == "rt") {
		listLogicalRoutersRoutes(ovnPod)
	} else if (len(inputParams) == 2 && inputParams[1] == "nat") || (len(inputParams) == 2 && inputParams[1] == "nats") {
		listLogicalRoutersNat(ovnPod)
	}*/
}

func (ls logicalSwitch) ShowLogicalSwitchDetail(ovnPod string, inputParams []string) {
	/*	if (len(inputParams) == 3 && inputParams[1] == "routes") || (len(inputParams) == 3 && inputParams[1] == "rt") {
			showLogicalRouterRoutes(inputParams[2], ovnPod)
		} else if (len(inputParams) == 3 && inputParams[1] == "nat") || (len(inputParams) == 3 && inputParams[1] == "nats") {
			showLogicalRouterNat(inputParams[2], ovnPod)
		} else {
			fmt.Println("Command not complete! Print Help")
		}
	*/
}
