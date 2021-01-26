package LogicalRouter

import (
	"log"
	"os/exec"
	"bufio"
	"strings"
	"strconv"
	"flownet/Tools"
	"fmt"
	"regexp"
)

type logicalRouterPort struct {
	name string
	uuid string
	mac string
	networks string
	redirectChassis string
}

type logicalRouter struct {
	name []string
	nat []string
	ports []string
	staticRoutes []string
}

type logicalRouterDetail struct {
	name string
	logicalRouter logicalRouter
	logicalRouterPorts []logicalRouterPort
}

type logicalRouterRoute struct {
	routerName string
	ipRoute string
	destination string
	origin	string
}

type logicalRouterNat struct {
	routerName string
	natType string
	externalIp string
	logicalIp string
}

func New() (logicalRouter) {
	lr := logicalRouter{}
	return lr
}

func refactorStringList(stringList []string) ([]string){
	newList := []string{}
	for _, item := range  stringList {
		if strings.HasPrefix(item, "[") {
			item = item[1:len(item)]		
		} 
		if strings.HasPrefix(item, "\"") {
			item = item[1:len(item)]		
		}  
		if ( strings.HasSuffix(item, ",") || strings.HasSuffix(item, "]") ) {
			item = item[0:len(item)-1]
		} 
		if strings.HasSuffix(item, "\"") {
			item = item[0:len(item)-1]
		} 
		newList = append (newList, item)
	}
	return newList
}

func refactorString(origString string) (string){
	if strings.HasPrefix(origString, "[") {
		origString = origString[1:len(origString)]		
	} 
	if strings.HasPrefix(origString, "\"") {
		origString = origString[1:len(origString)]		
	}  
	if ( strings.HasSuffix(origString, ",") || strings.HasSuffix(origString, "]") ) {
		origString = origString[0:len(origString)-1]
	} 
	if strings.HasSuffix(origString, "\"") {
		origString = origString[0:len(origString)-1]
	} 
	return origString
}

func getLogicalRouters(ovnPod string) ([]logicalRouter) {
	kubeCtlCmd := exec.Command("/usr/bin/kubectl" , "exec" , ovnPod , "ovn-nbctl" , "list", "Logical_Router")
	out, err := kubeCtlCmd.Output()
    if err != nil {
        log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	logicalRouters := []logicalRouter{}
	lr := logicalRouter{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) >= 3 {
			if line[0] == "name" {
				lr.name = refactorStringList(line[2:len(line)])
			} else if line[0] == "nat" {
				lr.nat = refactorStringList(line[2:len(line)])
			} else if line[0] == "ports" {
				lr.ports = refactorStringList(line[2:len(line)])
			} else if line[0] == "static_routes" {
				lr.staticRoutes = refactorStringList(line[2:len(line)])
			}
		} else {
			logicalRouters = append(logicalRouters, lr)
		}
	}
	logicalRouters = append(logicalRouters, lr)
	return logicalRouters
}

func getLogicalRouterPorts(ovnPod string) ([]logicalRouterPort) {
	kubeCtlCmd := exec.Command("/usr/bin/kubectl" , "exec" , "anc-ovn-0" , "ovn-nbctl" , "list", "Logical_Router_Port")
	redirectChassisRegex := regexp.MustCompile(".*redirect-chassis=\"([A-za-z0-9-]*)")
	//redirectChassisRegex := regexp.MustCompile("([A-za-z0-9-_])")
	out, err := kubeCtlCmd.Output()
    if err != nil {
        log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	logicalRouterPorts := []logicalRouterPort{}
	lri := logicalRouterPort{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) == 3 {
			if line[0] == "_uuid" {
				lri.uuid = refactorString(line[2])
			} else if line[0] == "mac" {
				lri.mac = refactorString(line[2])
			} else if line[0] == "name" {
				lri.name = refactorString(line[2])
			} else if line[0] == "networks" {
				lri.networks = refactorString(line[2])
			} else if (line[0] == "options") {
				if (strings.Contains(line[2], "redirect-chassis=")) {
					split := strings.Split(redirectChassisRegex.FindString(line[2]), "=")
					lri.redirectChassis = refactorString(split[1])
				} else {
					lri.redirectChassis = ""
				}		
			} 
		} else {
			logicalRouterPorts = append(logicalRouterPorts, lri)
			lri := logicalRouterPort{}
			if lri.name != "" {
				log.Fatal("Error in creation of struct")
			}
		}
	}
	logicalRouterPorts = append(logicalRouterPorts, lri)
	return logicalRouterPorts
}

