// 01   July  29, 2018   Created separate doc file, added comments from J. Chinneck


/*
Package gpx ("linear programming Go for Cplex") provides an interface to a small subset of CPX 
(C language) functions available in the CPLEX Callable Library. It is complementary to,
but independent of the lpo ("linear programming object") package.

The gpx package requires a copy of Cplex and a C compiler to be installed and configured. 
It can be used from Go language programs to create linear programming or 
mixed-integer linear programming models, or read them from files, modify them, 
solve them, and obtain the solution results. It can also be used from the lpo package 
(https://github.com/go-opt/lpo), which has a number of routines for retrieving 
and presolving linear models.

The naming convention for gpx functions tries to match the underlying CPX function name as
closely as possible (e.g. gpx.LpOpt calls CPXlpopt). In some cases, it was not possible
and/or practical to have a one-for-one correspondence between gpx functions and the CPX
functions they call. 

The CPX function(s) called are listed in the comments section of the
relevant gpx function calling them. Please refer to the CPLEX documentation for details 
(https://www.ibm.com/support/knowledgecenter/en/SSSA5P_12.4.0/ilog.odms.cplex.help/CPLEX/maps/CPLEX_1.html).

The executable provided with the package illustrates how the gpx package can be
used and contains an exerciser to allow each function to be tested independently.

Dependencies

Package gpx requires the following:
  - Cplex to be installed and configured correctly.
  - C compiler to be installed and configured correctly.
  - Package github.com/pkg/errors to be installed.

Example

The wpSolveFromFile function in the gpxrun executable shows which functions
need to be called, and the order in which to call them, if solving a problem from
a data file. Similarly, the wpSolveFromStruct function in gpxrun shows functions, 
and the order in which to call them, if solving a problem defined in the gpx 
data structures.

A simplified version of the wpSolveFromFile for reading an MPS data file and
solving it via Cplex would include the following statements:

  ...
  var sObjVal       float64   // solution value of objective function
  var sRows   []gpx.SolnRow   // solution rows provided by gpx
  var sCols   []gpx.SolnCol   // solution columns provided by gpx
  var fileName       string   // name of file containing model
  var fileType       string   // type of input file (MPS, LP, or SAV)
  var probName       string   // name of problem (Cplex does not get this from the MPS file)
  var dispToScreen     bool   // flag instructing Cplex to display output to screen
  var err             error   // string concatenating error conditions that occurred
  
  // Initialize variables.
  fileName     = "c:/myDataFile.txt"   // provide full path for file to be read
  fileType     = "MPS"                 // set file type to MPS
  probName     = "MyTestProblem"       // provide the name of the problem
  dispToScreen = true                  // print output to screen

  // Set up the Cplex environment and set problem name.
  if err = gpx.CreateProb(probName); err != nil {
       return errors.Wrap(err, "Failed to initialize environment")				
  }

  // Set the display to screen parameter.
  if err = gpx.OutputToScreen(dispToScreen); err != nil {
       return errors.Wrap(err, "Failed to set display to screen")						
  }

  // Read the problem from input file.
  if err = gpx.ReadCopyProb(fileName, fileType); err != nil {
      return errors.Wrap(err, "Opening input file failed")
  } 

  // Solve the problem as an LP.
  if err = gpx.LpOpt(); err != nil {
      return errors.Wrap(err, "Failed to solve the LP")
  }		

  // Populate the gpx data structures with the solution provided by Cplex.
  err = gpx.GetSolution(&sObjVal, &sRows, &sCols)
  if err != nil {
      return errors.Wrap(err, "Failed to get LP solution")
  } 

  // Close and clean up the Cplex environment.
  if err = gpx.CloseCplex(); err != nil {
      return errors.Wrap(err, "Failed to close Cplex")
  }

  // Do something useful with the solution. This will be defined by you to
  // print the solution, translate to different format, or perform some other task.
  processGpxSolution(sObjVal, sRows, sCols)

  ...


*/
package gpx
