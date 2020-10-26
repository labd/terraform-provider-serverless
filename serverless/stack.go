package serverless

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func setStackInfo(stacksOutput *cloudformation.DescribeStacksOutput, d *schema.ResourceData) error {
	if len(stacksOutput.Stacks) == 0 {
		return nil
	}

	stack := stacksOutput.Stacks[0]

	for _, output := range stack.Outputs {
		if *output.OutputKey == "HttpApiUrl" {
			d.Set("http_api_url", *output.OutputValue)
		}
	}

	return nil
}
