package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"
)

// A CID is a hunter identification number - it stands for "conservation identification".
// It uniquely identifies a hunter and his lottery application.
type CID string

// ApplicationType represents a hunter's application for the lottery.
// It contains the hunter's CID, name, and the hunt dates being applied for.
// The array of requested dates must be of the same length as the number of hunt dates being awarded in the lottery.
// Note that the application does not specify requested blinds - hunters are awarded a date with a blind randomly assigned.
type ApplicationType struct {
	CID            CID
	Name           string
	RequestedDates []bool
}

// LotteryResultType represents one lottery award that was won by a hunter.
// It contains the application that represents the hunter that submitted the application.
// It also contains the hunt date awarded along with the blind assigned for that hunt date.
type LotteryResultType struct {
	Application ApplicationType
	AwardedDate string
	BlindNumber int
}

// CIDMapType is a map of CIDs. This provides objects that need to contain a set of CIDs that are easy to lookup.
type CIDMapType map[CID]bool

// ApplicationByCIDType is a map of applications keyed on CID, thus making an application easy to lookup by CID.
type ApplicationByCIDType map[CID]ApplicationType

// ApplicationListType is simply an array consisting of a list of applications.
type ApplicationListType []ApplicationType

// contains looks up a CID in the given map and returns true if the CID exists.
func (m CIDMapType) contains(cid CID) bool {
	_, ok := m[cid]
	return ok
}

// contains looks up a CID in the given map and returns true if the CID exists.
func (m ApplicationByCIDType) contains(cid CID) bool {
	_, ok := m[cid]
	return ok
}

// contains looks up a CID in the given application list and returns true if an application for the CID exists.
func (apps ApplicationListType) contains(cid CID) bool {
	for _, app := range apps {
		if app.CID == cid {
			return true
		}
	}
	return false
}

// didCIDRequestDate determines if an application with the given CID requested to hunt on the given date.
// In other words, did the hunter apply for the lottery for the given hunt date?
func didCIDRequestDate(cid CID, huntDateToSearch string) bool {
	app := applicationsByCID[cid]
	for huntDateIndex, huntDate := range huntDates {
		if huntDate == huntDateToSearch {
			return app.RequestedDates[huntDateIndex]
		}
	}
	return false
}

// wasCIDAwardedBlindOnDate determines if an application with the given CID was assigned a blind on the given date.
// In other words, did the hunter win the lottery for the given hunt date?
func wasCIDAwardedBlindOnDate(cid CID, huntDateToSearch string) bool {
	for _, l := range lotteryResults {
		if cid == l.Application.CID && huntDateToSearch == l.AwardedDate {
			return true
		}
	}
	return false
}

const DEFAULT_INPUT_FILE_NAME string = "./data.csv"
const DEFAULT_OUTPUT_FILE_NAME string = "./lottery-results.csv"
const DEFAULT_NUMBER_OF_BLINDS int = 5

// NULL_APPLICATION is used when a blind was not awarded to any hunter - the blind has gone unassigned.
var NULL_APPLICATION = ApplicationType{
	CID:            "",
	Name:           "",
	RequestedDates: make([]bool, 0),
}

// applicationsList is the list of all applications as found in the input CSV file.
var applicationsList ApplicationListType

// applicationsByCID represents the same data as applicationsList but in a map form keyed on CID so applications can be looked up by  CID.
var applicationsByCID ApplicationByCIDType

// huntDates are a list of dates in which the set of blinds will be awarded in the lottery.
// These are taken directly from the column headers in the input CSV file.
var huntDates []string

// lotteryResults contain the results of the lottery winnings - each item in the list represents an awarded blind on a given hunt date.
var lotteryResults []*LotteryResultType

// selectedCIDs are a set of CIDs representing hunters that have been awarded a blind.
// This is used during the lottery processing to ensure fairness when picking lottery winners.
var selectedCIDs CIDMapType

// CIDAwards track how many blinds were awarded to a given CID. The Awards field is a counter that indicates
// how many blinds were awarded to the CID. This struct is designed in such a way as to be sortable.
type CIDAwards struct {
	CID    CID
	Awards int
}

// CIDAwardsList is a list that provides information on how many blinds were awarded to the different CIDs.
type CIDAwardsList []*CIDAwards

