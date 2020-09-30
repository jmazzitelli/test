#!/bin/bash

set -eu

install_pull_secret() {
  if [ "$(oc get secret/pull-secret -n openshift-config -o json | jq -r '.data.".dockerconfigjson"' | base64 -d | jq -M '.auths."brew.registry.redhat.io"')" == "null" ]; then
    echo "Decrypt GPG file to get the password - you will need to provide the GPG passphrase to do this."
    local pw=$(gpg --no-tty --quiet --decrypt ${PASSWORD_GPG_FILE} | grep password: | sed -r 's/^.*password: (.*)/\1/')

    if [ "${pw}" == "" ]; then
      echo "Failed to obtain password - make sure the file [${PASSWORD_GPG_FILE}] has 'password: <the password>' in it."
      exit 1
    fi

    echo "Updating the pull secret so the Brew registry password is now included - will run podman as root via sudo"
    oc get secret/pull-secret -n openshift-config -o json | jq -r '.data.".dockerconfigjson"' | base64 -d > ./authfile
    sudo podman login --authfile ./authfile --username "|shared-qe-temp.zmns.153b77" --password ${pw} brew.registry.redhat.io
    oc set data secret/pull-secret -n openshift-config --from-file=.dockerconfigjson=./authfile

    # We do not need this file anymore
    rm ./authfile
  else
    echo "There is already a pull secret for brew.registry.redhat.io - will not attempt to update it"
  fi
}

install_image_content_source_policy() {
  echo "Creating the ImageContentSourcePolicy for the Brew registry"
  cat <<EOM | oc apply -f -
apiVersion: operator.openshift.io/v1alpha1
kind: ImageContentSourcePolicy
metadata:
  name: brew-registry
spec:
  repositoryDigestMirrors:
  - mirrors:
    - brew.registry.redhat.io
    source: registry.redhat.io
  - mirrors:
    - brew.registry.redhat.io
    source: registry.stage.redhat.io
  - mirrors:
    - brew.registry.redhat.io
    source: registry-proxy.engineering.redhat.com
EOM
}

install_catalog_source() {
  echo "Disabling the out-of-box catalogs"
  oc patch OperatorHub cluster --type json -p '[{"op": "add", "path": "/spec/disableAllDefaultSources", "value": true}]'

  echo "Creating the Brew CatalogSource"
  cat <<EOM | oc apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: ${OPERATOR_CATALOG_PREFIX}-operator-catalog
  namespace: openshift-marketplace
spec:
  sourceType: grpc
  image: brew.registry.redhat.io/rh-osbs/iib:${IIB_TAG}
  displayName: ${OPERATOR_CATALOG_PREFIX^} Operator Catalog
  publisher: grpc
EOM
}

delete_image_content_source_policy() {
  echo "Deleting the ImageContentSourcePolicy for the Brew registry"
  oc delete ImageContentSourcePolicy brew-registry
}

delete_catalog_source() {
  echo "Deleting the Brew CatalogSource"
  oc delete CatalogSource --namespace openshift-marketplace ${OPERATOR_CATALOG_PREFIX}-operator-catalog

  echo "Re-enabling the out-of-box catalogs"
  oc patch OperatorHub cluster --type json -p '[{"op": "add", "path": "/spec/disableAllDefaultSources", "value": false}]'
}

status_pull_secret() {
  echo
  echo "===== PULL SECRET"
  if [ "$(oc get secret/pull-secret -n openshift-config -o json | jq -r '.data.".dockerconfigjson"' | base64 -d | jq -M '.auths."brew.registry.redhat.io"')" == "null" ]; then
    echo "There is no pull secret for brew.registry.redhat.io."
  else
    echo "A pull secret for brew.registry.redhat.io exists."
  fi
}

status_image_content_source_policy() {
  echo
  echo "===== IMAGE CONTENT SOURCE POLICY"
  if oc get ImageContentSourcePolicy brew-registry > /dev/null 2>&1; then
    echo "'brew-registry' ImageContentSourcePolicy is installed:"
    oc get ImageContentSourcePolicy brew-registry
  else
    echo "There is no 'brew-registry' ImageContentSourcePolicy installed."
  fi
}

