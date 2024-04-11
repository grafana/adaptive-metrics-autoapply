terraform {
  required_providers {
    grafana-adaptive-metrics = {
      source = "registry.terraform.io/grafana/grafana-adaptive-metrics"
    }
  }
}

provider "grafana-adaptive-metrics" {}

locals {
  recommendations = fileexists("recommendations.json") ? jsondecode(file("recommendations.json")) : []
}

resource "grafana-adaptive-metrics_rule" "rules" {
  for_each = {
    for rec in local.recommendations : rec.metric => rec
  }

  metric               = each.value.metric
  match_type           = lookup(each.value, "match_type", "")
  drop                 = lookup(each.value, "drop", false)
  drop_labels          = lookup(each.value, "drop_labels", [])
  keep_labels          = lookup(each.value, "keep_labels", [])
  aggregations         = lookup(each.value, "aggregations", [])
  aggregation_interval = lookup(each.value, "aggregation_interval", "")
  aggregation_delay    = lookup(each.value, "aggregation_delay", "")

  # We set auto_import=true to tell the provider to automatically import any
  # existing rules into Terraform state.
  auto_import = true
}

output "rules" {
  value = resource.grafana-adaptive-metrics_rule.rules
}
