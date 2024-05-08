# Auto-apply for Grafana Adaptive Metrics

A template repository for enabling auto-apply of Adaptive Metrics recommendations.

## Getting started

To start, create a new repository by navigating to "Use this template" → "Create a new repository" at the top right of the repository page.

Add the following variable to the new repository ("Settings" → "Secrets and variables" → "Actions" → "Variables" → " New repository variable")

- `grafana_am_api_url`: This is your Grafana Cloud prometheus URL. To find this URL, go to your `grafana.com` account (https://grafana.com → "My Account") and click on the "Details" button of your Grafana Cloud Prometheus stack.
  The URL is listed at the top of the page next to the Prometheus icon. **Make sure to use only the host part of this URL, e.g. remove anything after `grafana.net`**.

Add the following secrets to the new repository ("Settings" → "Secrets and variables" → "Actions" → "New repository secret"):

- `grafana_am_api_key`: This must be specified in the format `<your-numeric-instance-id>:<your-cloud-access-policy-token>`, where
    - `<your-numeric-instance-id>` is the numeric instance ID for which you want to enable auto-apply of your Adaptive Metrics. This value can be found in the "Query Endpoint" section of the Details page under "Username / Instance ID".
    - `<your-cloud-access-policy-token>` is a token from a [Grafana Cloud Access Policy](https://grafana.com/docs/grafana-cloud/account-management/authentication-and-permissions/access-policies/). Make sure the access policy has `metrics:read` and `metrics:write` scopes for the appropriate stack ID.
- `automerge_pat`: Optionally, provide a [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) to enable auto-apply of new recommendations. Without this, a PR will be opened with the latest recommendations, but that PR will still require a manual merge to take effect. **The personal access token should have access to the repo and read/write permissions for "Pull Requests" and "Contents" enabled.**

## What to expect

By default, auto-apply is scheduled to run at 04:00 UTC every day. This can be configured by editing the schedule time in `.github/workflows/main.yml`.

At the scheduled time, the GitHub Action will pull the latest recommendations and open a pull request titled "Scheduled refresh of the latest recommendations." with the changes.

If `automerge_pat` was provided, the pull request will automatically be merged. If not, the pull request will remain open for manual review.

Once the pull request is merged, the`rules.json`, `.terraform.lock.hcl`, and `terraform.tfstate` files will all be committed and pushed to main with the commit message "Auto-apply updated aggregation rules.".

## Controlling your recommendations

To control your recommendations, exemptions resources can be added to the `main.tf` file. See [grafana-adaptive-metrics_exemption (Resource)](https://registry.terraform.io/providers/grafana/grafana-adaptive-metrics/latest/docs) for more information.

The new Terraform will be automatically applied when the changes are pushed to `main`.

## Further reading

- [Grafana Adaptive Metrics](https://grafana.com/docs/grafana-cloud/cost-management-and-billing/reduce-costs/metrics-costs/control-metrics-usage-via-adaptive-metrics/)
- [Grafana Adaptive Metrics Terraform provider](https://registry.terraform.io/providers/grafana/grafana-adaptive-metrics/latest/docs)
