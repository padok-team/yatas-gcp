package instance

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func InstanceNoPublicIPAttached(resource commons.Resource) bool {
	instance, ok := resource.(*VMInstance)
	if !ok {
		return false
	}

	if instance.Instance.NetworkInterfaces != nil {
		for _, networkInterface := range instance.Instance.NetworkInterfaces {
			if networkInterface.AccessConfigs != nil {
				for _, accessConfig := range networkInterface.AccessConfigs {
					// If there is an external IPv4 or IPv6 address, return false
					if accessConfig.NatIP != nil && *accessConfig.NatIP != "" {
						return false
					}
					if accessConfig.ExternalIpv6 != nil {
						return false
					}
				}
			}
		}
	}
	return true
}

func DiskIsCustomerEncrypted(resource commons.Resource) bool {
	disk, ok := resource.(*VMDisk)
	if !ok {
		return false
	}
	return disk.Disk.DiskEncryptionKey != nil
}
