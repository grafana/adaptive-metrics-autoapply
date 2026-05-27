# Automatically apply recommendations from Adaptive Metrics

Use this template repository to set up auto-apply mode for Adaptive Metrics in Grafana Cloud.

> [!NOTE]
> [Server-side auto-apply mode for Adaptive Metrics](https://grafana.com/docs/grafana-cloud/cost-management-and-billing/adaptive-telemetry/adaptive-metrics/auto-apply/) is now available in [public preview](https://grafana.com/docs/release-life-cycle/). We encourage customers to use server-side auto-apply, as this is what Grafana Labs will be building upon going forward. Grafana Labs will plan to deprecate this version of auto-apply in the future.  

## What to expect

By default, auto-apply mode runs at 04:00 UTC Monday through Friday. To configure this setting, edit the `schedule` parameter in the `.github/workflows/pull_recommendations.yml` file.

At the scheduled time, the GitHub Action pulls the latest recommendations and creates a pull request named "Scheduled refresh of the latest recommendations".

After you merge this pull request, the GitHub Action uploads the updated rules to Grafana Cloud.

You can also set the pull request to merge automatically.

## Automatically apply recommendations

Create a new repository using this one as a template to automatically apply Adaptive Metrics recommendations in Grafana Cloud.

1. Create a new repository by navigating to "Use this template" → "Create a new repository" at the top-right of the repository page in GitHub.

2. Go to "Settings" → "Secrets and variables" → "Actions" → "Variables" → " New repository variable" and add the following variable to the new repository:

    - `grafana_am_api_url`: This is your Grafana Cloud Prometheus URL. To find this URL, go to your `grafana.com` account (https://grafana.com → "My Account") and click on the "Details" button of your Grafana Cloud Prometheus stack.
  The URL is listed at the top of the page next to the Prometheus icon. 
      > Make sure to use only the host part of this URL. Remove any parameters after `grafana.net`.

3. Go to "Settings" → "Secrets and variables" → "Actions" → "New repository secret" and add the following secret to the new repository:

    - `grafana_am_api_key`: You must specify this key in the format `<your-numeric-instance-id>:<your-cloud-access-policy-token>`, where:
      - `<your-numeric-instance-id>` is the numeric instance ID for which you want to enable auto-apply mode. You can find this value in the "Query Endpoint" section of the *Details* page under "Username / Instance ID".
      - `<your-cloud-access-policy-token>` is a token from a [Grafana Cloud Access Policy](https://grafana.com/docs/grafana-cloud/account-management/authentication-and-permissions/access-policies/). Make sure the access policy has `metrics:read` and `metrics:write` scopes for the appropriate stack ID.

4. Go to "Settings" → "Actions" → "General" → "Workflow permissions" and select the checkbox for "Allow GitHub Actions to create and approve pull requests". Then, click "Save".

After you set up this configuration, you can manually run the workflow named "Pull Adaptive Metrics recommendations".
By default, this workflow creates a pull request with the latest recommendations.
After you merge this pull request, the workflow automatically creates the corresponding set of aggregation rules.

## (Optional) Automatically merge rules

You can enable auto-merge mode to skip the manual pull request review and merge processes.

1. Create a [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens). The personal access token should have access to the repo and read/write permissions for "Pull Requests" and "Contents" enabled.

2. Go to "Settings" → "Secrets and variables" → "Actions" → "New repository secret" and add the following secret to the new repository:

    - `automerge_pat`: This is the personal access token you created in the previous step.

## See also

- [Grafana Adaptive Metrics](https://grafana.com/docs/grafana-cloud/cost-management-and-billing/reduce-costs/metrics-costs/control-metrics-usage-via-adaptive-metrics/)
- [Grafana Adaptive Metrics Terraform provider](https://registry.terraform.io/providers/grafana/grafana-adaptive-metrics/latest/docs)

Have feedback about this feature? Open an issue to let the team know!
