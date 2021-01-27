package main

import (
	"flownet/Chassis"
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
	lr := LogicalRouter.New()
	ls := LogicalSwitch.New()
	ch := Chassis.New()
	tools := Tools.New()
	switch resource[0] {
	case "logicalrouters":
		lr.ListLogicalRoutersDetail(ovnPod, resource)
	case "lr":
		lr.ListLogicalRoutersDetail(ovnPod, resource)
	case "logicalswitch":
		ls.ListLogicalSwitchDetail(ovnPod, resource)
	case "ls":
		ls.ListLogicalSwitchDetail(ovnPod, resource)
	case "chassis":
		ch.ListChassisDetail(ovnPod, resource)
	case "ch":
		ch.ListChassisDetail(ovnPod, resource)
	default:
		tools.PrintHelp()
		os.Exit(1)
	}
}

func runShowCommand(resource []string, ovnPod string) {
	lr := LogicalRouter.New()
	ch := Chassis.New()
	tools := Tools.New()
	switch resource[0] {
	case "logicalrouter":
		lr.ShowLogicalRoutersDetail(ovnPod, resource)
	case "lr":
		lr.ShowLogicalRoutersDetail(ovnPod, resource)
	case "chassis":
		ch.ShowChassisDetail(ovnPod, resource)
	case "ch":
		ch.ShowChassisDetail(ovnPod, resource)
	default:
		tools.PrintHelp()
		os.Exit(1)
	}
}
