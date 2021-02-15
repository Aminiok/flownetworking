package Tools

import (
	"flownet/Tools/ipsubnet"
	"flownet/Tools/tablewriter"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type params struct {
	ovnPod string
}

func New() params {
	param := params{}
	return param
}

func pipeCommands(commands ...*exec.Cmd) (out []byte, err error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		log.Fatal(err)
	}
	return final, err
}

func (param *params) GetVersion() {
	fmt.Println("1.0.0")
}

func (param *params) GetOvnPod() string {
	kubeCommand := exec.Command("/usr/bin/kubectl", "get", "pods")
	grepCommand := exec.Command("grep", "anc-ovn")
	awkCommand := exec.Command("awk", "{print $1}")
	out, err := pipeCommands(kubeCommand, grepCommand, awkCommand)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func (param *params) PrintHelp() {
	helpText := "\n" +
		"\n" +
		"Usage:\n" +
		"    flow <command> <resource>\n" +
		"Commands:\n" +
		"    list|ls               displays display a list of resource\n" +
		"    show|sh               displays detailed information of a resource\n" +
		"    version               displays version number\n" +
		"    help                  displays the help\n" +
		"Resources:\n" +
		"    logicalrouter|lr      logical router\n" +
		"    chassis|ch            chassis\n\n" +
		"Example:\n" +
		"    flow ls lr            displays the list of logical routers\n" +
		"    flow ls ch            displays the list of chassis\n" +
		"    flow ls lr nat        displays the list of NAT or each logicalrouter\n" +
		"    flow ls lr routes     displays the list of routes on each logicalrouter\n" +
		"    flow sh lr routes <router-id>\n" +
		"    flow sh ch <chassis-id>\n"
	fmt.Println(helpText)
	os.Exit(1)
}

func (param *params) ShowInTable(data [][]string, header []string, mergeColumnIndex []int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoMergeCellsByColumnIndex(mergeColumnIndex)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.AppendBulk(data)
	table.Render()
}

func (param *params) RefactorStringList(stringList []string) []string {
	newList := []string{}
	reg, err := regexp.Compile("[^a-zA-Z0-9/.:_-]+")
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range stringList {
		processedString := reg.ReplaceAllString(item, "")
		newList = append(newList, processedString)
	}
	return newList
}

func (param *params) RefactorString(origString string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9/.:_-]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(origString, "")
	return processedString
}

func (param *params) GetNetworkFromIP(ipAddress string) string {
	ipAddressArray := strings.Split(ipAddress, "/")
	if len(ipAddressArray) == 2 {
		netMask, err := strconv.Atoi(ipAddressArray[1])
		if err != nil {
			log.Fatal(err)
		}
		sub := ipsubnet.SubnetCalculator(ipAddressArray[0], netMask)
		network := sub.GetNetworkPortion()
		return network + "/" + strconv.Itoa(netMask)
	}
	return ""
}

func (param *params) Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
