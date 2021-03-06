// Wrapper functions demonstrating how some gpx functions are used.
// 01   July  5, 2018   Initial version uploaded to github
// 02   Aug. 28, 2018   Simplified to reduce complexity and remove functionality

package main

import (
	"bufio"
	"fmt"
	"github.com/go-opt/gpx"
	"github.com/pkg/errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// Variables controlling program input and output. The full absolute path for the
// data files must be provided if the files reside in a directory other than where
// the executable was launched.
var pauseAfter       int = 50       // Number of items to print before pausing
var sampleLpFile  string = "inputMpsLp1.txt"  // MPS file for LP example (afiro)
var sampleMipFile string = "inputGpxMip1.txt" // Text file for MIP example (noswot)
var fileNameSoln  string = "soln_file.txt"    // Solution file generated by Cplex
var fileNameMps   string = "mps_file.txt"     // MPS file of the model generated by Cplex

// Need to make gpx variables global to this package to make them available to all
// wrapper functions that need them without having to pass them as arguments.
var gName     string            // gpx input problem name
var gRows   []gpx.InputRow      // gpx input rows
var gCols   []gpx.InputCol      // gpx input cols
var gElem   []gpx.InputElem     // gpx input elems
var gObj    []gpx.InputObjCoef  // gpx input objective function coefficients
var sObjVal   float64           // Solution value of objective function
var sRows   []gpx.SolnRow       // Solution rows provided by gpx
var sCols   []gpx.SolnCol       // Solution columns provided by gpx

//==============================================================================

// printOptions displays the options that are available for testing. Package
// global flags control which menus are printed.
// The function accepts no arguments and returns no values.
func printOptions() {

	fmt.Printf("\nAvailable Options:\n\n")
	fmt.Printf(" 0 - EXIT program\n")
	fmt.Printf(" 1 - solve sample LP problem (afiro) from MPS data file\n")
	fmt.Printf(" 2 - solve sample MILP problem (noswot) from data structures\n")
	fmt.Printf(" 3 - display solution\n")

}

//==============================================================================

// wpIsMip checks if the problem is considered a MIP according to Cplex. It scans
// the columns, and if it detects any type other than "C" ("continuous" in Cplex),
// it is considered a MIP, and the function returns "true". If it does not find one, 
// it returns false and the problem can be solved as a pure linear problem (LP).
func wpIsMip() bool {
	
	for i := 0; i < len(gCols); i++ {
		if gCols[i].Type != "C" {
			return true
		}
	}
	
	return false
}

//==============================================================================

// wpReadGpxFile reads a text file written in the special gpx format and populates
// the input gpx data structures, which are passed in as arguments. 
// In case of failure, the function returns an error.
func wpReadGpxFile(rows *[]gpx.InputRow, cols *[]gpx.InputCol, elem *[]gpx.InputElem,
					obj *[]gpx.InputObjCoef, probName *string) error {

	var fileName  string           // name of file from which gpx data are read
	var lineNum   int = 0          // line number being processed
	var numTokens int              // number of tokens in current line
	var readState int              // data block currently being processed
	var bigInt    int64            // placeholder when reading numbers from file
	var eof       bool = false     // flag indicating if end of file reached
	var rowItem   gpx.InputRow     // row item used in constructing gRows list
	var colItem   gpx.InputCol     // col item used in constructing gCols list
	var elemItem  gpx.InputElem    // elem item used in constructing gElem list
	var objItem   gpx.InputObjCoef // item used in constructing coefficients in obj. func.


	// Initialize lists that will be returned.
	*rows = nil
	*cols = nil
	*elem = nil
	*obj  = nil
	readState = -1

	fileName = sampleMipFile
	
	// Check that the input file exists and open it for reading.
	inputFile, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "Open MPS file failed")
	}

	defer inputFile.Close()
	fileReader := bufio.NewReader(inputFile)


	// Create the token which will be used when reading the file.
	token := make([]string, 1)

	// The main file reading loop ----------------------------------------------

	for !eof {
		lineNum++

		curLine, err := fileReader.ReadString('\n') //reads in a line

		// EOF is not an error which would cause the function to abort, so ignore.

		if err == io.EOF {
			err = nil
			eof = true
			if readState != 5 {
				fmt.Printf("WARNING: End of data token missing, %d lines read.\n", lineNum)
				return nil
			}
		}
		if err != nil {
			return errors.Errorf("Problem reading line %d", lineNum)
		}

		if string(curLine[0]) == "#" {
			continue
		} // Skip lines with an asterisk in the first column
		
		// Split the line above into a slice of tokens using strings.Fields.
		token = strings.Fields(curLine) 
		numTokens = len(token)
		if numTokens == 0 {
			continue
		} //skip blank lines

		// Take the appropriate action for a new keyword ---------------------------

		switch strings.ToUpper(token[0]) {

		case "PROBLEM_NAME:":
			if numTokens == 1 {
				*probName = "NoName"
			} else {
				*probName = token[1]
			}
			readState = 0
			continue

		case "OBJECTIVE_START":
			readState = 1
			continue

		case "ROWS_START":
			readState = 2
			continue

		case "COLUMNS_START":
			readState = 3
			continue

		case "ELEMENTS_START":
			readState = 4
			continue

		case "END_DATA":
			readState = 5

		default:
			if readState < 0 {
				return errors.Errorf("Unexpected format at line %d", lineNum)				
			}

		} // end of switch on data block

		if readState == 5 {
			break
		}		

		switch readState {

		case 1:  // Objective function
			if numTokens != 2 {
				return errors.Errorf("Invalid OBJ items on line %d", lineNum)
			} 
			bigInt, _           = strconv.ParseInt(token[0], 10, 64)
			objItem.ColIndex    = int(bigInt)      
			objItem.Value,    _ = strconv.ParseFloat(token[1], 64)
			*obj = append(*obj, objItem)
		
		case 2: // Reading rows
			if numTokens != 4 {
				return errors.Errorf("Invalid ROW items on line %d", lineNum)
			} 
			rowItem.Name      = token[0]
			rowItem.Sense     = token[1]
			rowItem.Rhs,    _ = strconv.ParseFloat(token[2], 64)
			rowItem.RngVal, _ = strconv.ParseFloat(token[3], 64)
			*rows = append(*rows, rowItem)
		
		case 3: // Reading columns
			if numTokens != 4 {
				return errors.Errorf("Invalid COL items on line %d", lineNum)
			} 
			colItem.Name     = token[0]
			colItem.Type     = token[1]
			colItem.BndLo, _ = strconv.ParseFloat(token[2], 64)
			colItem.BndUp, _ = strconv.ParseFloat(token[3], 64)
			*cols = append(*cols, colItem)
		
		case 4: // Reading elements
			if numTokens != 3 {
				return errors.Errorf("Invalid ELEM items on line %d", lineNum)
			} 
			bigInt, _            = strconv.ParseInt(token[0], 10, 64)
			elemItem.RowIndex    = int(bigInt)
			bigInt, _            = strconv.ParseInt(token[1], 10, 64)
			elemItem.ColIndex    = int(bigInt)
			elemItem.Value,    _ = strconv.ParseFloat(token[2], 64)
			*elem = append(*elem, elemItem)
								
		} // end switch on readState			
	} // end of loop reading file	

	return nil	
}

