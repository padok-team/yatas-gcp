package instance

import (
	"testing"

	"cloud.google.com/go/compute/apiv1/computepb"
)

func stringPtr(s string) *string {
	return &s
}

type fakeResource struct{}

func (f fakeResource) GetID() string {
	return "fake-resource"
}

func TestInstanceNoPublicIPAttached(t *testing.T) {
	testCases := []struct {
		desc     string
		instance *computepb.Instance
		want     bool
	}{
		{
			desc: "no network interfaces",
			instance: &computepb.Instance{
				Name: stringPtr("test-instance-1"),
				Zone: stringPtr("projects/test-project/zones/us-central1-a"),
				Disks: []*computepb.AttachedDisk{
					{
						Source: stringPtr("projects/test-project/zones/us-central1-a/disks/test-disk"),
					},
				},
			},
			want: true,
		},
		{
			desc: "no access configs",
			instance: &computepb.Instance{
				Name: stringPtr("test-instance-2"),
				Zone: stringPtr("projects/test-project/zones/us-central1-a"),
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Network: stringPtr("projects/test-project/global/networks/default"),
					},
				},
			},
			want: true,
		},
		{
			desc: "no external IPv4 or IPv6 addresses",
			instance: &computepb.Instance{
				Name: stringPtr("test-instance-3"),
				Zone: stringPtr("projects/test-project/zones/us-central1-a"),
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Network: stringPtr("projects/test-project/global/networks/default"),
						AccessConfigs: []*computepb.AccessConfig{
							{},
						},
					},
				},
			},
			want: true,
		},
		{
			desc: "external IPv4 address",
			instance: &computepb.Instance{
				Name: stringPtr("test-instance-4"),
				Zone: stringPtr("projects/test-project/zones/us-central1-a"),
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Network: stringPtr("projects/test-project/global/networks/default"),
						AccessConfigs: []*computepb.AccessConfig{
							{
								NatIP: stringPtr("1.2.3.4"),
							},
						},
					},
				},
			},
			want: false,
		},
		{
			desc: "external IPv6 address",
			instance: &computepb.Instance{
				Name: stringPtr("test-instance-5"),
				Zone: stringPtr("projects/test-project/zones/us-central1-a"),
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Network: stringPtr("projects/test-project/global/networks/default"),
						AccessConfigs: []*computepb.AccessConfig{
							{
								ExternalIpv6: stringPtr("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
							},
						},
					},
				},
			},
			want: false,
		},
		{
			desc: "external IPv4 and IPv6 addresses",
			instance: &computepb.Instance{
				Name: stringPtr("test-instance-6"),
				Zone: stringPtr("projects/test-project/zones/us-central1-a"),
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Network: stringPtr("projects/test-project/global/networks/default"),
						AccessConfigs: []*computepb.AccessConfig{
							{
								NatIP:        stringPtr("1.2.3.4"),
								ExternalIpv6: stringPtr("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
							},
						},
					},
				},
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			instance := &VMInstance{Instance: *tc.instance}
			got := InstanceNoPublicIPAttached(instance)
			if got != tc.want {
				t.Errorf("InstanceNoPublicIPAttached(%v) = %v, want %v", instance, got, tc.want)
			}
		})
	}
}

func TestDiskIsEncrypted(t *testing.T) {
	// Create encrypted and unencrypted disks
	key := computepb.CustomerEncryptionKey{}
	diskWithKey := computepb.Disk{DiskEncryptionKey: &key}
	disk := computepb.Disk{}

	encryptedDisk := VMDisk{Disk: diskWithKey}
	unencryptedDisk := VMDisk{Disk: disk}

	// Ensure that the function correctly detects encrypted and unencrypted disks
	if !DiskIsCustomerEncrypted(&encryptedDisk) {
		t.Error("Expected encrypted disk to be detected as encrypted")
	}
	if DiskIsCustomerEncrypted(&unencryptedDisk) {
		t.Error("Expected unencrypted disk to be detected as unencrypted")
	}

	// Ensure that the function returns false for non-disk resources
	nonDiskResource := &fakeResource{}
	if DiskIsCustomerEncrypted(nonDiskResource) {
		t.Error("Expected non-disk resource to be detected as non-encrypted")
	}
}