status_catalog_source() {
  echo
  echo "===== CUSTOM CATALOG SOURCE"
  if oc get CatalogSource --namespace openshift-marketplace ${OPERATOR_CATALOG_PREFIX}-operator-catalog > /dev/null 2>&1; then
    echo "'${OPERATOR_CATALOG_PREFIX}-operator-catalog' CatalogSource is installed:"
    oc get CatalogSource --namespace openshift-marketplace ${OPERATOR_CATALOG_PREFIX}-operator-catalog
    local iib="$(oc get CatalogSource --namespace openshift-marketplace ${OPERATOR_CATALOG_PREFIX}-operator-catalog -o jsonpath='{.spec.image}')"
    echo "The IIB image being used is: ${iib}"
  else
    echo "There is no CatalogSource named [${OPERATOR_CATALOG_PREFIX}-operator-catalog] installed."
  fi

  echo
  echo "===== DEFAULT CATALOG SOURCES"
  if [ "$(oc get OperatorHub cluster -o jsonpath='{.spec.disableAllDefaultSources}')" == "true" ]; then
    echo "All out-of-box catalogs have been disabled as a group."
  else
    echo "All out-of-box catalogs have NOT been disabled as a group."
  fi
}

# Change to the directory where this script is
SCRIPT_ROOT="$( cd "$(dirname "$0")" ; pwd -P )"
cd ${SCRIPT_ROOT}

DEFAULT_PASSWORD_GPG_FILE="./brew-registry.gpg"
DEFAULT_OPERATOR_CATALOG_PREFIX="my"

_CMD=""
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in

    # COMMANDS

    install)
      _CMD="install"
      shift
      ;;

    delete)
      _CMD="delete"
      shift
      ;;

    status)
      _CMD="status"
      shift
      ;;

    # OPTIONS

    -it|--iib-tag)
      IIB_TAG="${2}"
      shift;shift
      ;;

    -ocp|--operator-catalog-prefix)
      OPERATOR_CATALOG_PREFIX="${2}"
      shift;shift
      ;;

    -pgf|--password-gpg-file)
      PASSWORD_GPG_FILE="${2}"
      if [ ! -f ${PASSWORD_GPG_FILE} ]; then
        echo "ERROR: The GPG file is missing. This file must exist and have the password in it: ${PASSWORD_GPG_FILE}"
        exit 1
      fi
      shift;shift
      ;;

    # HELP

    -h|--help)
      cat <<HELPMSG

$0 [option...] command

Valid options:

  -it|--iib-tag <tag>
      The image index location tag (e.g. '14261' in 'registry-proxy.engineering.redhat.com/rh-osbs/iib:14261')
      You get this value from the CVP emails in the "Index Image Location" section.
      Used for "install" and "delete".

  -ocp|--operator-catalog-prefix)
      The prefix used as part of the name of the brew operator catalog that will be created.
      Default: ${DEFAULT_OPERATOR_CATALOG_PREFIX}

  -pgf|--password-gpg-file <path>
      The path to the GPG file that contains the password to the Brew registry.
      Used for "install".
      Default: ${DEFAULT_PASSWORD_GPG_FILE}

The command must be one of:

  * install: Install the necessary resources.
  * delete: Delete any existing resources.
  * status: Provides details about resources that have been installed.

HELPMSG
      exit 1
      ;;
    *)
      echo "ERROR: Unknown argument [$key]. Aborting."
      exit 1
      ;;
  esac
done

# Setup environment

OPERATOR_CATALOG_PREFIX="${OPERATOR_CATALOG_PREFIX:-${DEFAULT_OPERATOR_CATALOG_PREFIX}}"
PASSWORD_GPG_FILE="${PASSWORD_GPG_FILE:-${DEFAULT_PASSWORD_GPG_FILE}}"

echo IIB_TAG=${IIB_TAG:-<unspecified>}
echo OPERATOR_CATALOG_PREFIX=${OPERATOR_CATALOG_PREFIX}
echo PASSWORD_GPG_FILE=$PASSWORD_GPG_FILE

# Make sure we are logged in

if ! oc whoami > /dev/null 2>&1; then
  echo "ERROR: You are not logged into the OpenShift cluster. Use 'oc login' to log into a cluster and then retry."
fi

if [ "${_CMD}" == "install" ]; then

  if [ -z "${IIB_TAG:-}" ]; then
    echo "ERROR: Please specify which IIB tag you want to use (--iib-tag)"
    exit 1
  fi

  install_pull_secret
  install_image_content_source_policy
  install_catalog_source

elif [ "${_CMD}" == "delete" ]; then

  delete_catalog_source
  delete_image_content_source_policy

elif [ "${_CMD}" == "status" ]; then

  status_pull_secret
  status_image_content_source_policy
  status_catalog_source

else 
  echo "ERROR: Missing or unknown command. See --help for usage."
  exit 1
fi
