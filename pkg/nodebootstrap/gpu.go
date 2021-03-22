package nodebootstrap

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// CONTEXT FOR THIS FILE: https://github.com/weaveworks/eksctl/issues/3446

const (
	// This was the date that GPUs were first released with the changes.
	// The changes were then rolled back until XXXX, but we have to process this
	// one too
	rogueDate = "2021-03-03"

	// This is the date from which GPU images have the current config changes
	newImageFormatDate = "2021-03-10T21:48:02.000Z"

	gpuImageNameGlob = "amazon-eks-gpu-node-1*"
)

func IsNewGPUAMI(ec2API ec2iface.EC2API, ami string) (bool, error) {
	images, err := describeGPUImages(ec2API)
	if err != nil {
		return false, err
	}
	for _, image := range images {
		if ami == aws.StringValue(image.ImageId) {
			isNew, err := isNewImage(aws.StringValue(image.CreationDate))
			if err != nil {
				return false, err
			}
			if isNew {
				return true, nil
			}
		}
	}
	return false, nil
}

func describeGPUImages(ec2API ec2iface.EC2API) ([]*ec2.Image, error) {
	input := ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{Name: aws.String("name"), Values: []*string{aws.String(gpuImageNameGlob)}},
		},
	}

	output, err := ec2API.DescribeImages(&input)
	if err != nil {
		return nil, err
	}

	return output.Images, nil
}

func isNewImage(imageCreatedAt string) (bool, error) {
	if strings.Contains(imageCreatedAt, rogueDate) {
		return true, nil
	}

	t1, err := time.Parse(time.RFC3339, newImageFormatDate)
	if err != nil {
		return false, nil
	}

	t2, err := time.Parse(time.RFC3339, imageCreatedAt)
	if err != nil {
		return false, nil
	}

	return t2.After(t1), nil
}

// TODO
// - manual testing, below
// - unit-testing for this file
// - refactoring make data
// - get new format date

// testing combinations:
// - gpu old format
// - gpu new format
// - non-gpu
// - instances distribution w/ gpu types
// - instances distribution w/o gpu types
// - instances distribution w/ array + instance type set