// These three methods provide the Interface for the Sort API.
func (c CIDAwardsList) Len() int           { return len(c) }
func (c CIDAwardsList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CIDAwardsList) Less(i, j int) bool { return c[i].Awards < c[j].Awards }

// incrementCIDAward will be called when a CID was awarded another blind.
// The CID's Awards counter will be incremented.
func incrementCIDAward(cid CID) {
	for _, item := range awardsByCID {
		if item.CID == cid {
			item.Awards++
		}
	}
	sort.Sort(awardsByCID)
}

// decrementCIDAward will be called when a CID was stripped of a previously awarded blind in order to transfer that blind to another CID.
// The CID's Awards counter will be decremented.
func decrementCIDAward(cid CID) {
	for _, item := range awardsByCID {
		if item.CID == cid && item.Awards > 0 {
			item.Awards--
		}
	}
	sort.Sort(awardsByCID)
}

// currentCIDAwards returns the current number of blinds that have been awarded to the given CID.
func currentCIDAwards(cid CID) int {
	for _, item := range awardsByCID {
		if item.CID == cid {
			return item.Awards
		}
	}
	return 0
}

// The list of CIDs and the number of blinds awarded to each CID.
var awardsByCID CIDAwardsList

// command line options passed by the user
var (
	argInput     = flag.String("input", DEFAULT_INPUT_FILE_NAME, "Path to the input CSV file.")
	argOutput    = flag.String("output", DEFAULT_OUTPUT_FILE_NAME, "Path to the output CSV file containing the lottery results.")
	argNumBlinds = flag.Int("numBlinds", DEFAULT_NUMBER_OF_BLINDS, "Number of blinds to be awarded per hunt day.")
)

// processArguments will parse the command line options passed to the program by the user.
func processArguments() {
	flag.Parse()

	if *argNumBlinds <= 0 {
		log.Fatalf("Number of blinds must be greater than 0")
	}

	log.Printf("Input file: %v", *argInput)
	log.Printf("Output file: %v", *argOutput)
	log.Printf("Number of blinds: %v", *argNumBlinds)
}

// readApplicationData will read the input CSV file and parse the application data.
func readApplicationData() {
	// open data file
	f, err := os.Open(*argInput)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read all csv data
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// initialize empty lists and maps
	applicationsByCID = make(ApplicationByCIDType, 0)
	selectedCIDs = make(CIDMapType, 0)
	awardsByCID = make(CIDAwardsList, 0)

	// process all rows in the application CSV file - the first row is the header row, the rest of the rows are each individual hunter's application data.
	for rowCount, row := range data {
		numberOfNonHuntDateColumns := 2 // non-hunt-date column headers are: "CID", "Name" (the remaining column headers are the actual hunt dates)
		if rowCount == 0 {
			huntDates = make([]string, len(row)-numberOfNonHuntDateColumns)
			for i, huntDate := range row {
				if i >= numberOfNonHuntDateColumns {
					huntDates[i-numberOfNonHuntDateColumns] = huntDate
				}
			}
		} else {
			// if the row doesn't have the CID column, name column, and boolean selector columns for all the hunt dates then abort - the data is bad
			if len(row) != numberOfNonHuntDateColumns+len(huntDates) {
				log.Fatalf("The row does not have the expected number of columns: %+v", row)
			}
			// build the current Application object from the current row
			var rec ApplicationType
			for i, column := range row {
				switch i {
				case 0:
					rec.CID = CID(column)
					rec.RequestedDates = make([]bool, len(huntDates))
				case 1:
					rec.Name = column
				default:
					if column != "1" && column != "0" && column != "" {
						log.Fatalf("Application for CID [%v] has invalid hunt date values. Must be '0', '1', or empty string ''.\n", rec.CID)
					}
					rec.RequestedDates[i-numberOfNonHuntDateColumns] = (column == "1") // if column is "1" that means the hunter wants to enter the lottery for that date
				}
			}
			if applicationsByCID.contains(rec.CID) {
				log.Fatalf("There are multiple applications with the same CID [%v] - aborting", rec.CID)
			}

			// we got another hunter's application - store it in the appropriate places
			applicationsList = append(applicationsList, rec)
			applicationsByCID[rec.CID] = rec
			awardsByCID = append(awardsByCID, &CIDAwards{
				CID:    rec.CID,
				Awards: 0,
			})
		}
	}

	log.Printf("===== START INPUT DATA =====")
	for i, j := range huntDates {
		fmt.Printf("Hunt Day #%v: %+v\n", i, j)
	}

	for i, j := range applicationsList {
		fmt.Printf("Hunter Application #%v: %+v\n", i, j)
	}
	log.Printf("===== END INPUT DATA =====")
}

