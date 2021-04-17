NUM_NAMESPACES=${1:-2}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==NAMESPACE [test${i}]"
  kubectl get sa -n test${i} -l app=kiali -o name
done
