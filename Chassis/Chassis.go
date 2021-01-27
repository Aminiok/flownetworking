package Chassis

import (
	"bufio"
	"flownet/Tools"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type chassisPort struct {
	portName string
}

type chassis struct {
	name     string
	hostname string
	encap    string
	ip       string
}

func New() chassis {
	ch := chassis{}
	return ch
}

func listChassis(ovnPod string) {
	kubeCtlCmd := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-sbctl", "show")
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
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) == 2 {
			if line[0] == "Chassis" {
				ch.name = tools.RefactorString(line[1])
			} else if line[0] == "hostname:" {
				ch.hostname = tools.RefactorString(line[1])
			} else if line[0] == "Encap" {
				ch.encap = tools.RefactorString(line[1])
			} else if line[0] == "ip:" {
				ch.ip = tools.RefactorString(line[1])
				data = []string{ch.name, ch.hostname, ch.encap, ch.ip}
				outputData = append(outputData, data)
			}
		}
	}
	header := []string{"Chassis Name", "Hostname", "Encap", "IP"}
	tools.ShowInTable(outputData, header, []int{0})
}

func showChassis(chassisName string, ovnPod string) {
	kubeCtlCmd := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-sbctl", "show")
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
			} else if line[0] == "Port_Binding" && isChassisInfo {
				data = []string{ch.name, ch.hostname, ch.encap, ch.ip, tools.RefactorString(line[1])}
				outputData = append(outputData, data)
			}
		}
	}
	if len(outputData) > 0 {
		header := []string{"Chassis Name", "Hostname", "Encap", "IP", "Ports"}
		tools.ShowInTable(outputData, header, []int{0, 1, 2, 3})
	} else {
		log.Println("Chassis not found!")
		tools.PrintHelp()
	}

}

func (ch chassis) ListChassisDetail(ovnPod string, inputParams []string) {
	if len(inputParams) == 1 {
		listChassis(ovnPod)
	} else {
		fmt.Println("Print Help")
	}
}

func (ch chassis) ShowChassisDetail(ovnPod string, inputParams []string) {
	if len(inputParams) == 2 {
		showChassis(inputParams[1], ovnPod)
	} else {
		fmt.Println("Command not complete! Print Help")
	}
}