// performLotteryDraw will randomly award all blinds for all hunt dates to all hunters who have applied for the lottery.
func performLotteryDraw() {
	log.Println("RANDOMLY ASSIGNING BLINDS NOW")
	rand.Seed(time.Now().UnixNano())

	// the main loop processes each hunt date - the inner loop will then award each blind for the current hunt date
	for huntDateIndex, huntDate := range huntDates {
		log.Printf("START LOTTERY FOR HUNT DATE #%v: [%v]\n", huntDateIndex, huntDate)

		// determine the pool of applications to select from for this date (only include those that haven't been selected yet and requested this date)
		datePool := make(ApplicationListType, 0)
		for _, app := range applicationsList {
			if app.RequestedDates[huntDateIndex] == true && !selectedCIDs.contains(app.CID) {
				datePool = append(datePool, app)
			}
		}
		log.Printf("Number of applications to choose from for the hunt date is [%v]\n", len(datePool))

		// keep track of who was selected for the current hunt date - ensures we don't select the same application twice this date after resetting the main pool
		selectedFromDatePool := make(CIDMapType, 0)

		// the inner loop will randomly assign each blind for the current hunt date
		for blindNum := 1; blindNum <= *argNumBlinds; blindNum++ {
			log.Printf("START LOTTERY FOR HUNT DATE [%v], BLIND [%v]\n", huntDate, blindNum)

			// if everyone has been awarded at least one blind, reset and put everyone back into the main pool
			if len(selectedCIDs) == len(applicationsList) {
				selectedCIDs = make(CIDMapType, 0)
				log.Printf("All applications have been awarded blinds. Resetting as of HUNT DATE [%v], BLIND [%v]\n", huntDate, blindNum)
			}

			// check to see if everyone in the current hunt date pool has been selected - if so we may have run out of hunters to award the current blind
			if len(datePool) == len(selectedFromDatePool) {
				log.Printf("All applications from the current pool have been selected for HUNT DATE [%v], BLIND [%v]\n", huntDate, blindNum)

				// We may have another application that was previously selected who also wants this date.
				// We will allow this application to enter the pool again in order to avoid an unclaimed blind.
				// But don't put the application back in the pool if it was already selected for this hunt date because we can't award multiple blinds
				// for the same hunt date to the same application.
				for _, app := range applicationsList {
					if !datePool.contains(app.CID) && app.RequestedDates[huntDateIndex] == true && !selectedFromDatePool.contains(app.CID) {
						log.Printf("Adding application [%v] back into the pool to avoid an unclaimed blind for HUNT DATE [%v]\n", app.CID, huntDate)
						datePool = append(datePool, app)
					}
				}

				// If we still don't have enough applications in the pool, this blind is going to have to go unclaimed.
				if len(datePool) == len(selectedFromDatePool) {
					log.Printf("No more applications left to select - HUNT DATE [%v], BLIND [%v] will go unclaimed!\n", huntDate, blindNum)
					lotteryResults = append(lotteryResults, &LotteryResultType{
						Application: NULL_APPLICATION,
						AwardedDate: huntDate,
						BlindNumber: blindNum,
					})
					continue
				}
			}

			// randomly pick from the pool to award the current date/blind - keep picking until we pick an application that hasn't been selected yet.
			var lotteryBall int
			for pickAgain := true; pickAgain; {
				lotteryBall = rand.Intn(len(datePool))
				if !selectedFromDatePool.contains(datePool[lotteryBall].CID) {
					// the lottery has picked a winner - mark this CID as having been awarded a blind
					selectedCIDs[datePool[lotteryBall].CID] = true
					selectedFromDatePool[datePool[lotteryBall].CID] = true
					pickAgain = false
				} else {
					log.Printf("CID [%v] was already selected; another lottery ball will be drawn for HUNT DATE [%v], BLIND [%v]\n", datePool[lotteryBall].CID, huntDate, blindNum)
				}
			}

			// store the results of this individual lottery award
			result := LotteryResultType{
				Application: datePool[lotteryBall],
				AwardedDate: huntDate,
				BlindNumber: blindNum,
			}

			lotteryResults = append(lotteryResults, &result)

			// indicate this CID has been awarded another blind by incrementing its awards counter
			incrementCIDAward(result.Application.CID)

			lotteryWinnerAwards := currentCIDAwards(result.Application.CID)
			lotteryWinnerCID := result.Application.CID
			log.Printf("WINNER! CID [%v] has been awarded HUNT DATE [%v], BLIND [%v] (total awards=[%v])\n", lotteryWinnerCID, huntDate, blindNum, lotteryWinnerAwards)

			// If this hunter has 2 or more awards compared to another hunter, strip one of his awards and give it to that other hunter.
			// Start at the opposite end of the histogram (give it to a hunter with the least amount of awards compared to everyone else).
			// But only give the hunt to the other hunter if that other hunter requested the hunt date that this hunter previously won.
		out:
			// loop through awardsByCID which is sorted ascending - thus the first list item has the least amount of blinds awarded to it
			for _, loser := range awardsByCID {
				// if this loser has 2 or more fewer lottery wins than the current winner, we will want to give the loser one of the current winner's assigned blinds
				if lotteryWinnerAwards-loser.Awards >= 2 {
					log.Printf("Winner [%v] of HUNT DATE [%v], BLIND [%v] has a surplus of awards [%v] compared to CID [%v] who only has [%v] awards. Will attempt to transfer one of the awards.\n", lotteryWinnerCID, huntDate, blindNum, lotteryWinnerAwards, loser.CID, loser.Awards)
					// find all the current winner's lottery wins - the first one that is of a hunt date in which the loser requested will be the one we re-assign
					for _, aLotteryResult := range lotteryResults {
						if aLotteryResult.Application.CID == lotteryWinnerCID {
							// if the current winner won a blind on a date the loser requested, give that to the loser, unless the loser already has one for that date
							if didCIDRequestDate(loser.CID, aLotteryResult.AwardedDate) && !wasCIDAwardedBlindOnDate(loser.CID, aLotteryResult.AwardedDate) {
								// strip the previous award from the current lottery winner and reassign it to the loser
								aLotteryResult.Application = applicationsByCID[loser.CID]
								decrementCIDAward(lotteryWinnerCID)
								incrementCIDAward(loser.CID)
								selectedCIDs[loser.CID] = true
								log.Printf("FAIRNESS AWARD: CID [%v] gave to CID [%v] a hunt on [%v] in blind [%v]\n", lotteryWinnerCID, loser.CID, aLotteryResult.AwardedDate, aLotteryResult.BlindNumber)
								break out
							}
						}
					}
				}
			}
		}
	}

	log.Println("==================== ALL BLINDS HAVE BEEN ASSIGNED ====================")
}

