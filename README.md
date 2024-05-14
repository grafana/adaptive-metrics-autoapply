# Auto-apply for Grafana Adaptive Metrics

A template repository for enabling auto-apply of Adaptive Metrics recommendations.

## Getting started

To start, create a new repository by navigating to "Use this template" → "Create a new repository" at the top right of the repository page.

Add the following variable to the new repository ("Settings" → "Secrets and variables" → "Actions" → "Variables" → " New repository variable")

- `grafana_am_api_url`: This is your Grafana Cloud prometheus URL. To find this URL, go to your `grafana.com` account (https://grafana.com → "My Account") and click on the "Details" button of your Grafana Cloud Prometheus stack.
  The URL is listed at the top of the page next to the Prometheus icon. **Make sure to use only the host part of this URL, e.g. remove anything after `grafana.net`**.

Add the following secret to the new repository ("Settings" → "Secrets and variables" → "Actions" → "New repository secret"):

- `grafana_am_api_key`: This must be specified in the format `<your-numeric-instance-id>:<your-cloud-access-policy-token>`, where
    - `<your-numeric-instance-id>` is the numeric instance ID for which you want to enable auto-apply of your Adaptive Metrics. This value can be found in the "Query Endpoint" section of the Details page under "Username / Instance ID".
    - `<your-cloud-access-policy-token>` is a token from a [Grafana Cloud Access Policy](https://grafana.com/docs/grafana-cloud/account-management/authentication-and-permissions/access-policies/). Make sure the access policy has `metrics:read` and `metrics:write` scopes for the appropriate stack ID.

Once you have added the required variables and secrets, go to "Settings" → "Actions" → "General" → "Workflow permissions" and enable the checkbox for "Allow Github Actions to create and approve pull requests", Then click "Save".

## Enable auto-merge (optional)

Once the above configuration is set up, you can manually run the workflow named "Pull Adaptive Metrics recommendations".
By default, this will create a PR with the current recommendations.
Once you merge this PR, the corresponding set of aggregation rules is automatically created.

To skip the manual PR review and merge step, you can define a repository variable named `grafana_am_automerge_enabled` with value `true`.
This will make the workflow automatically merge the pull request it creates.

## What to expect

By default, auto-apply is scheduled to run at 04:00 UTC Monday to Friday. This can be configured by editing the schedule time in `.github/workflows/pull_recommendations.yml`.
At the scheduled time, the GitHub Action will pull the latest recommendations and create (and if enabled merge) a pull request titled "Scheduled refresh of the latest recommendations." with the changes.

Once the pull request is merged, the`rules.json`, `.terraform.lock.hcl`, and `terraform.tfstate` files will be created or updated and pushed to main with the commit message "Auto-apply updated aggregation rules.".

## Controlling your recommendations

To control your recommendations, exemptions resources can be added to the `main.tf` file. See [grafana-adaptive-metrics_exemption (Resource)](https://registry.terraform.io/providers/grafana/grafana-adaptive-metrics/latest/docs) for more information.

The new Terraform will be automatically applied when the changes are pushed to `main`.

## Further reading

- [Grafana Adaptive Metrics](https://grafana.com/docs/grafana-cloud/cost-management-and-billing/reduce-costs/metrics-costs/control-metrics-usage-via-adaptive-metrics/)
- [Grafana Adaptive Metrics Terraform provider](https://registry.terraform.io/providers/grafana/grafana-adaptive-metrics/latest/docs)
