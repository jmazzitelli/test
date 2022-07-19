clean:
	rm -f lottery

build:
	go build lottery.go

run: build
	./lottery -input testdata/testdata.csv -output testdata/lottery-results.csv -numBlinds 5

big-data-run: build
	./lottery -input testdata/big-testdata.csv -output testdata/big-lottery-results.csv -numBlinds 5

build-container: build
	docker build -f ./Dockerfile -t quay.io/jmazzitelli/blind-lottery:v1.0

run-container: build-container
	docker run -v $(PWD)/testdata:/testdata quay.io/jmazzitelli/blind-lottery:v1.0 -input /testdata/testdata.csv -output /testdata/lottery-results.csv
