NUM_NAMESPACES=${1:-2}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==NAMESPACE [test${i}]"
  kubectl delete sa -n test${i} -l app=kiali
done
