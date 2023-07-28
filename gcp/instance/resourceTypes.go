package instance

import (
	"fmt"
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
)

// This type implements commons.Resource
type VMInstance struct {
	Instance computepb.Instance
}

func (i *VMInstance) GetID() string {
	zoneURLSplit := strings.Split(i.Instance.GetZone(), "/")
	zoneName := zoneURLSplit[len(zoneURLSplit)-1]
	return fmt.Sprintf("%s/%s (%d)", zoneName, i.Instance.GetName(), i.Instance.GetId())
}

// This type implements commons.Resource
type VMDisk struct {
	Disk computepb.Disk
}

func (d *VMDisk) GetID() string {
	zoneURLSplit := strings.Split(d.Disk.GetZone(), "/")
	zoneName := zoneURLSplit[len(zoneURLSplit)-1]
	return fmt.Sprintf("%s/%s (%d)", zoneName, d.Disk.GetName(), d.Disk.GetId())
}

// This type implements commons.Resource
type InstanceGroup struct {
	InstanceGroup computepb.InstanceGroupManager
}

func (i *InstanceGroup) GetID() string {
	zoneURLSplit := strings.Split(i.InstanceGroup.GetZone(), "/")
	zoneName := zoneURLSplit[len(zoneURLSplit)-1]
	if zoneName == "" {
		regionURLSplit := strings.Split(i.InstanceGroup.GetRegion(), "/")
		zoneName = regionURLSplit[len(regionURLSplit)-1]
	}
	return fmt.Sprintf("%s/%s (%d)", zoneName, i.InstanceGroup.GetName(), i.InstanceGroup.GetId())
}
