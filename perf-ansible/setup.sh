NUM_NAMESPACES=${1:-2}
for i in $(seq ${NUM_NAMESPACES}); do
  echo "==CREATE NAMESPACE [test${i}]"
  kubectl create namespace test${i}
done
