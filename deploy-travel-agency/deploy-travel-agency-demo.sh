#/bin/bash

: ${CLIENT_EXE:="oc"}

# Create the demo namespaces

${CLIENT_EXE} create namespace travel-agency
${CLIENT_EXE} label namespace travel-agency istio-injection=enabled

${CLIENT_EXE} create namespace travel-portal
${CLIENT_EXE} label namespace travel-portal istio-injection=enabled

# Prepare the new demo namespaces for CNI

if [ "${CLIENT_EXE}" == "oc" ]; then
${CLIENT_EXE} adm policy add-scc-to-group privileged system:serviceaccounts:travel-agency
${CLIENT_EXE} adm policy add-scc-to-group anyuid system:serviceaccounts:travel-agency
cat <<EOF | ${CLIENT_EXE} -n travel-agency create -f -
apiVersion: "k8s.cni.cncf.io/v1"
kind: NetworkAttachmentDefinition
metadata:
  name: istio-cni
EOF

${CLIENT_EXE} adm policy add-scc-to-group privileged system:serviceaccounts:travel-portal
${CLIENT_EXE} adm policy add-scc-to-group anyuid system:serviceaccounts:travel-portal
cat <<EOF | ${CLIENT_EXE} -n travel-portal create -f -
apiVersion: "k8s.cni.cncf.io/v1"
kind: NetworkAttachmentDefinition
metadata:
  name: istio-cni
EOF
fi

# Deploy the demo

${CLIENT_EXE} apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/master/travel_agency.yaml) -n travel-agency
${CLIENT_EXE} apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/master/travel_portal.yaml) -n travel-portal

# Set up metric classification

${CLIENT_EXE} -n istio-system get envoyfilter stats-filter-1.6 -o yaml > stats-filter-1.6.yaml
cat <<EOF | patch -o - | ${CLIENT_EXE} -n istio-system apply -f - && rm stats-filter-1.6.yaml
--- stats-filter-1.6.yaml	2020-06-02 11:10:29.476537126 -0400
+++ stats-filter-1.6.yaml.new	2020-06-02 09:59:26.434300000 -0400
@@ -95,7 +95,14 @@
               configuration: |
                 {
                   "debug": "false",
-                  "stat_prefix": "istio"
+                  "stat_prefix": "istio",
+                  "metrics": [
+                   {
+                     "name": "requests_total",
+                     "dimensions": {
+                       "request_operation": "istio_operationId"
+                     }
+                   }]
                 }
               root_id: stats_inbound
               vm_config:
EOF

cat <<EOF | ${CLIENT_EXE} -n istio-system apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: attribgen-travelagency-hotels
spec:
  workloadSelector:
    labels:
      app: hotels
  configPatches:
  - applyTo: HTTP_FILTER
    match:
      context: SIDECAR_INBOUND
      proxy:
        proxyVersion: '1\.6.*'
      listener:
        filterChain:
          filter:
            name: "envoy.http_connection_manager"
            subFilter:
              name: "istio.stats"
    patch:
      operation: INSERT_BEFORE
      value:
        name: istio.attributegen
        typed_config:
          "@type": type.googleapis.com/udpa.type.v1.TypedStruct
          type_url: type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm
          value:
            config:
              configuration: |
                {
                  "attributes": [
                    {
                      "output_attribute": "istio_operationId",
                      "match": [
                        {
                          "value": "LocalRentals",
                          "condition": "request.url_path.matches('^/hotels/[[:alnum:]]*$') && request.method == 'GET'"
                        }
                      ]
                    }
                  ]
                }
              vm_config:
                runtime: envoy.wasm.runtime.null
                code:
                  local: { inline_string: "envoy.wasm.attributegen" }
---
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: attribgen-travelagency-cars
spec:
  workloadSelector:
    labels:
      app: cars
  configPatches:
  - applyTo: HTTP_FILTER
    match:
      context: SIDECAR_INBOUND
      proxy:
        proxyVersion: '1\.6.*'
      listener:
        filterChain:
          filter:
            name: "envoy.http_connection_manager"
            subFilter:
              name: "istio.stats"
    patch:
      operation: INSERT_BEFORE
      value:
        name: istio.attributegen
        typed_config:
          "@type": type.googleapis.com/udpa.type.v1.TypedStruct
          type_url: type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm
          value:
            config:
              configuration: |
                {
                  "attributes": [
                    {
                      "output_attribute": "istio_operationId",
                      "match": [
                        {
                          "value": "LocalRentals",
                          "condition": "request.url_path.matches('^/cars/[[:alnum:]]*$') && request.method == 'GET'"
                        }
                      ]
                    }
                  ]
                }
              vm_config:
                runtime: envoy.wasm.runtime.null
                code:
                  local: { inline_string: "envoy.wasm.attributegen" }
---
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: attribgen-travelagency-flights
spec:
  workloadSelector:
    labels:
      app: flights
  configPatches:
  - applyTo: HTTP_FILTER
    match:
      context: SIDECAR_INBOUND
      proxy:
        proxyVersion: '1\.6.*'
      listener:
        filterChain:
          filter:
            name: "envoy.http_connection_manager"
            subFilter:
              name: "istio.stats"
    patch:
      operation: INSERT_BEFORE
      value:
        name: istio.attributegen
        typed_config:
          "@type": type.googleapis.com/udpa.type.v1.TypedStruct
          type_url: type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm
          value:
            config:
              configuration: |
                {
                  "attributes": [
                    {
                      "output_attribute": "istio_operationId",
                      "match": [
                        {
                          "value": "LongDistanceTransportation",
                          "condition": "request.url_path.matches('^/flights/[[:alnum:]]*$') && request.method == 'GET'"
                        }
                      ]
                    }
                  ]
                }
              vm_config:
                runtime: envoy.wasm.runtime.null
                code:
                  local: { inline_string: "envoy.wasm.attributegen" }
EOF
