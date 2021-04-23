NUM_NAMESPACES=${1:-$(cat num-namespaces.txt)}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==CREATE NAMESPACE [test${i}]"
  kubectl create namespace test${i}
done
