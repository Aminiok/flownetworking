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
	lsTunnelID      string
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
	logicalSwitchPortConfigRegex := regexp.MustCompile(".*id=\"([0-9]*)\"")
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
			} else if line[0] == "other_config" {
				split := strings.Split(logicalSwitchPortConfigRegex.FindString(line[2]), "=")
				ls.lsTunnelID = tools.RefactorString(split[1])
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
				} else {
					lsp.portIP = ""
				}
			} else if line[0] == "tag" {
				lsp.portTag = tools.RefactorString(line[2])
			}
		} else if len(line) == 0 {
			// remove localnet
			if lsp.portType != "localnet" {
				lspList = append(lspList, lsp)
			}
			lsp := logicalSwitchPort{}
			if lsp.portName != "" {
				log.Fatal("Error in creation of struct")
			}
		}
	}
	if lsp.portType != "localnet" {
		lspList = append(lspList, lsp)
	}
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
				data = []string{lSwitch.switchName, lSwitch.lsTunnelID, switchPort.portIP, switchPort.portMac, switchPort.portType}
				outputData = append(outputData, data)
			}
		}
	}
	header := []string{"Logical Switch Name", "Tunnel ID", "Port IP", "Port MAC", "Port Type"}
	tools.ShowInTable(outputData, header, []int{0, 1})
}

func (ls logicalSwitch) ListLogicalSwitchDetail(ovnPod string, inputParams []string) {
	if len(inputParams) == 1 {
		getLogicalSwitchDetail(ovnPod)
	}
}

func (ls logicalSwitch) ShowLogicalSwitchDetail(ovnPod string, inputParams []string) {
	/*	TBD
	 */
}
