package LogicalPort

import (
	"bufio"
	"flownet/Tools"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// List of LogicalPort maping
type LogicalPortList struct {
	portList []logicalPort
	portDict PortDict
}

// PortDict is a dict containing PortID: PortIP | portID: PortMAC
type PortDict struct {
	PortIPDict  map[string]string
	PortMACDict map[string]string
}

type logicalPort struct {
	portID         string
	portIP         string
	portMac        string
	gatewayChassis string
	portType       string
	portChassis    string
	portLPName     string
}

// New logical port object
func New(ovnPod string) LogicalPortList {
	lp := getLogicalPortList(ovnPod)
	return lp
}

func getLogicalPortList(ovnPod string) LogicalPortList {
	tools := Tools.New()
	kubeCtlCmd := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-sbctl", "list", "port_binding")
	logicalPortExternalRegex := regexp.MustCompile(".*\"neutron:cidrs\"=\"([0-9./]*)")
	out, err := kubeCtlCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	logicalPortList := LogicalPortList{}
	lp := logicalPort{}
	portIPDict := make(map[string]string)
	portMACDict := make(map[string]string)
	portDict := PortDict{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) >= 3 {
			if line[0] == "_uuid" {
				lp.portID = tools.RefactorString(line[2])
			} else if line[0] == "chassis" {
				lp.portChassis = tools.RefactorString(line[2])
			} else if line[0] == "gateway_chassis" {
				lp.gatewayChassis = tools.RefactorString(line[2])
			} else if line[0] == "mac" {
				netInfo := tools.RefactorStringList(line[2:len(line)])
				lp.portMac = netInfo[0]
				if len(netInfo) == 2 {
					lp.portIP = netInfo[1]
				}
			} else if line[0] == "external_ids" {
				split := strings.Split(logicalPortExternalRegex.FindString(line[2]), "=")
				if len(split) == 2 {
					lp.portIP = tools.RefactorString(split[1])
				}
			} else if line[0] == "type" {
				lp.portType = tools.RefactorString(line[2])
			} else if line[0] == "logical_port" {
				lp.portLPName = tools.RefactorString(line[2])
			}
		} else {
			portIPDict[lp.portLPName] = lp.portIP
			portMACDict[lp.portLPName] = lp.portMac
			logicalPortList.portList = append(logicalPortList.portList, lp)
		}
	}
	portIPDict[lp.portLPName] = lp.portIP
	portMACDict[lp.portLPName] = lp.portMac
	portDict.PortIPDict = portIPDict
	portDict.PortMACDict = portMACDict
	logicalPortList.portList = append(logicalPortList.portList, lp)
	logicalPortList.portDict = portDict
	return logicalPortList
}

// GetPortDict returns a Dict of PortID: PortIP
func (lps *LogicalPortList) GetPortDict() PortDict {
	return lps.portDict
}

// ShowLogicalPort shows port list
func (lps *LogicalPortList) ShowLogicalPort() {
	fmt.Println(lps)
}
