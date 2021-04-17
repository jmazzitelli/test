NUM_NAMESPACES=${1:-2}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==DELETE NAMESPACE [test${i}]"
  kubectl delete namespace test${i}
done