// publishResults writes the lottery results to the output CSV file.
func publishResults() {
	outputCSVRecords := make([][]string, 0)
	outputCSVRecords = append(outputCSVRecords, []string{"Date", "Blind", "CID", "Name"})
	for _, award := range lotteryResults {
		csvRow := make([]string, 4)
		csvRow[0] = award.AwardedDate
		csvRow[1] = fmt.Sprintf("%v", award.BlindNumber)
		csvRow[2] = string(award.Application.CID)
		csvRow[3] = award.Application.Name
		outputCSVRecords = append(outputCSVRecords, csvRow)
	}

	outputFile, err := os.Create(*argOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	w := csv.NewWriter(outputFile)
	w.WriteAll(outputCSVRecords)
	if err := w.Error(); err != nil {
		log.Fatalln("Error writing CSV file:", err)
	}
	log.Printf("Lottery results written to file [%v]\n", *argOutput)
}

// logStats prints a bunch of statistics to stdout.
func logStats() {
	histo := make(map[int]int, 0)
	cidAwardsMap := make(map[CID]int, len(lotteryResults))
	unassignedBlinds := 0
	losers := make(map[CID]string, len(applicationsList)) // value is the hunter name

	// seed the losers list with all names; names will be removed except for those who lost every lottery draw
	for _, a := range applicationsList {
		losers[a.CID] = a.Name
	}

	// loop through the actual lottery results - we want to see the numbers from the actual results list
	for _, l := range lotteryResults {
		if l.Application.CID != "" {
			cidAwardsMap[l.Application.CID]++
			delete(losers, l.Application.CID)
		} else {
			unassignedBlinds++
		}
	}
	// dump how many blinds were awarded to each CID and
	// a histogram showing how often N blinds were awarded
	log.Println("=====START AWARD COUNTS=====")
	fmt.Println("CID,Awards")
	for cid, awards := range cidAwardsMap {
		fmt.Printf("%v,%v\n", cid, awards)
		histo[awards]++
	}
	histo[0] = len(losers)
	log.Println("=====END AWARD COUNTS=====")

	log.Println("=====START CIDS WITHOUT AN AWARD=====")
	fmt.Println("CID,Name")
	for cid, name := range losers {
		fmt.Printf("%v,%v\n", cid, name)
	}
	log.Println("=====END CIDS WITHOUT AN AWARD=====")

	log.Println("=====START HISTOGRAM=====")
	fmt.Println("Awards,Count")
	for key, value := range histo {
		fmt.Printf("%v,%v\n", key, value)
	}
	log.Println("=====END HISTOGRAM=====")

	log.Printf("Unassigned Blinds=[%v]\n", unassignedBlinds)
	log.Printf("Number of Hunt Dates=[%v]\n", len(huntDates))
	log.Printf("Number of Blinds=[%v]\n", *argNumBlinds)
	log.Printf("Total Hunts Available in Lottery=[%v]\n", len(huntDates)*(*argNumBlinds))
	log.Printf("Total Hunts Assigned/Unassigned in Lottery=[%v]\n", len(lotteryResults))
	log.Printf("Unassigned Blinds=[%v]\n", unassignedBlinds)
	log.Printf("Number of Applications=[%v]\n", len(applicationsList))
	log.Printf("Number of CIDs awarded blinds=[%v]\n", len(cidAwardsMap))
	log.Printf("Number of CIDs not awarded any blinds=[%v]\n", len(losers))
}

// validateLotteryResults will confirm the results make sense - if there is a bug or misbehavior in the lottery algorithm, this hopefully will catch the problem.
func validateLotteryResults() {
	for _, award := range lotteryResults {
		if award.Application.CID != "" && !didCIDRequestDate(award.Application.CID, award.AwardedDate) {
			log.Fatalf("A CID was awarded a hunt on a date that was not requested: [%+v]", award)
		}
	}

	if len(lotteryResults) != (len(huntDates) * *argNumBlinds) {
		log.Fatalf("Number of lottery results did not match the total number of hunts available")
	}

	for _, hd := range huntDates {
		blindNumForDate := make(map[int]bool, *argNumBlinds) // key is blind number
		cidForDate := make(map[CID]bool, *argNumBlinds)
		for _, award := range lotteryResults {
			if award.AwardedDate == hd && award.Application.CID != "" {
				if _, ok := blindNumForDate[award.BlindNumber-1]; ok {
					log.Fatalf("Hunt date [%v] awarded a blind [%v] multiple times", hd, award.BlindNumber)
				}
				blindNumForDate[award.BlindNumber-1] = true
				if _, ok := cidForDate[award.Application.CID]; ok {
					log.Fatalf("Hunt date [%v] awarded a blind to the same CID [%v] multiple times", hd, award.Application.CID)
				}
				cidForDate[award.Application.CID] = true
			}
		}
	}

	log.Println("Lottery results have been validated")
}

func init() {
	// log everything to stderr so that it can be easily gathered by logs, separate log files are problematic with containers
	log.SetOutput(os.Stderr)
}

// main is the program entry point
func main() {
	processArguments()
	readApplicationData()
	performLotteryDraw()
	logStats()
	publishResults()
	validateLotteryResults()
}