//==============================================================================

// wpInitGpx initializes all global input and solution variables. It accepts
// no input and returns no values.
func wpInitGpx() {

	// Initialize all global gpx data structures.	
	gName   = ""
	gRows   = nil
	gCols   = nil
	gElem   = nil
	sObjVal = 0.0
	sRows   = nil
	sCols   = nil
	
}

//==============================================================================

// wpPrintGpxSoln prints the gpx solution data structures. It accepts no arguments
// and returns no values.
func wpPrintGpxSoln() {
	var userString string  // user input
	var counter    int     // counter keeping track of number of lines printed
	
	fmt.Printf("\nObjective function value = %f\n\n", sObjVal)
	
	userString = ""
	fmt.Printf("Display additional results [Y|N]: ")
	fmt.Scanln(&userString)

	if userString == "y" || userString == "Y" {

		userString = ""
		fmt.Printf("\nDisplay variable list [Y|N]: ")
		fmt.Scanln(&userString)
		if userString == "y" || userString == "Y" {
			if len(sCols) != 0 {
				counter = 0
				for i := 0; i < len(sCols); i++ {
					fmt.Printf("Col %4d: %15s, Val = %13e,  Reduced cost = %13e\n", 
								i, sCols[i].Name, sCols[i].Value, sCols[i].RedCost)
					counter++
					userString = ""
					if counter == pauseAfter {
						fmt.Printf("\nPAUSED... <CR> continue, any key to quit: ")
						fmt.Scanln(&userString)
						if userString != "" {
							break 
						}		
					} // end if pause needed
				} // end for printing variables
			} else {
				fmt.Printf("List of solved variables is empty.\n")
			} // end else varibable list is empty			
		} // end if displaying variables

		userString = ""
		fmt.Printf("\nDisplay constraint list [Y|N]: ")
		fmt.Scanln(&userString)
		if userString == "y" || userString == "Y" {			
			if len(sRows) != 0 {
				counter = 0
				for i := 0; i < len(sRows); i++ {
					fmt.Printf("Row %4d: %15s, Pi = %13e,  Slack = %13e\n", 
								i, sRows[i].Name, sRows[i].Pi, sRows[i].Slack)
					counter++
					userString = ""
					if counter == pauseAfter {
						fmt.Printf("\nPAUSED... <CR> continue, any key to quit: ")
						fmt.Scanln(&userString)
						if userString != "" {
							break 
						}			
					} // end if pause needed
				} // end for printing constraints
			} else {
				fmt.Printf("List of solved constraints is empty.\n")
			} // end else constraints list is empty						
		} // end if displaying constraints
	} // end if printing results
		
}

