# adaptive-metrics-autoapply

A template repository for enabling auto-apply of Adaptive Metrics recommendations.

## Getting started

To start, create a new repository by navigating to "Use this template" -> "Create a new repository" up top.

Add the following secrets to the new repository ("Settings" -> "Secrets and variables" -> "Actions" -> "New repository secret"):

- `grafana_am_api_url`: This is your Grafana Cloud prometheus URL. To find this URL, go to your `grafana.com` account and check the **Details** page of your hosted Prometheus endpoint.
- `grafana_am_api_key`: This looks like `<your-numeric-instance-id>:<your-cloud-access-policy-token>`.
    - `<your-numeric-instance-id>` is the numeric instance ID where you want to enable auto-apply of your Adaptive Metrics. To find this value, go to your `grafana.com` account and check the **Details** page of your hosted Prometheus endpoint for **Username/Instance ID**.
    - `<your-cloud-access-policy-token>` is a token from a [Grafana Cloud Access Policy](https://grafana.com/docs/grafana-cloud/account-management/authentication-and-permissions/access-policies/). Make sure the access policy has `metrics:read` and `metrics:write` scopes for the appropriate stack ID.

## What to expect

By default, auto-apply is scheduled to run at 04:00 UTC every day. This can be configured by editing the schedule time in `.github/workflows/main.yml`.

At the scheduled time, the GitHub Action will pull the latest recommendations and apply them as aggregation rules using Terraform.

For easy perusal, the latest set of rules will be saved as `rules.json`.

The`rules.json`, `.terraform.lock.hcl`, and `terraform.tfstate` files will all be committed and pushed to main with the commit message "Auto-apply updated aggregation rules.".

**Note that the expectation is that no Terraform should be run locally.**

## Controlling your recommendations

To control your recommendations, exemptions resources can be added to the `main.tf` file. See [grafana-adaptive-metrics_exemption (Resource)](https://registry.terraform.io/providers/grafana/grafana-adaptive-metrics/latest/docs) for more information.

The new Terraform will be automatically applied when the changes are pushed to `main`.

## Further reading

- [Grafana Adaptive Metrics](https://grafana.com/docs/grafana-cloud/cost-management-and-billing/reduce-costs/metrics-costs/control-metrics-usage-via-adaptive-metrics/)
- [Grafana Adaptive Metrics Terraform provider](https://registry.terraform.io/providers/grafana/grafana-adaptive-metrics/latest/docs)
