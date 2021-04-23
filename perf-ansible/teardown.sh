NUM_NAMESPACES=${1:-$(cat num-namespaces.txt)}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==DELETE NAMESPACE [test${i}]"
  kubectl delete namespace test${i}
done
