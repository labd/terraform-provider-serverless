package serverless

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func setStackInfo(stacksOutput *cloudformation.DescribeStacksOutput, d *schema.ResourceData) error {
	if len(stacksOutput.Stacks) == 0 {
		return nil
	}

	stack := stacksOutput.Stacks[0]

	functionArns, err := getFunctionArns(stack)
	if err != nil {
		return err
	}
	if err := d.Set("function_arns", functionArns); err != nil {
		return err
	}

	for _, output := range stack.Outputs {
		if *output.OutputKey == "HttpApiUrl" {
			if err := d.Set("http_api_url", *output.OutputValue); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func getFunctionArns(stack *cloudformation.Stack) (map[string]string, error) {
	result := map[string]string{}

	arnPattern := regexp.MustCompile(`^(.*)LambdaFunctionQualifiedArn$`)

	for _, output := range stack.Outputs {
		match := arnPattern.FindStringSubmatch(*output.OutputKey)
		if match == nil {
			continue
		}

		functionArnName := deNormalizeFunctionName(match[1])
		result[functionArnName] = *output.OutputValue
	}

	return result, nil
}

// We're looking for a normalized lambda function name here.
// The original function in serverless looks like this:
//    normalizeName(name) {
//       return `${_.upperFirst(name)}`;
//    },
// 	  getNormalizedFunctionName(functionName) {
// 	     return this.normalizeName(functionName.replace(/-/g, 'Dash').replace(/_/g, 'Underscore'));
// 	  },
//
func deNormalizeFunctionName(value string) string {
	value = strings.Replace(value, "Dash", "-", -1)
	value = strings.Replace(value, "Underscore", "_", -1)
	return lowerFirst(value)
}

func lowerFirst(value string) string {
	for i, v := range value {
		return string(unicode.ToLower(v)) + value[i+1:]
	}
	return ""
}