func getLogicalRoutersDetail(ovnPod string) () {
	tools := Tools.New()
	lRouterInterface := getLogicalRouterPorts(ovnPod)
	lRouters := getLogicalRouters(ovnPod)
	logicalRouters := []logicalRouterDetail{}
	outputData := [][]string{}
	for _, router := range lRouters {
		i := 1
		logicalRouterDetail := logicalRouterDetail{}
		logicalRouterDetail.name = router.name[0]
		for _, routerPort := range router.ports {	
			for _, detailedPort := range lRouterInterface {
				if routerPort == detailedPort.uuid {
					logicalRouterDetail.logicalRouterPorts = append(logicalRouterDetail.logicalRouterPorts, detailedPort)
					outputData = append(outputData, []string{router.name[0] , strconv.Itoa(i), detailedPort.mac, detailedPort.networks, detailedPort.redirectChassis})
					i++
				}
			}
		}
		if len(logicalRouterDetail.logicalRouterPorts) == 0 {
			outputData = append(outputData, []string{logicalRouterDetail.name, "", "", "", ""})
		}
		logicalRouters = append(logicalRouters, logicalRouterDetail)
	}
	header := []string{"Router Name", "Port no.", "MAC address", "IP address", "Redirect Chassis"}
	tools.ShowInTable(outputData, header, []int{0})
}

func listLogicalRoutersRoutes(ovnPod string) () {
	name := ""
	outputData := [][]string{}
	tools := Tools.New()
	lRouters := getLogicalRouters(ovnPod)
	logicalRouterRoutes := []logicalRouterRoute{}
	for _, router := range lRouters {
		getRouteCommand := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "lr-route-list", router.name[0])
		out, err := getRouteCommand.Output()
    	if err != nil {
        	log.Fatal(err)
		}
		output := string(out)
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			lrr := logicalRouterRoute{}
			line := strings.Fields(scanner.Text())
			data := []string{"", "", "", ""}
			if len(line) >= 1 {
				if line[0] == "IPv4" {
					name = router.name[0]
				} else {
					lrr.routerName = name
					lrr.ipRoute = line[0]
					lrr.destination = line[1]
					data = []string{lrr.routerName, lrr.ipRoute, lrr.destination, ""}
				} 
				if len(line) > 3 {
					lrr.origin = line[3]
					data[3] = lrr.origin
				} else {
					lrr.origin = ""
					data[3] = ""
				}
			}
			if lrr.routerName != "" {
				logicalRouterRoutes = append(logicalRouterRoutes, lrr)
				outputData = append(outputData, data)
			}
		}	
	}
	header := []string{"Router Name", "Destination", "Next Hop", "Origin"}
	tools.ShowInTable(outputData, header, []int{0})
}

func showLogicalRouterRoutes(routerName string, ovnPod string) () {
	outputData := [][]string{}
	tools := Tools.New()
	getRouteCommand := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "lr-route-list", routerName)
	out, err := getRouteCommand.Output()
    if err != nil {
    	log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	lrr := logicalRouterRoute{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		data := []string{"", "", "", ""}
		if len(line) >= 1 {
			if line[0] != "IPv4" {
				lrr.routerName = routerName
				lrr.ipRoute = line[0]
				lrr.destination = line[1]
				data = []string{lrr.routerName, lrr.ipRoute, lrr.destination, ""}
			} 
			if len(line) > 3 {
				lrr.origin = line[3]
				data[3] = lrr.origin
			} else {
				lrr.origin = ""
				data[3] = ""
			}
		}
		if lrr.routerName != "" {
			outputData = append(outputData, data)
		}
	}
	header := []string{"Router Name", "Destination", "Next Hop", "Origin"}
	tools.ShowInTable(outputData, header, []int{0})
}

