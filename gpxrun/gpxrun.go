// This file contains the main function and needed only by this executable.
// Additional functions, which are shared with the lpo executable (functions copied,
// not exported) are included in a separate file.

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

// Flags to control display of menus and use of customized environment.
var mainMenuOn bool = true    // Flag for main LPO functions
var gpxMenuOn  bool = false   // Flag for enabling gpx functions   
var custEnvOn  bool = false   // Flag to enable custom paths and names
var pauseAfter int = 50       // Number of items to print before pausing

// Customized environment used if custEnvOn = true.
// It is intended to reduce the amount of typing for SOME (not all) user input,
// and to build names of internal files related to the "base" name specified.
// If disabled, user must enter complete directory and file names when prompted.

var dSrcDev       string = "C:/Users/Admin/go/src/gpx/gpxrun/" // Development source data dir
var fExtension    string = ".txt"   // Extension of source data files in development dir.  


//==============================================================================

// printOptions displays the options that are available for testing. Package
// global flags control which menus are printed.
// The function accepts no arguments and returns no values.
func printOptions() {

	fmt.Println("\nAvailable Options (0 to EXIT):")
	fmt.Println("")
	fmt.Println(" g - toggle gpx function exerciser           c - toggle custom environment")
	
  if mainMenuOn {
	fmt.Println("")
	fmt.Println(" 1 - solve problem from data structures      2 - solve problem from data file")
	fmt.Println(" 3 - initialize internal data structures     4 - populate input data structures")
	fmt.Println(" 5 - display input data                      6 - display solution")
  }

  if gpxMenuOn {
	fmt.Println("")
	fmt.Println("61 - ChgCoefList      62 - ChgObjSen        63 - ChgProbName      64 - CloseCplex")
	fmt.Println("65 - CreateProb       66 - GetColName       67 - GetMipSolution   68 - GetNumCols")
	fmt.Println("69 - GetNumRows       70 - GetObjVal        71 - GetRowName       72 - GetSlack")
	fmt.Println("73 - GetSolution      74 - GetX             75 - LpOpt            76 - MipOpt")
	fmt.Println("77 - NewCols          78 - NewRows          79 - OutputToScreen   80 - ReadCopyProb")
	fmt.Println("81 - SolWrite         82 - WriteProb")
  }

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

	fmt.Printf("Enter name of GPX file to be read: ")
	fmt.Scanln(&fileName)
	if custEnvOn {
		fileName = dSrcDev + fileName + fExtension
	}



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

// wpSolveFromStruct illustrates an example of a problem solved using the internal
// data structures. It reads data from file, populates the internal data structures,
// solves the problem, prints the solution, and gives user the option to save
// the model and solution to file. Function accepts no arguments.
// In case of failure, function returns an error.
func wpSolveFromStruct() error {
	var err error
	var userString string
	var dispToScreen bool

	fmt.Printf("\nThis example illustrates how to solve a problem by using internal\n")
	fmt.Printf("gpx data structures defining the model.\n\n")
	fmt.Printf("Here, the data structures are populated by reading a text file written\n")
	fmt.Printf("in the gpx format. Under 'normal' conditions, users would use their own\n")
	fmt.Printf("functions to populate the needed structures, or use the lpo package\n")
	fmt.Printf("to read an MPS file or populate the lpo data structures and translate\n")
	fmt.Printf("them to gpx using the lpo.TransToGpx function.\n\n")

	// Initialize data structures and variables
	wpInitGpx()
	dispToScreen = false	

	
	if err = wpReadGpxFile(&gRows, &gCols, &gElem, &gObj, &gName); err != nil {
		return errors.Wrap(err, "wpReadGpxFile failed")		
	}

	// Set the display to screen parameter
	userString = ""
	fmt.Printf("\nShould Cplex display output on the screen [Y|N]: ")
	fmt.Scanln(&userString)
	if userString == "y" || userString == "Y" {
		fmt.Printf("Cplex output will be displayed on the screen.\n")
		dispToScreen = true	
	} else {
		fmt.Printf("Cplex output will NOT be displayed on the screen.\n")		
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


	// Print the solution using a separate function	
	wpPrintGpxSoln()

	userString = ""
	fmt.Printf("\nDo you wish to save the solution to a file [Y|N]: ")
	fmt.Scanln(&userString)
	if userString == "y" || userString == "Y" {
		userString = ""
		fmt.Printf("Enter name of solution file: ")
		fmt.Scanln(&userString)
		if custEnvOn {
			userString = dSrcDev + userString + fExtension
		}

		fmt.Printf("Running SolWrite       - write solution to file '%s'...\n", userString)
		if err = gpx.SolWrite(userString); err != nil {
			return errors.Wrap(err, "Failed to write solution file")	
		}
		fmt.Printf("\n")	
	} // end if saving results to file

	userString = ""
	fmt.Printf("Do you wish to save the model to a file [Y|N]: ")
	fmt.Scanln(&userString)
	if userString == "y" || userString == "Y" {
		fmt.Printf("Running ReadCopyProb   - writing problem to file.\n")
		wpWriteProb()
		fmt.Printf("\n")	
	} // end if saving model to file


	fmt.Printf("\nRunning CloseCplex     - clean up the environment...\n")
	if err = gpx.CloseCplex(); err != nil {
		return errors.Wrap(err, "Failed to close Cplex")
	}

	fmt.Printf("\nThis example is now concluded. The functions demonstrated, in addition\n")
	fmt.Printf("to others that were not used, can be accessed individually by toggling\n")
	fmt.Printf("the gpx functions exerciser ('g') of this tutorial.\n\n")
	fmt.Printf("Users are encouraged to use the tutorial provided in the lpo package\n")
	fmt.Printf("which provides access to the same gpx functions, but also allows\n")
	fmt.Printf("the model to be populated by reading files in MPS format (in lpo)\n")
	fmt.Printf("and then translating them to gpx format (via lpo.TransToGpx).\n\n")
		
	return nil
}

//==============================================================================

// wpSolveFromFile illustrates an example of a problem solved by reading a data
// file directly by Cplex. After reading the file, the function solves the problem, 
// prints the solution, and gives user the option to save the model and solution to 
// file. Function accepts no arguments.
// In case of failure, function returns an error.
func wpSolveFromFile() error {
	var err error
	var userString string
	var dispToScreen bool
	var isMip bool

	fmt.Printf("\nThis example illustrates how to solve a problem by reading data\n")
	fmt.Printf("containing the model from a file.\n\n")

	// Initialize all variables.
	wpInitGpx()
	
	userString = ""
	isMip      = false

	// Get problem type (LP or MIP) and problem name from user.

	fmt.Printf("Enter the name of the problem: ")
	fmt.Scanln(&gName)

	fmt.Printf("\nIs the problem a mixed integer problem (MIP) [Y|N]: ")
	fmt.Scanln(&userString)
	if userString == "y" || userString == "Y" {
		isMip = true
	}


	// Set the display to screen parameter
	dispToScreen = false	
	userString = ""
	fmt.Printf("\nShould Cplex display output on the screen [Y|N]: ")
	fmt.Scanln(&userString)
	if userString == "y" || userString == "Y" {
		fmt.Printf("\nCplex output will be displayed on the screen.\n")
		dispToScreen = true	
	} else {
		fmt.Printf("\nCplex output will NOT be displayed on the screen.\n")		
	}

	fmt.Printf("Running CreateProb     - initialize environment for problem '%s'...\n", gName)
	if err = gpx.CreateProb(gName); err != nil {
		return errors.Wrap(err, "Failed to initialize environment")				
	}
	
	fmt.Printf("Running OutputToScreen - set echo to '%t'...\n", dispToScreen)
	if err = gpx.OutputToScreen(dispToScreen); err != nil {
		return errors.Wrap(err, "Failed to set display to screen")						
	}

	fmt.Printf("Running ReadCopyProb   - to read data from file...\n\n")
	if err = wpReadDataFile(); err != nil {
		return errors.Wrap(err, "Failed to read data from file")		
	}

	if dispToScreen {
		// Add a blank line between our output and Cplex output.
		fmt.Printf("\n\n")
	}

	if isMip {
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


	// Print the solution using a separate function	
	wpPrintGpxSoln()

	userString = ""
	fmt.Printf("\nDo you wish to save the solution to a file [Y|N]: ")
	fmt.Scanln(&userString)
	if userString == "y" || userString == "Y" {
		userString = ""
		fmt.Printf("Enter name of solution file: ")
		fmt.Scanln(&userString)
		if custEnvOn {
			userString = dSrcDev + userString + fExtension
		}
		fmt.Printf("Running SolWrite       - write solution to file '%s'...\n", userString)
		if err = gpx.SolWrite(userString); err != nil {
			return errors.Wrap(err, "Failed to write solution file")	
		}	
		fmt.Printf("\n")	
	} // end if saving results to file

	userString = ""
	fmt.Printf("Do you wish to save the model to a file [Y|N]: ")
	fmt.Scanln(&userString)
	if userString == "y" || userString == "Y" {
		fmt.Printf("Running ReadCopyProb   - writing problem to file.\n")
		wpWriteProb()
		fmt.Printf("\n")	
	} // end if saving model to file


	fmt.Printf("\nRunning CloseCplex     - clean up the environment...\n")
	if err = gpx.CloseCplex(); err != nil {
		return errors.Wrap(err, "Failed to close Cplex")
	}

	fmt.Printf("\nThis example is now concluded. The functions demonstrated, in addition\n")
	fmt.Printf("to others that were not used, can be accessed individually by toggling\n")
	fmt.Printf("the gpx functions exerciser ('g') of this tutorial.\n\n")
		
	return nil
}

//==============================================================================

// runMainWrapper displays the menu of options available, prompts the user to enter
// one of the options, and executes the command specified. The main wrapper controls
// the main commands, and in turn calls secondary wrappers to execute additional
// commands. The flags which control the display of menu options have no impact on
// the available commands. All commands are available even if the corresponding menu
// item is "hidden". The function accepts no arguments and returns no values.
func runMainWrapper() {

	var cmdOption     string  // command option
	var err            error  // error returned by called functions


	// Print header and options, and enter infinite loop until user quits.

	fmt.Println("\nTUTORIAL AND EXERCISER FOR GPX FUNCTIONS.")
	printOptions()
	
	for {

		// Initialize variables, read command, and execute command.
		
		cmdOption    = ""		
		fmt.Printf("\nEnter a new option: ")
		fmt.Scanln(&cmdOption)

		switch cmdOption {

		//---------------- Commands for toggles --------------------------------

		case "g":
			if gpxMenuOn {
				gpxMenuOn = false
				fmt.Println("\nFunctions menu commands will be disabled.")
			} else {
				gpxMenuOn = true
				fmt.Println("\nFunctions menu commands will be enabled.")
				printOptions()				
			}
						
		case "c":
			if custEnvOn {
				fmt.Printf("\nCustomized environment disabled.\n")
				fmt.Printf("Full file paths must be entered when needed.\n")
				custEnvOn = false
			} else {
				fmt.Printf("\nWARNING: Customized environment enabled.\n\n")
				fmt.Printf("When prompted for file names, only the base name (without path\n")
				fmt.Printf("to directory or file extension) needs to be entered.\n\n")
				fmt.Printf("Directory for files = '%s'\n", dSrcDev)
				fmt.Printf("File extension      = '%s'\n", fExtension)
				custEnvOn = true				
			}

		//------------- Functions exercised by main wrapper --------------------

		case "0":
			fmt.Println("\n===> NORMAL PROGRAM TERMINATION <===\n")
			return


		case "1":
			// Solve problem from internal structures
			err = wpSolveFromStruct()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Example solving problem from data structures completed.\n")
			}

		case "2":
			// Solve problem from data file
			err = wpSolveFromFile()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Example solving problem from data file completed.\n")
			}
			
		case "3":
			wpInitGpx()
			fmt.Printf("All data structures have been initialized.\n")
					
		case "4":			
			err = wpReadGpxFile(&gRows, &gCols, &gElem, &gObj, &gName)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Data structures successfully populated from file.\n")
			}

		case "5":			
			wpPrintGpxIn()
			fmt.Printf("\nDisplay of input data structures completed.\n")

		case "6":			
			wpPrintGpxSoln()
			fmt.Printf("\nDisplay of solution completed.\n")
			
			
		default:

			// If the command was not present in this wrapper, check the other ones.
			// Only if the command cannot be satisfied by any of the secondary
			// wrappers treat this as an "error" and display the available commands.
						

			if gpxMenuOn {
				if err = runGpxWrapper(cmdOption); err == nil {
					// Found the command in gpx menu, continue
					continue
				}
			}

			fmt.Printf("Unsupported option: '%s'\n", cmdOption)
			printOptions()
						
		} // end of switch on cmdOption
	} // end for looping over commands

}

//==============================================================================

// main function calls the main wrapper. It accepts no arguments and returns
// no values.
func main() {
	
	runMainWrapper()
}
