package LogicalSwitch

import (
	"bufio"
	"flownet/LogicalRouter"
	"flownet/Tools"
	"fmt"
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
	switchName      string `json:"network_name"`
	switchPortsName []string
	lsTunnelID      string
	networkCIDR     string
}

//Port is the port information used for JSON printing
type Port struct {
	PortIP   string `json:"port_ip"`
	PortMAC  string `json:"port_mac"`
	PortType string `json:"port_type"`
}

// Network is the network information used for JSON printing
type Network struct {
	NetworkName     string `json:"network_name"`
	NetworkTunnelID string `json:"tunnel_id"`
	NetworkCIDR     string `json:"network_cidr"`
	NetworkPorts    []Port `json:"ports"`
}

// Networks is the list of networks and its ports used for JSON printing
type Networks struct {
	NetworkList []Network `json:"networks"`
}

func New() logicalSwitch {
	ls := logicalSwitch{}
	return ls
}

func generatePortHostDict(ovnPod string) []logicalSwitch {
	tools := Tools.New()
	kubeCtlCmd := exec.Command("sudo", "/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch")
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
	kubeCtlCmd := exec.Command("sudo", "/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch")
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

func getAllLogicalSwitchPorts(ovnPod string, portMACDict map[string]string) []logicalSwitchPort {
	tools := Tools.New()
	logicalSwitchPortExternalRegex := regexp.MustCompile(".*\"neutron:cidrs\"=\"([0-9./]*)")
	lspList := []logicalSwitchPort{}
	kubeCtlCmd := exec.Command("sudo", "/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch_Port")
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
			// exclude localnet
			if lsp.portType != "localnet" {
				portName := ""
				if lsp.portMac == "" {
					if strings.HasPrefix(lsp.portName, "router") {
						portName = "lrp-" + lsp.portName
					} else if strings.HasPrefix(lsp.portName, "ext_gw") {
						portName = "lrp-" + lsp.portName
					}
					lsp.portMac = portMACDict[portName]
				}
				lspList = append(lspList, lsp)
			}
			lsp := logicalSwitchPort{}
			if lsp.portName != "" {
				log.Fatal("Error in creation of struct")
			}
		}
	}
	if lsp.portType != "localnet" {
		portName := ""
		if lsp.portMac == "" {
			if strings.HasPrefix(lsp.portName, "router") {
				portName = "lrp-" + lsp.portName
			} else if strings.HasPrefix(lsp.portName, "ext_gw") {
				portName = "lrp-" + lsp.portName
			}
			lsp.portMac = portMACDict[portName]
		}
		lspList = append(lspList, lsp)
	}
	return lspList
}

func getLogicalSwitchPort(ovnPod string, portName string) logicalSwitchPort {
	tools := Tools.New()
	logicalSwitchPortExternalRegex := regexp.MustCompile(".*\"neutron:cidrs\"=\"([0-9./]*)")
	lsp := logicalSwitchPort{}
	kubeCtlCmd := exec.Command("sudo", "/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "list", "Logical_Switch_Port", portName)
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
func printLogicalSwitchTable(logicalSwitches []logicalSwitch, logicalSwitchPortList []logicalSwitchPort) {
	tools := Tools.New()
	outputData := [][]string{}
	for _, lSwitch := range logicalSwitches {
		data := []string{}
		for _, switchPort := range logicalSwitchPortList {
			if tools.Contains(lSwitch.switchPortsName, switchPort.portID) {
				data = []string{lSwitch.switchName, lSwitch.lsTunnelID, tools.GetNetworkFromIP(switchPort.portIP), switchPort.portIP, switchPort.portMac, switchPort.portType}
				outputData = append(outputData, data)
			}
		}
	}
	header := []string{"Logical Switch Name", "Tunnel ID", "Network CIDR", "Port IP", "Port MAC", "Port Type"}
	tools.ShowInTable(outputData, header, []int{0, 1, 2})
}

func printLogicalSwitchJSON(logicalSwitches []logicalSwitch, logicalSwitchPortList []logicalSwitchPort) {
	tools := Tools.New()
	networkList := Networks{}
	for _, lSwitch := range logicalSwitches {
		network := Network{NetworkName: lSwitch.switchName, NetworkTunnelID: lSwitch.lsTunnelID}
		for _, switchPort := range logicalSwitchPortList {
			if tools.Contains(lSwitch.switchPortsName, switchPort.portID) {
				network.NetworkCIDR = tools.GetNetworkFromIP(switchPort.portIP)
				port := Port{PortIP: switchPort.portIP, PortMAC: switchPort.portMac, PortType: switchPort.portType}
				network.NetworkPorts = append(network.NetworkPorts, port)
			}
		}
		networkList.NetworkList = append(networkList.NetworkList, network)
	}
	tools.PrintInJSON(networkList)
}

func getLogicalSwitchDetail(ovnPod string, portMACDict map[string]string, jsonOutput bool) {
	logicalSwitches := getLogicalSwitches(ovnPod)
	logicalSwitchPortList := getAllLogicalSwitchPorts(ovnPod, portMACDict)
	if jsonOutput {
		printLogicalSwitchJSON(logicalSwitches, logicalSwitchPortList)

	} else {
		printLogicalSwitchTable(logicalSwitches, logicalSwitchPortList)
	}
}

func (ls logicalSwitch) ListLogicalSwitchDetail(ovnPod string, inputParams []string, portMACDict map[string]string, jsonOutput bool) {
	if len(inputParams) == 2 && inputParams[1] == "detail" {
		lr := LogicalRouter.New()
		_ = lr.GetLogicalRoutersDetail(ovnPod)
		routerDict := lr.GetRouterDict()
		fmt.Println(routerDict.PortMACRouterDict)
	} else {
		getLogicalSwitchDetail(ovnPod, portMACDict, jsonOutput)
	}
}

func (ls logicalSwitch) ShowLogicalSwitchDetail(ovnPod string, inputParams []string, jsonOutput bool) {
	/*	TBD
	 */
}
