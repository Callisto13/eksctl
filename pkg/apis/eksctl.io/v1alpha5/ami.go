package v1alpha5

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Values for `NodeAMIFamily`
// All valid values should go in this block
const (
	// DefaultNodeImageFamily (default)
	DefaultNodeImageFamily      = NodeImageFamilyAmazonLinux2
	NodeImageFamilyAmazonLinux2 = "AmazonLinux2"
	NodeImageFamilyUbuntu2004   = "Ubuntu2004"
	NodeImageFamilyUbuntu1804   = "Ubuntu1804"
	NodeImageFamilyBottlerocket = "Bottlerocket"

	NodeImageFamilyWindowsServer2019CoreContainer = "WindowsServer2019CoreContainer"
	NodeImageFamilyWindowsServer2019FullContainer = "WindowsServer2019FullContainer"
	NodeImageFamilyWindowsServer1909CoreContainer = "WindowsServer1909CoreContainer"
	NodeImageFamilyWindowsServer2004CoreContainer = "WindowsServer2004CoreContainer"
)

const (
	// ownerIDUbuntuFamily is the owner ID used for Ubuntu AMIs
	ownerIDUbuntuFamily = "099720109477"

	// ownerIDWindowsFamily is the owner ID used for Ubuntu AMIs
	ownerIDWindowsFamily = "801119661308"

	// eksResourceAccountStandard defines the AWS EKS account ID that provides node resources in default regions
	// for standard AWS partition
	eksResourceAccountStandard = "602401143452"

	// eksResourceAccountAPEast1 defines the AWS EKS account ID that provides node resources in ap-east-1 region
	eksResourceAccountAPEast1 = "800184023465"

	// eksResourceAccountMESouth1 defines the AWS EKS account ID that provides node resources in me-south-1 region
	eksResourceAccountMESouth1 = "558608220178"

	// eksResourceAccountCNNorthWest1 defines the AWS EKS account ID that provides node resources in cn-northwest-1 region
	eksResourceAccountCNNorthWest1 = "961992271922"

	// eksResourceAccountCNNorth1 defines the AWS EKS account ID that provides node resources in cn-north-1
	eksResourceAccountCNNorth1 = "918309763551"

	// eksResourceAccountAFSouth1 defines the AWS EKS account ID that provides node resources in af-south-1
	eksResourceAccountAFSouth1 = "877085696533"

	// eksResourceAccountEUSouth1 defines the AWS EKS account ID that provides node resources in eu-south-1
	eksResourceAccountEUSouth1 = "590381155156"

	// eksResourceAccountUSGovWest1 defines the AWS EKS account ID that provides node resources in us-gov-west-1
	eksResourceAccountUSGovWest1 = "013241004608"

	// eksResourceAccountUSGovEast1 defines the AWS EKS account ID that provides node resources in us-gov-east-1
	eksResourceAccountUSGovEast1 = "151742754352"
)

var eksResourceAccount ResourceAccount

type ResourceAccount struct {
	ID string
}

// SetEKSResourceAccount looks up the Region account ID based on the region
func SetEKSResourceAccount(region *string) error {
	var err error
	eksResourceAccount.ID, err = resourceAccountID(*region)
	if err != nil {
		return err
	}
	return nil
}

// EKSResourceAccountID provides worker node resources(ami/ecr image) in different aws account
// for different aws partitions & opt-in regions.
func EKSResourceAccountID() string {
	return eksResourceAccount.ID
}

// OwnerAccountID returns the AWS account ID that owns worker AMI.
func OwnerAccountID(imageFamily, region string) (string, error) {
	switch imageFamily {
	case NodeImageFamilyUbuntu2004, NodeImageFamilyUbuntu1804:
		return ownerIDUbuntuFamily, nil
	case NodeImageFamilyAmazonLinux2:
		return eksResourceAccount.ID, nil
	default:
		if IsWindowsImage(imageFamily) {
			return ownerIDWindowsFamily, nil
		}
		return "", fmt.Errorf("unable to determine the account owner for image family %s", imageFamily)
	}
}

// IsWindowsImage reports whether the AMI family is for Windows
func IsWindowsImage(imageFamily string) bool {
	switch imageFamily {
	case NodeImageFamilyWindowsServer2019CoreContainer,
		NodeImageFamilyWindowsServer2019FullContainer,
		NodeImageFamilyWindowsServer1909CoreContainer,
		NodeImageFamilyWindowsServer2004CoreContainer:
		return true

	default:
		return false
	}
}

func resourceAccountID(region string) (string, error) {
	switch region {
	case RegionUSIsoEast1:
		return isoResourceAccountID(region, "ISO_EAST1_ACCOUNT_ID")
	case RegionUSIsobEast1:
		return isoResourceAccountID(region, "ISO_B_EAST1_ACCOUNT_ID")
	default:
		return publicResourceAccountID(region), nil
	}
}

func isoResourceAccountID(region, varName string) (string, error) {
	id, ok := os.LookupEnv(varName)
	if !ok {
		return "", errors.Errorf("%s not set, required for use of region: %s", varName, region)
	}

	return id, nil
}

func publicResourceAccountID(region string) string {
	switch region {
	case RegionAPEast1:
		return eksResourceAccountAPEast1
	case RegionMESouth1:
		return eksResourceAccountMESouth1
	case RegionCNNorthwest1:
		return eksResourceAccountCNNorthWest1
	case RegionCNNorth1:
		return eksResourceAccountCNNorth1
	case RegionUSGovWest1:
		return eksResourceAccountUSGovWest1
	case RegionUSGovEast1:
		return eksResourceAccountUSGovEast1
	case RegionAFSouth1:
		return eksResourceAccountAFSouth1
	case RegionEUSouth1:
		return eksResourceAccountEUSouth1
	default:
		return eksResourceAccountStandard
	}
}
