# Nutanix Flow Networking CLI tools
## Description
This tool is developed to help the user to execute the commands on Flow Networking Atlas Network Controllers (ANC)

## Prerequirements
The *flow* command requires to be run on the *master* node of ANC and requires *kubectl* already installed. To run the *flow* command, the *flow* binary file needs to be copied the `/usr/local/bin` path of the *master nodes*.

## Help
    Usage:
        flow <command> <resource>
    Commands:
        list|ls                           displays display a list of resource
        show|sh                           displays detailed information of a resource
        version                           displays version number
        help                              displays the help
    Resources:
        logicalswitch|ls                  logical switch
        logicalrouter|lr                  logical router
        chassis|ch                        chassis


## Example:
        flow ls lr                        displays the list of logical routers and their ports
        flow ls ls                        displays the list of logical switches and their ports
        flow ls ch                        displays the list of chassis
        flow ls lr nat                    displays the list of NAT on each logicalrouter
        flow ls lr routes                 displays the list of routes on each logicalrouter
        flow sh lr routes <router-id>     displays the list of routes on an specific logicalrouter
        flow sh ch <chassis-id>           displays detailed information about an specific chassis

        [nutanix@anc-1e21947f1cea-e88643-default-0 ~]$ flow ls lr
        +---------------------------------------------+----------+-------------------+-----------------+------------------+
        |                 ROUTER NAME                 | PORT NO  |    MAC ADDRESS    |   IP ADDRESS    | REDIRECT CHASSIS |
        +---------------------------------------------+----------+-------------------+-----------------+------------------+
        | router_5d0b5798-4962-4e12-b902-dbb7ebaaa322 |    1     | e0:19:95:c5:66:41 | 192.168.12.1/24 |                  |
        +                                             +----------+-------------------+-----------------+------------------+
        |                                             |    2     | e0:19:95:8e:d9:c3 | 192.168.11.1/24 |                  |
        +---------------------------------------------+----------+-------------------+-----------------+------------------+
        
        [nutanix@anc-1e21947f1cea-e88643-default-0 ~]$ flow ls net
        +----------------------------------------------+-----------+-------------------+-------------------+-----------+
        |             LOGICAL SWITCH NAME              | TUNNEL ID |      PORT IP      |     PORT MAC      | PORT TYPE |
        +----------------------------------------------+-----------+-------------------+-------------------+-----------+        
        | network_de1213e1-946a-48d4-b1d0-44205b219764 |   10001   | 192.168.11.129/24 | 50:6b:8d:30:e9:51 |           |
        +                                              +           +-------------------+-------------------+-----------+
        |                                              |           |  192.168.11.1/24  |                   |  router   |
        +                                              +           +-------------------+-------------------+-----------+
        |                                              |           | 192.168.11.136/24 | 50:6b:8d:40:35:44 |           |
        +----------------------------------------------+-----------+-------------------+-------------------+-----------+
        | network_8456e6c7-d995-4c5f-9ea1-953e954c3862 |   10002   |  192.168.12.1/24  |                   |  router   |
        +----------------------------------------------+-----------+-------------------+-------------------+-----------+
        
        [nutanix@anc-1e21947f1cea-e88643-default-0 ~]$ flow ls ports
        +------------------------------------------------------+--------------------------------------+-----------------+-------------------+-------+-----------------+
        |                      PORT NAME                       |               CHASSIS                |       IP        |        MAC        | TYPE  | GATEWAY CHASSIS |
        +------------------------------------------------------+--------------------------------------+-----------------+-------------------+-------+-----------------+
        |      port_3d1a688e-78e7-4af3-8d3a-813671129bcc       | ea8a21c4-0b07-4037-8580-7904e1da920d | 192.168.11.129  | 50:6b:8d:30:e9:51 |       |                 |
        +------------------------------------------------------+--------------------------------------+-----------------+-------------------+-------+-----------------+
        | lrp-router-port_de1213e1-946a-48d4-b1d0-44205b219764 |                                      | 192.168.11.1/24 | e0:19:95:8e:d9:c3 | patch |                 |
        +------------------------------------------------------+--------------------------------------+-----------------+-------------------+-------+-----------------+
        | lrp-router-port_8456e6c7-d995-4c5f-9ea1-953e954c3862 |                                      | 192.168.12.1/24 | e0:19:95:c5:66:41 | patch |                 |
        +------------------------------------------------------+--------------------------------------+-----------------+-------------------+-------+-----------------+
        |      port_0d8a2764-817f-4d6f-809e-b53eaeb6d5af       | 68b72a31-a951-475d-ac16-e7040e3ca118 | 192.168.11.136  | 50:6b:8d:40:35:44 |       |                 |
        +------------------------------------------------------+--------------------------------------+-----------------+-------------------+-------+-----------------+
        
        [nutanix@anc-1e21947f1cea-e88643-default-0 ~]$ flow ls chassis
        +--------------------------------------+-----------------+--------+--------------+
        |             CHASSIS NAME             |    HOSTNAME     | ENCAP  |      IP      |
        +--------------------------------------+-----------------+--------+--------------+
        | ea8a21c4-0b07-4037-8580-7904e1da920d | aa-nested-ahv-2 | geneve | 10.48.29.201 |
        +--------------------------------------+-----------------+--------+--------------+
        | 0c014e57-d249-4a61-8cb6-938dff5aac4d | aa-nested-ahv-1 | geneve | 10.48.29.200 |
        +--------------------------------------+-----------------+--------+--------------+
        | 68b72a31-a951-475d-ac16-e7040e3ca118 | aa-nested-ahv-3 | geneve | 10.48.29.202 |
        +--------------------------------------+-----------------+--------+--------------+