//==============================================================================

// wpSolveFromStruct illustrates an example of a problem solved using the internal
// data structures. It reads data from file, populates the internal data structures,
// solves the problem, prints the solution, and gives user the option to save
// the model and solution to file. Function accepts no arguments.
// In case of failure, function returns an error.
func wpSolveFromStruct() error {
	var fileType     string   // file type as recognized by Cplex
	var dispToScreen bool     // flag indicating if Cplex output should be displayed
	var err          error    // error returned from functions called

	fmt.Printf("\nThis example illustrates how to solve a problem by using internal\n")
	fmt.Printf("gpx data structures defining the model.\n\n")

	// Initialize data structures and variables
	wpInitGpx()
	dispToScreen = true	
	fileType     = "MPS"
	
	fmt.Printf("Populating data        - translating file '%s' to data structures...\n", sampleMipFile)
	if err = wpReadGpxFile(&gRows, &gCols, &gElem, &gObj, &gName); err != nil {
		return errors.Wrap(err, "wpReadGpxFile failed")		
	}

	fmt.Printf("Running CreateProb     - initialize environment for problem '%s'...\n", gName)
	if err = gpx.CreateProb(gName); err != nil {
		return errors.Wrap(err, "Failed to initialize environment")				
	}
	
	fmt.Printf("Running OutputToScreen - set echo to '%t'...\n", dispToScreen)
	if err = gpx.OutputToScreen(dispToScreen); err != nil {
		return errors.Wrap(err, "Failed to set display to screen")						
	}
	
	fmt.Printf("Running NewRows        - create rows in Cplex...\n")
	if err = gpx.NewRows(gRows); err != nil {
		return errors.Wrap(err, "Failed to create new rows")								
	}
	
	fmt.Printf("Running NewCols        - create new columns in Cplex...\n")
	if err = gpx.NewCols(gObj, gCols); err != nil {
		return errors.Wrap(err, "Failed to create new columns")										
	}
	
	fmt.Printf("Running ChgCoefList    - create non-zero coefficients in Cplex...\n")
	if err = gpx.ChgCoefList(gElem); err != nil {
		return errors.Wrap(err, "Failed to create new columns")		
	}
	
	if dispToScreen {
		// Add a blank line between our output and Cplex output.
		fmt.Println("")
	}

	if wpIsMip() {
		fmt.Printf("Running MipOpt         - solve the MIP...\n")
		if err = gpx.MipOpt(); err != nil {
			return errors.Wrap(err, "Failed to solve the MIP")			
		}		
		// To make things pretty, separate our output from Cplex output by blank line.
		if dispToScreen {
			fmt.Printf("\n")
		}
	
		fmt.Printf("Running GetMipSolution - get MIP solution from Cplex...\n")
		err = gpx.GetMipSolution(&sObjVal, &sRows, &sCols)
		if err != nil {
			return errors.Wrap(err, "Failed to get MIP solution")
		} 
	} else {
		fmt.Printf("Running LpOpt          - solve the LP...\n")
		if err = gpx.LpOpt(); err != nil {
			return errors.Wrap(err, "Failed to solve the LP")
		}		

		// To make things pretty, separate our output from Cplex output by blank line.
		if dispToScreen {
			fmt.Printf("\n")
		}
	
		fmt.Printf("Running GetSolution    - get LP solution from Cplex...\n")
		err = gpx.GetSolution(&sObjVal, &sRows, &sCols)
		if err != nil {
			return errors.Wrap(err, "Failed to get LP solution")
		} 
	} // end else problem is LP


	fmt.Printf("Running SolWrite       - write solution to file '%s'...\n", fileNameSoln)
	if err = gpx.SolWrite(fileNameSoln); err != nil {
		return errors.Wrap(err, "Failed to write solution file")	
	}

	fmt.Printf("Running WriteProb      - write model to MPS file '%s'...\n", fileNameMps)
	if err = gpx.WriteProb(fileNameMps, fileType); err != nil {
		return errors.Wrap(err, "Failed to write model file")	
	}	

	fmt.Printf("Running CloseCplex     - clean up the environment...\n")
	if err = gpx.CloseCplex(); err != nil {
		return errors.Wrap(err, "Failed to close Cplex")
	}

	// Print the solution using a separate function	
	wpPrintGpxSoln()

	return nil
}

