package main

import (
	"flownet/Chassis"
	"flownet/LogicalPort"
	"flownet/LogicalRouter"
	"flownet/LogicalSwitch"
	"flownet/Tools"
	"log"
	"os"
)

func main() {
	jsonOutput := false
	tools := Tools.New()
	ovnPod := tools.GetOvnPod()
	if ovnPod == "" {
		log.Fatal("OVN pod is not found!")
	}
	if len(os.Args) < 2 {
		tools.PrintHelp()
	}
	if os.Args[len(os.Args)-1] == "--json" {
		jsonOutput = true
	}

	switch os.Args[1] {
	case "help":
		tools.PrintHelp()
	case "version":
		tools.GetVersion()
	case "list", "ls":
		runListCommand(os.Args[2:], ovnPod, jsonOutput)
	case "show", "sh":
		runShowCommand(os.Args[2:], ovnPod, jsonOutput)
	default:
		tools.PrintHelp()
	}
}

func runListCommand(resource []string, ovnPod string, jsonOutput bool) {
	tools := Tools.New()
	switch resource[0] {
	case "logicalrouters", "lr", "vpc":
		ch := Chassis.New(ovnPod)
		lr := LogicalRouter.New()
		lr.ListLogicalRoutersDetail(ovnPod, resource, ch.GetChassisDict(), jsonOutput)
	case "logicalswitch", "ls", "network", "net":
		lp := LogicalPort.New(ovnPod)
		ls := LogicalSwitch.New()
		ls.ListLogicalSwitchDetail(ovnPod, resource, lp.GetPortDict().PortMACDict, jsonOutput)
	case "chassis", "ch":
		ch := Chassis.New(ovnPod)
		ch.ListChassisDetail(ovnPod, resource, jsonOutput)
	case "port", "ports":
		ch := Chassis.New(ovnPod)
		lp := LogicalPort.New(ovnPod)
		lp.ListPortsDetail(ch.GetChassisDict().ChassisIDDict, jsonOutput)
	default:
		tools.PrintHelp()
		os.Exit(1)
	}
}

func runShowCommand(resource []string, ovnPod string, jsonOutput bool) {
	lr := LogicalRouter.New()
	tools := Tools.New()
	switch resource[0] {
	case "logicalrouter", "lr":
		lr.ShowLogicalRoutersDetail(ovnPod, resource, jsonOutput)
	case "chassis", "ch":
		lp := LogicalPort.New(ovnPod)
		ch := Chassis.New(ovnPod)
		ch.ShowChassisDetail(ovnPod, resource, lp.GetPortDict(), jsonOutput)
	default:
		tools.PrintHelp()
		os.Exit(1)
	}
}
