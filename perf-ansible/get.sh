NUM_NAMESPACES=${1:-2}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==NAMESPACE [test${i}]"
  echo "  Labels=$(kubectl get ns test${i} -o jsonpath='{.metadata.labels}')"
  echo "  Annotations=$(kubectl get ns test${i} -o jsonpath='{.metadata.annotations}')"
  echo "  Resources:"
  kubectl get sa -n test${i} -l app=kiali -o name
done