//==============================================================================

// wpSolveFromFile illustrates an example of a problem solved by reading a data
// file directly by Cplex. After reading the file, the function solves the problem, 
// prints the solution, and gives user the option to save the model and solution to 
// file. Function accepts no arguments.
// In case of failure, function returns an error.
func wpSolveFromFile() error {
	var fileNameIn   string   // name of the input file containing the model
	var fileType     string   // type of input file
	var dispToScreen bool     // flag indicating if Cplex output should be displayed
	var isMip        bool     // flag differentiating between MIP and LP problems
	var err          error    // error returned from functions called

	fmt.Printf("\nThis example illustrates how to solve an LP problem by reading data\n")
	fmt.Printf("containing the model from an MPS data file.\n\n")

	// Initialize all variables. In a previous incarnation of this executable,
	// this information was provided by the user. Now it is hard-coded.
	wpInitGpx()	
	gName        = "SampleLP01"
	isMip        = false
	dispToScreen = true
	fileType     = "MPS"
	fileNameIn   = sampleLpFile
	
	fmt.Printf("Running CreateProb     - initialize environment for problem '%s'...\n", gName)
	if err = gpx.CreateProb(gName); err != nil {
		return errors.Wrap(err, "Failed to initialize environment")				
	}
	
	fmt.Printf("Running OutputToScreen - set echo to '%t'...\n", dispToScreen)
	if err = gpx.OutputToScreen(dispToScreen); err != nil {
		return errors.Wrap(err, "Failed to set display to screen")						
	}

	fmt.Printf("Running ReadCopyProb   - read %s data file %s...\n\n", fileType, fileNameIn)
	if err = gpx.ReadCopyProb(fileNameIn, fileType); err != nil {
		return errors.Wrap(err, "Open MPS file failed")
	} 

	if dispToScreen {
		// Add a blank line between our output and Cplex output.
		fmt.Printf("\n\n")
	}

	if isMip {
		fmt.Printf("Running MipOpt         - solve the MIP...\n\n")
		if err = gpx.MipOpt(); err != nil {
			return errors.Wrap(err, "Failed to solve the MIP")			
		}		
		// To make things pretty, separate our output from Cplex output by blank line.
		if dispToScreen {
			fmt.Printf("\n")
		}
	
		fmt.Printf("Running GetMipSolution - get MIP solution from Cplex...\n")
		err = gpx.GetMipSolution(&sObjVal, &sRows, &sCols)
		if err != nil {
			return errors.Wrap(err, "Failed to get MIP solution")
		} 
	} else {
		fmt.Printf("Running LpOpt          - solve the LP...\n\n")
		if err = gpx.LpOpt(); err != nil {
			return errors.Wrap(err, "Failed to solve the LP")
		}		

		// To make things pretty, separate our output from Cplex output by blank line.
		if dispToScreen {
			fmt.Printf("\n")
		}
	
		fmt.Printf("Running GetSolution    - get LP solution from Cplex...\n")
		err = gpx.GetSolution(&sObjVal, &sRows, &sCols)
		if err != nil {
			return errors.Wrap(err, "Failed to get LP solution")
		} 
	} // end else problem is LP

	fmt.Printf("Running SolWrite       - write solution to file '%s'...\n", fileNameSoln)
	if err = gpx.SolWrite(fileNameSoln); err != nil {
		return errors.Wrap(err, "Failed to write solution file")	
	}	

	fmt.Printf("Running CloseCplex     - clean up the environment...\n")
	if err = gpx.CloseCplex(); err != nil {
		return errors.Wrap(err, "Failed to close Cplex")
	}

	// Print the solution using a separate function	
	wpPrintGpxSoln()
		
	return nil
}

