terraform {
  required_providers {
    grafana-adaptive-metrics = {
      source  = "registry.terraform.io/grafana/grafana-adaptive-metrics"
      version = "0.3.0-alpha.5"
    }
  }
}

provider "grafana-adaptive-metrics" {}

locals {
  recommendations = jsondecode(fileexists("${path.module}/recommendations.json") ? file("${path.module}/recommendations.json") : "[]")

  segments = jsondecode(fileexists("${path.module}/segments.json") ? file("${path.module}/segments.json") : "[]")
  segments_lookup = {
    for segment in local.segments : segment.name => segment
  }

  segment_recommendations_files = fileset(path.module, "recommendations-*.json")
  segment_recommendations = {
    for f in local.segment_recommendations_files : regex("recommendations-(.+).json", f)[0] => jsondecode(file("${path.module}/${f}"))
  }
}

resource "grafana-adaptive-metrics_ruleset" "default_segment_rules" {
  rules = [
    for rule in local.recommendations : {
      metric               = rule.metric
      match_type           = lookup(rule, "match_type", "")
      drop                 = lookup(rule, "drop", false)
      drop_labels          = lookup(rule, "drop_labels", [])
      keep_labels          = lookup(rule, "keep_labels", [])
      aggregations         = lookup(rule, "aggregations", [])
      aggregation_interval = lookup(rule, "aggregation_interval", "")
      aggregation_delay    = lookup(rule, "aggregation_delay", "")
    }
  ]
}

resource "grafana-adaptive-metrics_ruleset" "segmented_rules" {
  for_each = local.segment_recommendations

  segment = local.segments_lookup[each.key].id
  rules = [
    for rule in each.value : {
      metric               = rule.metric
      match_type           = lookup(rule, "match_type", "")
      drop                 = lookup(rule, "drop", false)
      drop_labels          = lookup(rule, "drop_labels", [])
      keep_labels          = lookup(rule, "keep_labels", [])
      aggregations         = lookup(rule, "aggregations", [])
      aggregation_interval = lookup(rule, "aggregation_interval", "")
      aggregation_delay    = lookup(rule, "aggregation_delay", "")
    }
  ]
}
