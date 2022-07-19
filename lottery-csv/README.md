# Hunt Blind Lottery

A sample Go app/container that processes hunter lottery applications and randomly assigns numbered blinds to the hunters who have applied for the lottery.

The input data consists of a CSV file whose first line contains the column headers and subsequent lines represent one hunter application per row.

You can have this data in an Excel, LibreOffice, or Google spreadsheet - just export the data as a CSV (comma-separated values) file and use that.

The first column is "CID" or the hunter's ID number.
The second column is the hunter's name.
The remaining columns are each date for which the blinds can be awarded. There can be any number of hunt date columns.
Each subsequent row represents a single hunter's application, with a "1" in a hunt date column in which the hunter wishes to enter the lottery for that date.
Only one application may be submitted per hunter, though they can request any number of dates in the lottery.

For example, a spreadsheet that looks like this:

|CID|Name|10/15/2022|10/19/2022|10/20/2022|01/02/2023|01/03/2023|
|---|----|----------|----------|----------|----------|----------|
|000-123-456|Bob Jones|1|1|1|1|
|321-987-654|John Smith|1|1| |1|

will be represented in a CSV file like this:

```
CID,Name,10/15/2022,10/19/2022,10/20/2022,01/02/2023,01/03/2023
000-123-456,Bob Jones,1,1,1,1,,
321-987-654,John Smith,1,1,,1,,
```

Once all lottery applications are entered in the spreadsheet and a CSV file is available in a file called "data.csv", run the program like this:

```bash
make build
./lottery -input ./data.csv -output ./lottery-results.csv
```

If you can't build the app on your machine, you can run the container image like this (requires that you have [Docker](https://docs.docker.com/get-docker/) installed on your machine):

```bash
docker run -v .:/testdata quay.io/jmazzitelli/blind-lottery:v1.0 -input /testdata/data.csv -output /testdata/lottery-results.csv
```

This will output a file `lottery-results.csv` which can be imported into your spreadsheet program of choice.
It will contain all the assignments for each blind on each hunt date.
An example lottery result will look something like this:

```
Date,Blind,CID,Name
10/15/2022,1,000-123-456,Bob Jones
10/15/2022,2,321-987-654,John Smith
...
```

When imported into a spreadsheet, it will look something like this:

|Date|Blind|CID|Name|
|----|-----|---|----|
|10/15/2022|1|000-123-456|Bob Jones|
|10/15/2022|2|321-987-654|John Smith|
...

Note that by default the number of blinds assigned per hunt date is 5, though
you can change that via the `-numBlinds` command line argument.

# Lottery Selection Algorithm

The lottery algorithm aims to assign every hunter to at least one blind, but
also aims to assign all the blinds on any given hunt date. It also attempts to
be fair - if hunters have been assigned several blinds over and above other
hunters, blinds will be reassigned to those hunters that have fewer.

For each hunt date, each blind will be assigned randomly by choosing one hunter
from a pool of all the hunters that have requested a hunt for that date.  Once
a hunter has been awarded a blind, that hunter will be removed from the pool of
available hunters thus giving everyone else a chance to win a blind.  The only
time a hunter will be placed back into the pool is if: (a) Everyone has been
assigned a blind, or (b) There is an empty pool of hunters for a particular
date (due to all hunters either not having selected that date or having already
been assigned a blind).  Therefore, it is possible for a hunter to be assigned
a blind on multiple hunt dates, but only after all other hunters have been
assigned a blind OR there are too few available hunters to fill up all the
blinds on a given hunt date.  Once a hunter has been assigned a blind, the
number of blinds that hunter has been awarded will be compared to other hunters
who have been awarded a fewer number of blinds. If the difference is 2 or
greater (that is, one hunter has been awarded 2 or more blinds over another
hunter), one of those "surplus" blinds will be re-assigned to the hunter that
has 2 or more fewer blinds.

A hunter will not be awarded multiple blinds on the same hunt date. If too few
hunters are available to hunt all blinds for a given hunt date, some blinds
will remain unassigned.

## Building The Program

* Run `make build` to build the Go app.
* Run `make run` to run the Go app and process the test data found in the `testdata` directory.
* Run `make build-container` to build the container image.
* Run `make run-container` to run the Go app from within the container image and process the test data found in the `testdata` directory.
