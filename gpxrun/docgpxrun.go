/* 

Executable provides examples of gpx use and wrapper for exported functions.

SUMMARY


This executable provides examples of how the gpx package can be used to solve 
linear programming (LP) and mixed integer linear programming (MIP) problems
via Cplex and the callable C routines.
It also provides a wrapper to individually call and test each function exported by gpx.

The user must select one of the provided options to perform the desired task.
The options are grouped into a set of commands which illustrate main gpx functionality

	0 - exit program
	1 - solve problem from data structures
	2 - solve problem from data file
	3 - initialize internal data structures
	4 - populate input data structures
	5 - display input data
	6 - display solution

	
and a set of toggles which controls program behaviour as follows:

	c - toggle customized user environment
	g - toggle gpx function wrapper

To select an option, enter the corresponding letter or number when prompted.


MAIN COMMANDS

The main command options are displayed at all times and are as follows. To redisplay
the available options, enter a blank line or any other "unsupported" option.

Exit 

This option is used to terminate execution of the program. This option is displayed
as part of the command prompt and is not included in the lists showing other options.

Solve problem from data structures 

This option is used to populate the internal data structures by reading the model
definition from a text file prepared in a special format understood by gpx, and
solving the problem. Typically, gpx would be used in conjunction with the lpo
package which provides functions for reading MPS data files and populating the
needed data structures, or users would write their own functions to populate gpx.

This tutorial is intended to operate independently of lpo or other more sophisticated
methods of populating input data structures, and the simple file-reading mechanism is
used. The data structures are populated by reading the file, interpretting the
data read based on specific key words, parsing each line of input based on that
key word, and populating the lists in the order in which they are read. Manual
creation of the input text files is not recommended, but lpo may be used to read,
for example, an MPS file, and generating the corresponding gpx input data file.

The sequence of operations and gpx functions exercised with this option is as follows:
	
	wpReadGpxFile    - private function to read file and populate data structures
	CreateProb       - initialize Cplex environment and create the problem
	OutputToScreen   - set Cplex output to be displayed to screen or remain hidden
	NewRows          - create new rows
	NewCols          - create new columns
	ChgCoefList      - set non-zero coefficients for rows and columns
	
	if the problem is an LP
	LpOpt            - have Cplex solve the LP
	GetSolution      - populate the data structures with the LP solution
	
	or if the problem is a MIP
	MipOpt           - have Cplex solve the MIP
	GetMipSolution   - populate the data structures with the MIP solution
	
	wpPrintSoln      - private function to display objective function value and,
	                   optionally, variable values, reduced cost, pi, and slack
	SolWrite         - (optional) save the Cplex solution in a file
	ReadCopyProb     - (optional) save the model in a file
	CloseCplex       - clean up and close the Cplex environment
	
The data files, which can be used when running this example, and which are included
with the executable are as follows:

	inputGpxLp1.txt  - gpx format, small LP obtained from Cplex tutorial
	inputGpxLp2.txt  - gpx format, AFIRO LP obtained from netlib
	inputGpxMip1.txt - gpx format, small MIP obtained from Cplex tutorial
	inputGpxMip2.txt - gpx format, NOSWOT MIP obtained from miplib
	inputMpsLp1.txt  - MPS format, small LP obtained from the Cplex tutorial
	inputMpsLp2.txt  - MPS format, AFIRO LP obtained from netlib
	inputMpsMip1.txt - MPS format, small MIP obtained from the Cplex tutorial
	inputMpsMip2.txt - MPS format, NOSWOT MIP obtained from miplib

Solve problem from data file 

This option is used to have Cplex directly read a data file which defines the
model. Since the program does not know whether the source data file defines an
LP or a MIP, and does not know the problem name when the Cplex environment is
created in order to read the file, the user must provide this information.
Otherwise, the behaviour is the same as for the other option which uses internal 
data structures for model input.

The sequence of operations and gpx functions exercised with this option is as follows:
	
	user input       - get problem name and if problem is MIP or LP
	CreateProb       - initialize Cplex environment and create the problem
	OutputToScreen   - set Cplex output to be displayed to screen or remain hidden
	ReadCopyProb     - read model definition directly into Cplex
	
	if the problem is an LP
	LpOpt            - have Cplex solve the LP
	GetSolution      - populate the data structures with the LP solution
	
	or if the problem is a MIP
	MipOpt           - have Cplex solve the MIP
	GetMipSolution   - populate the data structures with the MIP solution
	
	wpPrintSoln      - private function to display objective function value and,
	                   optionally, variable values, reduced cost, pi, and slack
	SolWrite         - (optional) save the Cplex solution in a file
	ReadCopyProb     - (optional) save the model in a file
	CloseCplex       - clean up and close the Cplex environment


Initialize internal data structures

This option explicitly initializes the internal data structures. Initialization
is done automatically if solving the problem by reading a data file or populating
internal data structures, and need not be used in conjunction with those options.
However, if individual gpx functions are executed independently, it may be
necessary to explicitly initialize the input and solution data structures prior
to using some of those functions.

Populate input data structures

This option is used in conjunction with gpx functions executed independently in
order to populate the input data structures. Once populated, other gpx functions
can be used. This option is not needed if solving the problem, as the model is
automatically loaded into Cplex as part of that option.

Display input data

This option is used to display the model that resides in the internal data structures.
It is not needed for other operations, but is useful when running individual gpx
functions.

Display solution

This option is used to display the solution provided by Cplex. It is useful when
running individual gpx functions which do not automatically show the solution when
it is obtained.

TOGGLES

This section describes the toggles which control program behaviour.

Toggle gpx function exerciser

This toggle is used to enable or disable the options which are available to exercise
individual gpx functions. By default, the gpx function exerciser is disabled.
Changing the default state requires the following global variable to be changed
from "false" to "true":

	var gpxMenuOn  bool = false   // Flag for enabling gpx functions   


Toggle for custom environment

This toggle controls how file names are handled by this program. If all files are
located in the same directory and if all files have the same extension, this
option reduces the amount of typing needed to answer various prompts. By default,
the custom environment is disabled and the full file name (including path and
extension) must be specified. It is controlled by the following variable:

	var custEnvOn  bool = false    // Flag to enable custom paths and names

If custom environment is enabled (variable set to "true"), the directory name is
added as a prefix to the base file name, and the extension is added as a suffix
to that name. The default values are:

	var dSrcDev       string = "C:/Users/Admin/go/src/gpx/gpxrun/"  // Directory path
	var fExtension    string = ".txt"                               // File extension  

Caution is advised if using a custom environment.

FUNCTION WRAPPER

This section lists the options used to exercise individual gpx functions. Please
refer to the main documentation for details on function input, output, and behaviour.

The user interaction provided in the examples that solve problems from internal
data structures shows the functions being called, and the order in which they
are called. These functions, as well as several not used in the examples, can
be invoked independently via the options provided in this section. Care must be
taken that the required data structures have been correctly initialized and populated,
and that the functions are not called out of sequence.

The list of available functions, listed in alphabetical order, and some things to 
watch out for, are listed below.

   ChgCoefList     - Sets non-zero coefficients, must be used after NewCols and NewRows.
   ChgObjSen       - Sets problem to be treated as "maximize" or "minimize".
   ChgProbName     - Sets the problem name.
   CloseCplex      - Cleans up and closed the Cplex environment, must be called last.
   CreateProb      - Initializes the Cplex environment, must be called first.
   GetColName      - Creates the column solution list of the correct size and 
                     populates it with the column names.
   GetMipSolution  - Creates and populates solution structures for MIP problem.
   GetNumCols      - Gets the number of columns in the problem.
   GetNumRows      - Gets the number of rows in the problem.
   GetObjVal       - Gets the obj. func. value, assumes problem has been solved.
   GetRowName      - Creates the row solution list of the correct size and populates 
                     it with the row names.
   GetSlack        - Adds slack values to row solution list, which must exist.
   GetSolution     - Creates and populates solution structures for LP problem.
   GetX            - Adds values to column solution list, which must exist.
   LpOpt           - Optimizes an LP loaded into Cplex.
   MipOpt          - Optimizes a MIP loaded into Cplex.
   NewCols         - Creates new columns in Cplex from the internal data structures.
   NewRows         - Creates new rows in Cplex from the internal data structures.
   OutputToScreen  - Specifies whether Cplex should display output to screen or not.
   ReadCopyProb    - Populates the problem in Cplex directly from the file specified.
   SolWrite        - Writes the Cplex solution to a file.
   WriteProb       - Writes the problem loaded into Cplex to a file using the
                     format specified.



*/
package main
