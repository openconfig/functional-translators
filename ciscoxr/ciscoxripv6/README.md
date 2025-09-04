The paths parsed here are sourced from an unkeyed list of addresses.

It is easier to just raw manipulate the gNMI notifications, rather than using
ygot since it can't unmarshal the notifications.

*   Output Path:
    /openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/addresses/address/state/ip
    *   Input Paths:
        *   /Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address
*   Output Path:
    /openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/addresses/address/state/prefix-length
    *   Input Paths:
        *   /Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address
        *   /Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/prefix-length
