package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nexodus-io/nexodus/internal/api/public"
)

func listOrgDevices(c *public.APIClient, organizationID uuid.UUID, fullDisplay bool, encodeOut string) error {
	devices, _, err := c.DevicesApi.ListDevicesInOrganization(context.Background(), organizationID.String()).Execute()
	if err != nil {
		log.Fatal(err)
	}
	if encodeOut == encodeColumn || encodeOut == encodeNoHeader {
		w := newTabWriter()

		var fs string
		if fullDisplay {
			fs = "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n"
		} else {
			fs = "%s\t%s\t%s\t%s\t%s\n"
		}
		if encodeOut != encodeNoHeader && !fullDisplay {
			fmt.Fprintf(w, fs, "DEVICE ID", "HOSTNAME", "TUNNEL IPS", "ORGANIZATION ID", "RELAY")
		}
		if encodeOut != encodeNoHeader && fullDisplay {
			fmt.Fprintf(w, fs, "DEVICE ID", "HOSTNAME",
				"ENDPOINT IP", "PUBLIC KEY", "ORGANIZATION ID",
				"LOCAL IP", "ALLOWED IPS", "TUNNEL IPV4", "TUNNEL IPV6",
				"CHILD PREFIX", "ORG PREFIX IPV4", "ORG PREFIX IPV6",
				"REFLEXIVE IPv4", "ENDPOINT LOCAL IPv4", "OS", "SECURITY GROUP ID", "RELAY")
		}
		for _, dev := range devices {
			localIp := ""
			var reflexiveIp4 []string
			for _, endpoint := range dev.Endpoints {
				if endpoint.Source == "local" {
					localIp = endpoint.Address
				} else {
					reflexiveIp4 = append(reflexiveIp4, endpoint.Address)
				}
			}
			if !fullDisplay {
				fmt.Fprintf(w, fs, dev.Id, dev.Hostname, fmt.Sprintf("[%s %s]", dev.TunnelIp, dev.TunnelIpV6), dev.OrganizationId, fmt.Sprintf("%t", dev.Relay))
			} else {
				fmt.Fprintf(w, fs, dev.Id, dev.Hostname, localIp, dev.PublicKey, dev.OrganizationId,
					localIp, dev.AllowedIps, dev.TunnelIp, dev.TunnelIpV6, dev.ChildPrefix, dev.OrganizationPrefix,
					dev.OrganizationPrefixV6, reflexiveIp4, dev.EndpointLocalAddressIp4, dev.Os, dev.SecurityGroupId, fmt.Sprintf("%t", dev.Relay))
			}
		}
		w.Flush()

		return nil
	}

	err = FormatOutput(encodeOut, devices)
	if err != nil {
		log.Fatalf("failed to print output: %v", err)
	}

	return nil
}

func listAllDevices(c *public.APIClient, fullDisplay bool, encodeOut string) error {
	devices, _, err := c.DevicesApi.ListDevices(context.Background()).Execute()
	if err != nil {
		log.Fatal(err)
	}
	if encodeOut == encodeColumn || encodeOut == encodeNoHeader {
		w := newTabWriter()
		var fs string
		if fullDisplay {
			fs = "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n"
		} else {
			fs = "%s\t%s\t%s\t%s\t%s\n"
		}
		if encodeOut != encodeNoHeader && !fullDisplay {
			fmt.Fprintf(w, fs, "DEVICE ID", "HOSTNAME", "TUNNEL IPS", "ORGANIZATION ID", "RELAY")
		}
		if encodeOut != encodeNoHeader && fullDisplay {
			fmt.Fprintf(w, fs, "DEVICE ID", "HOSTNAME",
				"ENDPOINT IP", "PUBLIC KEY", "ORGANIZATION ID",
				"LOCAL IP", "ALLOWED IPS", "TUNNEL IPV4", "TUNNEL IPV6",
				"CHILD PREFIX", "ORG PREFIX IPV4", "ORG PREFIX IPV6",
				"REFLEXIVE IPv4", "ENDPOINT LOCAL IPv4", "OS", "SECURITY GROUP ID", "RELAY")
		}
		for _, dev := range devices {
			localIp := ""
			var reflexiveIp4 []string
			for _, endpoint := range dev.Endpoints {
				if endpoint.Source == "local" {
					localIp = endpoint.Address
				} else {
					reflexiveIp4 = append(reflexiveIp4, endpoint.Address)
				}
			}
			if !fullDisplay {
				fmt.Fprintf(w, fs, dev.Id, dev.Hostname, fmt.Sprintf("[%s %s]", dev.TunnelIp, dev.TunnelIpV6), dev.OrganizationId, fmt.Sprintf("%t", dev.Relay))
			} else {
				fmt.Fprintf(w, fs, dev.Id, dev.Hostname, localIp, dev.PublicKey, dev.OrganizationId,
					localIp, dev.AllowedIps, dev.TunnelIp, dev.TunnelIpV6, dev.ChildPrefix, dev.OrganizationPrefix,
					dev.OrganizationPrefixV6, reflexiveIp4, dev.EndpointLocalAddressIp4, dev.Os, dev.SecurityGroupId, fmt.Sprintf("%t", dev.Relay))
			}
		}
		w.Flush()

		return nil
	}

	err = FormatOutput(encodeOut, devices)
	if err != nil {
		log.Fatalf("failed to print output: %v", err)
	}

	return nil
}

func deleteDevice(c *public.APIClient, encodeOut, devID string) error {
	devUUID, err := uuid.Parse(devID)
	if err != nil {
		log.Fatalf("failed to parse a valid UUID from %s %v", devUUID, err)
	}

	res, _, err := c.DevicesApi.DeleteDevice(context.Background(), devUUID.String()).Execute()
	if err != nil {
		log.Fatalf("device delete failed: %v\n", err)
	}

	if encodeOut == encodeColumn || encodeOut == encodeNoHeader {
		fmt.Printf("successfully deleted device %s\n", res.Id)
		return nil
	}

	err = FormatOutput(encodeOut, res)
	if err != nil {
		log.Fatalf("failed to print output: %v", err)
	}

	return nil
}