func listLogicalRoutersNat(ovnPod string) () {
	outputData := [][]string{}
	tools := Tools.New()
	lRouters := getLogicalRouters(ovnPod)
	logicalRoutersNat := []logicalRouterNat{}
	for _, router := range lRouters {
		getRouteCommand := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "lr-nat-list", router.name[0])
		out, err := getRouteCommand.Output()
    	if err != nil {
        	log.Fatal(err)
		}
		output := string(out)
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			lrn := logicalRouterNat{}
			line := strings.Fields(scanner.Text())
			data := []string{"", "", "", ""}
			if len(line) >= 3 {
				if line[0] != "TYPE" {
					lrn.routerName = router.name[0]
					lrn.natType = line[0]
					lrn.externalIp = line[1]
					lrn.logicalIp = line[2]
					data = []string{lrn.routerName, lrn.natType, lrn.externalIp, lrn.logicalIp}
				} 
			}
			if lrn.routerName != "" {
				logicalRoutersNat = append(logicalRoutersNat, lrn)
				outputData = append(outputData, data)
			}
		}	
	}
	header := []string{"Router Name", "Nat Type", "External IP", "Logical IP"}
	tools.ShowInTable(outputData, header, []int{0})
}

func showLogicalRouterNat(routerName string, ovnPod string) () {
	outputData := [][]string{}
	data := []string{"", "", "", ""}
	tools := Tools.New()
	getRouteCommand := exec.Command("/usr/bin/kubectl", "exec", ovnPod, "ovn-nbctl", "lr-nat-list", routerName)
	out, err := getRouteCommand.Output()
    if err != nil {
        log.Fatal(err)
	}
	output := string(out)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		lrn := logicalRouterNat{}
		line := strings.Fields(scanner.Text())
		data = []string{"", "", "", ""}
		if len(line) >= 3 {
			if line[0] != "TYPE" {
				lrn.routerName = routerName
				lrn.natType = line[0]
				lrn.externalIp = line[1]
				lrn.logicalIp = line[2]
				data = []string{lrn.routerName, lrn.natType, lrn.externalIp, lrn.logicalIp}
			} 
		}
		outputData = append(outputData, data)
	}	
	header := []string{"Router Name", "Nat Type", "External IP", "Logical IP"}
	tools.ShowInTable(outputData, header, []int{0})
}

func (lr logicalRouter) ListLogicalRoutersDetail(ovnPod string, inputParams []string) () {
	if (len(inputParams) == 1){
		getLogicalRoutersDetail(ovnPod)
	} else if (len(inputParams) == 2 && inputParams[1] == "routes") || (len(inputParams) == 2 && inputParams[1] == "rt") {
		listLogicalRoutersRoutes(ovnPod)
	}  else if (len(inputParams) == 2 && inputParams[1] == "nat") || (len(inputParams) == 2 && inputParams[1] == "nats") {
		listLogicalRoutersNat(ovnPod)
	}
}

func (lr logicalRouter) ShowLogicalRoutersDetail(ovnPod string, inputParams []string) () {
	if (len(inputParams) == 3 && inputParams[1] == "routes") || (len(inputParams) == 3 && inputParams[1] == "rt") {
		showLogicalRouterRoutes(inputParams[2], ovnPod)
	} else if (len(inputParams) == 3 && inputParams[1] == "nat") || (len(inputParams) == 3 && inputParams[1] == "nats") {
		showLogicalRouterNat(inputParams[2], ovnPod)
	} else {
		fmt.Println("Command not complete! Print Help")
	}
}