# Nutanix Flow Networking CLI tools
## Description
This tool is developed to help the user to execute the commands on Flow Networking Atlase Network Controllers (ANC)

## Prerequirements
The *flow* command requires to be run on the *master* node of ANC and requires *kubectl* already installed. To run the *flow* command, the *flow* binary file needs to be copied the `</usr/local/bin>` path of the *master nodes*.

## Help
    Usage:
        flow <command> <resource>
    Commands:
        list|ls                           displays display a list of resource
        show|sh                           displays detailed information of a resource
        version                           displays version number
        help                              displays the help
    Resources:
        logicalrouter|lr                  logical router
        chassis|ch                        chassis


## Example:
        flow ls lr                        displays the list of logical routers
        flow ls ch                        displays the list of chassis
        flow ls lr nat                    displays the list of NAT on each logicalrouter
        flow ls lr routes                 displays the list of routes on each logicalrouter
        flow sh lr routes <router-id>     displays the list of routes on an specific logicalrouter
        flow sh ch <chassis-id>           displays detailed information about an specific chassis
