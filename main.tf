terraform {
  required_providers {
    grafana-adaptive-metrics = {
      source = "registry.terraform.io/grafana/grafana-adaptive-metrics"
    }
  }
}

provider "grafana-adaptive-metrics" {}

data "grafana-adaptive-metrics_recommendations" "all" {
  verbose = false
}

output "recommendations" {
  value = data.grafana-adaptive-metrics_recommendations.all.recommendations
}

resource "grafana-adaptive-metrics_rule" "recommendations" {
  for_each = {
    for rec in data.grafana-adaptive-metrics_recommendations.all.recommendations : rec.metric => rec
  }

  metric               = each.value.metric
  match_type           = each.value.match_type
  drop                 = each.value.drop
  drop_labels          = each.value.drop_labels
  keep_labels          = each.value.keep_labels
  aggregations         = each.value.aggregations
  aggregation_interval = each.value.aggregation_interval
  aggregation_delay    = each.value.aggregation_delay
}