//==============================================================================

// runMainWrapper displays the menu of options available, prompts the user to enter
// one of the options, and executes the command specified. 
// The function accepts no arguments and returns no values.
func runMainWrapper() {
	var cmdOption     string  // command option
	var err            error  // error returned by called functions

	// Print header and options, and enter infinite loop until user quits.

	fmt.Printf("\nDEMONSTRATION OF GPX FUNCTIONALITY\n")
	
	for {

		// Initialize variables, read command, and execute command.		
		printOptions()
		cmdOption    = ""		
		fmt.Printf("\nEnter a new option: ")
		fmt.Scanln(&cmdOption)

		switch cmdOption {

		case "0":
			fmt.Println("\n===> NORMAL PROGRAM TERMINATION <===\n")
			return

		case "1":
			// Solve problem from data file
			err = wpSolveFromFile()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("\nExample solving problem from data file completed.\n")
			}

		case "2":
			// Solve problem from internal structures
			err = wpSolveFromStruct()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("\nExample solving problem from data structures completed.\n")
			}

		case "3":
			// Print gpx solution			
			wpPrintGpxSoln()
			fmt.Printf("\nDisplay of solution completed.\n")
			
			
		default:
			fmt.Printf("Unsupported option: '%s'\n", cmdOption)
						
		} // end of switch on cmdOption
	} // end for looping over commands

}

//==============================================================================

// main function calls the main wrapper. It accepts no arguments and returns
// no values.
func main() {

	runMainWrapper()
}

//============================ END OF FILE =====================================