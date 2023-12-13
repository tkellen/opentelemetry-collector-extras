# Metrics as Attributes Processor

## Description

The metrics as attributes processor is used to add infrastructure metrics as
attributes to spans and logs. This processor will cache metrics from the metrics
pipeline that match a given selector. This cache will be shared with the traces
and logs pipelines. If a span or log record matches a target selector, the
cached metrics will be added as attributes to the span or log record.

## Caveats

The metrics and spans (or logs) need to be routed through the same OpenTelemetry
Collector instance for this processor to work. 

It is recommended you have an OpenTelemetry Collector deployed in an agent mode
near or on the host that will emit metrics. Spans (or logs) that need to have
metrics attributes will need to be routed through this same Collector instance.

For shared components such as a message queue or database, it is recommended to
have a Collector instance deployed near the shared component to collect metrics.
Spans (or logs) that need to have the metrics attributes added will be routed to
this Collector instance. You can do the filter processor, and mulitple pipelines
to route spans (or logs) accordingly.

> **Warning:**
> Using this processor with a Collector instance that is deployed in a cluster
> behind a load balancer is not recommended.

## Configuration

Configuration is specified through a list of metrics groups. Each metrics group
has target selectors for spans and/or logs and metric selectors to identify
how matches are to be made. It also includes a metrics to add section used to 
specify which metrics will be cached and added.

```yaml
  metricsasattributes:
    # time before state metrics are removed from cache (default 5m)
    cache_ttl: 5m

    metrics_groups:
      # name of the group (must be unique across all groups)
      - name: arbritrary_name

        # This will create a single unique string with all the attributes
        # order is important
        # e.g. k8s.pod.name="pod1",k8s.namespace.name="namespace1" 
        # would result in pod1,namespace1
        target_selectors:
          spans:
            - attribute_type: resource
              name: k8s.pod.name
            - attribute_type: resource
              name: k8s.namespace.name
          logs:
            - attribute_type: resource
              name: k8s.pod.name
            - attribute_type: resource
              name: k8s.namespace.name

        # This will create a single unique string with all the attributes
        # order is important
        # e.g. k8s.pod.name="pod1",k8s.namespace.name="namespace1" 
        # would result in pod1,namespace1
        metrics_selector:
          - attribute_type: resource
            name: k8s.pod.name
          - attribute_type: resource
            name: k8s.namespace.name

        metrics_to_add:
          - instrumentation_scope: otelcol/hostmetrics*
            metrics:
              - name: system.cpu.usage
                include_only_attributes: {cpu: total, state: user}
              - name: system.memory.usage
                include_only_attributes: {state: used}
                new_name: system.memory.usage.used
          - instrumentation_scope: jvm
            metrics:
              - name: jvm.*
              
      # - name: metricsGroup2
      # ...

```

