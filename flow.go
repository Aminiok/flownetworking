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
	tools := Tools.New()
	ovnPod := tools.GetOvnPod()
	if ovnPod == "" {
		log.Fatal("OVN pod is not found!")
	}
	if len(os.Args) < 2 {
		tools.PrintHelp()
	}
	switch os.Args[1] {
	case "help":
		tools.PrintHelp()
	case "version":
		tools.GetVersion()
	case "list":
		runListCommand(os.Args[2:], ovnPod)
	case "ls":
		runListCommand(os.Args[2:], ovnPod)
	case "show":
		runShowCommand(os.Args[2:], ovnPod)
	case "sh":
		runShowCommand(os.Args[2:], ovnPod)
	default:
		tools.PrintHelp()
	}
}

func runListCommand(resource []string, ovnPod string) {
	//lp := logicalport.New(ovnPod)
	//lp.GetPortIPDict()
	tools := Tools.New()
	switch resource[0] {
	case "logicalrouters":
		ch := Chassis.New(ovnPod)
		lr := LogicalRouter.New()
		lr.ListLogicalRoutersDetail(ovnPod, resource, ch.GetChassisDict())
	case "lr":
		ch := Chassis.New(ovnPod)
		lr := LogicalRouter.New()
		lr.ListLogicalRoutersDetail(ovnPod, resource, ch.GetChassisDict())
	case "logicalswitch":
		ls := LogicalSwitch.New()
		ls.ListLogicalSwitchDetail(ovnPod, resource)
	case "ls":
		ls := LogicalSwitch.New()
		ls.ListLogicalSwitchDetail(ovnPod, resource)
	case "chassis":
		ch := Chassis.New(ovnPod)
		ch.ListChassisDetail(ovnPod, resource)
	case "ch":
		ch := Chassis.New(ovnPod)
		ch.ListChassisDetail(ovnPod, resource)
	default:
		tools.PrintHelp()
		os.Exit(1)
	}
}

func runShowCommand(resource []string, ovnPod string) {
	lr := LogicalRouter.New()
	tools := Tools.New()
	switch resource[0] {
	case "logicalrouter":
		lr.ShowLogicalRoutersDetail(ovnPod, resource)
	case "lr":
		lr.ShowLogicalRoutersDetail(ovnPod, resource)
	case "chassis":
		lp := LogicalPort.New(ovnPod)
		ch := Chassis.New(ovnPod)
		ch.ShowChassisDetail(ovnPod, resource, lp.GetPortDict())
	case "ch":
		lp := LogicalPort.New(ovnPod)
		ch := Chassis.New(ovnPod)
		ch.ShowChassisDetail(ovnPod, resource, lp.GetPortDict())
	default:
		tools.PrintHelp()
		os.Exit(1)
	}
}
