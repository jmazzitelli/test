NUM_NAMESPACES=${1:-$(cat num-namespaces.txt)}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==NAMESPACE [test${i}]"
  kubectl delete sa -n test${i} -l app=kiali
  kubectl label namespace test${i} kiali.io/member-of-
done
