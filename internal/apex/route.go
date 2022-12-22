package apex

import (
	"net"
)

// handlePeerRoute when a new configuration is deployed, delete/add the peer allowedIPs
func (ax *Apex) handlePeerRoute(wgPeerConfig wgPeerConfig) {
	switch ax.os {
	case Darwin.String():
		// Darwin maps to a utunX address which needs to be discovered (currently hardcoded to utun8)
		devName, err := getInterfaceByIP(net.ParseIP(ax.wgLocalAddress))
		if err != nil {
			ax.logger.Debugf("failed to find the darwin interface with the address [ %s ] %v", ax.wgLocalAddress, err)
		}
		// If child prefix split the two prefixes (host /32 and child prefix
		for _, allowedIP := range wgPeerConfig.AllowedIPs {
			_, err := RunCommand("route", "-q", "-n", "delete", "-inet", allowedIP, "-interface", devName)
			if err != nil {
				ax.logger.Debugf("no route deleted: %v", err)
			}
			if err := AddRoute(allowedIP, devName); err != nil {
				ax.logger.Debugf("%v", err)
			}
		}

	case Linux.String():
		for _, allowedIP := range wgPeerConfig.AllowedIPs {
			routeExists, err := RouteExists(allowedIP)
			if err != nil {
				ax.logger.Info(err)
			}
			if !routeExists {
				if err := AddRoute(allowedIP, wgIface); err != nil {
					ax.logger.Errorf("route add failed: %v", err)
				}
			}
		}

	case Windows.String():
		for _, allowedIP := range wgPeerConfig.AllowedIPs {
			if err := AddRoute(allowedIP, wgIface); err != nil {
				ax.logger.Debugf("route add failed: %v", err)
			}
		}
	}
}

func (ax *Apex) addChildPrefixRoute(childPrefix string) {
	var dev string

	if ax.os == Darwin.String() {
		dev = darwinIface
	} else {
		dev = wgIface
	}

	routeExists, err := RouteExists(childPrefix)
	if err != nil {
		ax.logger.Warn(err)
	}

	if routeExists {
		ax.logger.Debugf("unable to add the child-prefix route [ %s ] as it already exists on this linux host", childPrefix)
		return
	}

	if err := AddRoute(childPrefix, dev); err != nil {
		ax.logger.Infof("error adding the child prefix route: %v", err)
	}
}
