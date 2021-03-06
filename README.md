# Terraform Provider for Serverless

Deploy and manage your serverless stack with Terraform.

**This is a fork** from the great work delivered by [FormidableLabs](https://github.com/FormidableLabs/terraform-provider-serverless).  
It adds:

- Extra functionality (see [Changelog](./CHANGELOG.md))
- Publish to Terraform registry for easier integration into our projects

## Why?

It's common to augment the simplicity and developer experience of Serverless with the power and flexibility of Terraform, either to fill in gaps in Serverless or to provision its supporting resources. However, Serverless and Terraform don't interoperate out-of-the-box. Even if you get the two to communicate, neither of them are aware of their dependencies on each other. This can lead to chicken-and-egg problems where one Terraform resource must exist before Serverless, but the rest of your resources depend on a Serverless deploy.

The Serverless provider integrates Serverless deploys as a Terraform resource so that Terraform can resolve resource dependencies correctly. It's a great tool for shimming Serverless into a Terraform-heavy project or easing the migration cost away from Serverless to pure Terraform.

## Installation

The Serverless provider is not an official provider, so `terraform init` doesn't automatically download and install it. See HashiCorp's instructions for [installing third-party providers.](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)

## Packaging requirements

This provider avoids deploying Serverless when code or configuration hasn't changed. This saves deployment time, supports deployment of bit-identical package artifacts across different environments, and prevents application code changes from reaching users when you only need to update configuration for other Terraform resources.

To avoid spurious deploys, Serverless requires you to package out-of-band from `sls deploy`. To update the package and trigger a deploy on the next `terraform apply`, run `sls package -p .terraform-serverless`. You should no longer use the `sls deploy` command.

## Resources

```hcl
data "aws_caller_identity" "current" {}

resource "serverless_deployment" "example" {
  # Used to support assume role policies.
  # Whenever this resource is used within a certain Terraform AWS context,
  # we must make sure the correct AWS credentails are used
  aws_config {
    account_id  = data.aws_caller_identity.current.account_id
    caller_arn  = data.aws_caller_identity.current.arn
    caller_user = data.aws_caller_identity.current.user_id
  }

  # The Serverless stage. Usually corresponds to the stage/environment of the Terraform module.
  stage               = "sandbox"

  # The directory where your `serverless.yml` config lives. Must be an absolute path.
  config_dir         = abspath("example")

  # The directory where your Serverless package lives. Defaults to `.terraform-serverless`.
  # **NOTE:** must be relative to `config_dir`!
  package_dir         = ".terraform-serverless"

  # The directory to look for the `serverless` binary. Otherwise uses the binary in your current $PATH
  # serverless_bin_dir = abspath("example/node_modules/.bin")

  # All environment variables are exposed to serverless and can be used in `serverless.yml`
  env = {
    "AWS_REGION": data.aws_region.current.name,
    "FOO": "BAR",
  }
}

output "serverless_url" {
  value = serverless_deployment.example.http_api_url
}

output "serverless_arns" {
  value = {
    "hello": serverless_deployment.example.function_arns.hello,
    "bye": serverless_deployment.example.function_arns.bye,
  }
}
```

## Contributing

See our contribution guidelines [here!](CONTRIBUTING.md)

## Maintenance Status

**Experimental:** This project is quite new. We're not sure what our ongoing maintenance plan for this project will be. Bug reports, feature requests and pull requests are welcome. If you like this project, let us know!
