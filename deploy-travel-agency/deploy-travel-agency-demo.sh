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
