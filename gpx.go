/*
Package gpx provides an interface to a small subset of CPX (C language) functions 
available in the CPLEX Callable Library. This package also requires a copy of Cplex
and a C compiler to be installed and configured.

The naming convention for gpx functions tries to match the underlying CPX function name as
closely as possible (e.g. gpx.LpOpt calls CPXlpopt). In some cases, it was not possible
and/or practical to have a one-for-one correspondence between gpx functions and the CPX
functions they call. 

The CPX function(s) called are listed in the comments section of the
relevant gpx function calling them. Please refer to the CPLEX documentation for details 
(https://www.ibm.com/support/knowledgecenter/en/SSSA5P_12.4.0/ilog.odms.cplex.help/CPLEX/maps/CPLEX_1.html).
*/
package gpx

/*
// Everything in comments above the import "C" is C code and will be compiled with the GCC. 
// Make sure you have a GCC installed.

#cgo LDFLAGS: -LD:/pk_cplex/include -lcplex1271

#include <string.h>
#include <stdio.h>
#include <D:/pk_cplex/include/ilcplex/cplex.h>

// Global variables local to this C section.

CPXENVptr env = NULL;
CPXLPptr lp = NULL;

#define BUFSIZE 80;

//==============================================================================
// C HELPER FUNCTIONS
//==============================================================================

//------------------------------------------------------------------------------
char**makeCharArray(int size) {
        return calloc(sizeof(char*), size);
}

char *makeNameStore(int size) {
	return malloc (size);	
}

static void setArrayString(char **a, char *s, int n) {
	a[n] = s;
}

static void freeCharArray(char **a, int size) {
	int i;
	for (i = 0; i < size; i++) {
		free(a[i]);			
	}
	free(a);
}

char *getArrayString(char **a, int n) {
	char *aItem;
				
	aItem = strdup(a[n]);
	return aItem;	
}

char *cGetArrayItem(char **a, int n) {
	int i;
	char *arrayItem;
	
	arrayItem = strdup(a[n]);
	return arrayItem;
}

//==============================================================================
// CPX FUNCTIONS
//==============================================================================

//------------------------------------------------------------------------------
int cOpenCplex() {
	
	int status = 0;	
	env = NULL;
	
    env = CPXopenCPLEX(&status);

	if (env == NULL) {
		char errmsg[CPXMESSAGEBUFSIZE];
		fprintf (stderr, "Could not open CPLEX environment.\n");
		CPXgeterrorstring (env, status, errmsg);
		fprintf (stderr, "%s", errmsg);
	}
		
	return status;	
}

//------------------------------------------------------------------------------
// Turn on output to the screen.
int cOutputToScreen(int state) {
	int status = 0;
	int echoState = 0;
	
	if (state) {
		echoState = CPX_ON;	
	} else {
		echoState = CPX_OFF;	
	}
	
	status = CPXsetintparam(env, CPXPARAM_ScreenOutput, echoState);
	if (status) {
		fprintf(stderr, "Failed to set screen output, error %d.\n", status);
	}
	return status;
}

//------------------------------------------------------------------------------
// Turn on data checking.
int cCheckData() {

	int status = 0;
	status = CPXsetintparam(env, CPXPARAM_Read_DataCheck, CPX_DATACHECK_WARN);
	if (status) {
		fprintf(stderr, "Failed to turn on data checking, error %d.\n", status);
	}
	return status;	
}


//------------------------------------------------------------------------------
// Create the LP.
int cCreateProb(char *probName) {

	char buffer[256];
	int status = 0;
		
	strcpy(buffer, probName);
	
	lp = CPXcreateprob(env, &status, buffer);
	
	if (lp == NULL) {
		fprintf(stderr, "Unable to create problem.\n");	
	}
	
	return status;
}

//------------------------------------------------------------------------------
// Create the LP.
int cChgProbName(char *probName) {

	char buffer[256];
	int status = 0;
		
	strcpy(buffer, probName);
	
	status = CPXchgprobname(env, lp, buffer);
	
	if (status) {
		fprintf(stderr, "Unable to change problem name, error %d.\n", status);	
	}
	
	return status;
}

//------------------------------------------------------------------------------
// Change objective function sense (default is minimize).
int cChgObjSen(int state) {
	int status = 0;
	int sense;
	
	if (state == 1) {
		sense = CPX_MIN;	
	} else {
		sense = CPX_MAX;	
	}
	
	status = CPXchgobjsen(env, lp, sense);
	if (status) {
		fprintf(stderr, "Failed to set objective function sense, error %d.\n", status);
	}
	return status;
}

//------------------------------------------------------------------------------
// Create new rows
int cCreateRows(int numRows, char *senseArray, char **rowName, double *rhs, double *rngVal) {

	int status = 0;
			
	status = CPXnewrows(env, lp, numRows, rhs, senseArray, rngVal, rowName);
	if (status) {
		fprintf(stderr, "Failed to create rows, error %d.\n", status);
	}

	return status;	
}


//------------------------------------------------------------------------------
// Create new columns
int cCreateCols(int isMip, int numCols, double *obj, char **colName,  char *type, double *lb, double *ub) {

	int status = 0;

	// If this is a MIP, we must pass tye variable type array to Cplex.
	// Otherwise, we must pass a NULL even if all variables are of the C type so CPXlpopt
	// does not complain when we solve the problem.
	
	if (isMip) {
		status = CPXnewcols(env, lp, numCols, obj, lb, ub, type, colName);	
	} else {	
		status = CPXnewcols(env, lp, numCols, obj, lb, ub, NULL, colName);	
    }
	
	if (status) {
		fprintf(stderr, "Failed to create columns, error %d.\n", status);
	}
	
	return status;	
}

//------------------------------------------------------------------------------
// Create new columns
int cChgCoefList(int numNZ, int *rowlist, int *collist, double *vallist) {

	int status = 0;
	
	status = CPXchgcoeflist(env, lp, numNZ, rowlist, collist, vallist); 
	if (status) {
		fprintf(stderr, "Failed to change coefficients, error %d.\n", status);
	}
	
	return status;	

}

//------------------------------------------------------------------------------
// Optimize a linear problem (continuous variables only).
int cLpOpt() {

	int status = 0;
	
	status = CPXlpopt(env, lp);	
	if (status) {
		fprintf(stderr, "CPXlpopt failed with error %d.\n", status);
	}

	return status;	
}

//------------------------------------------------------------------------------
// Optimize a mixed integer problem.
int cMipOpt() {

	int status = 0;
	
	status = CPXmipopt(env, lp);	
	if (status) {
		fprintf(stderr, "CPXmipopt failed with error %d.\n", status);
	}

	return status;	
}

//------------------------------------------------------------------------------
// Get number of rows
int cGetNumRows(int *numRows) {

	int cur_numrows;
	
	*numRows = 0;
	cur_numrows = CPXgetnumrows(env, lp);
	
	*numRows = cur_numrows;
	
	return 0;
}

//------------------------------------------------------------------------------
// Get number of columns
int cGetNumCols(int *numCols) {

	int cur_numcols;

	*numCols = 0;	
	cur_numcols = CPXgetnumcols(env, lp);
	
	*numCols = cur_numcols;
	
	return 0;
}

//------------------------------------------------------------------------------
int cGetColNameSurplus(int numCols, int *surplus) {
	int status;
	
	status = CPXgetcolname(env, lp, NULL, NULL, 0, surplus, 0, numCols - 1);
    if (( status != CPXERR_NEGATIVE_SURPLUS ) && ( status != 0 ))  {
    	fprintf (stderr, "Could not determine amount of space for column names.\n");
    	return status;
	}

	// Switch "negative surplus 'error'" to "this is not an error".
	return 0;	
}

//------------------------------------------------------------------------------
int cGetRowNameSurplus(int numRows, int *surplus) {
	int status;
	
	status = CPXgetrowname(env, lp, NULL, NULL, 0, surplus, 0, numRows - 1);
    if (( status != CPXERR_NEGATIVE_SURPLUS ) && ( status != 0 ))  {
    	fprintf (stderr, "Could not determine amount of space for row names.\n");
    	return status;
	}

	// Switch "negative surplus 'error'" to "this is not an error".
	return 0;	
}

//------------------------------------------------------------------------------
// Get objective function value
int cGetObjVal(double *objval) {

	int status;

	status = CPXgetobjval(env, lp, objval);
	if (status) {
		fprintf(stderr, "Failed to obtain objective value, error %d.\n", status);
		return status;
	}	
		
	return status;
}

//------------------------------------------------------------------------------
// Get Solution
int cGetSolution(double *objval, double *x, double *dj, double *pi, double *slack) {

	int status, solstat;

	status = CPXsolution(env, lp, &solstat, objval, x, pi, slack, dj);
	if (status) {
		fprintf(stderr, "Failed to obtain solution, error %d.\n", status);
		return status;
	}	
		
	return status;
}

//------------------------------------------------------------------------------

int cGetColNames(int numCols, char **cur_colname, char *cur_colnamestore, int storeSize) {
	
	int status;
	int surplus, cur_colspace;
	int j;	

	if (storeSize > 0) {

	    status = CPXgetcolname (env, lp, cur_colname, cur_colnamestore, 
					storeSize, &surplus, 0,numCols-1);

		if (status) {
	    	fprintf (stderr, "CPXgetcolname failed.\n");
			return status;	
    	}
	}
	else {
    	fprintf (stderr, "No names associated with problem\n");
		return 1;
	}

	return status;	
}

//------------------------------------------------------------------------------

int cGetRowNames(int numRows, char **cur_rowname, char *cur_rownamestore, int storeSize) {
	
	int status = 0;
	int surplus;

	if (storeSize > 0) {

	    status = CPXgetrowname (env, lp, cur_rowname, cur_rownamestore, 
					storeSize, &surplus, 0,numRows-1);

		if (status) {
	    	fprintf (stderr, "CPXgetrowname failed.\n");
			return status;	
    	}
	}
	else {
    	fprintf (stderr, "No names associated with problem\n");
		return 1;
	}

	return status;	
}


//------------------------------------------------------------------------------
// Get the value of the variables
int cGetX(int numCols, double *cColVal) {

	int status = 0;
		
	status = CPXgetx(env, lp, cColVal, 0, numCols - 1);
	if ( status ) {
		fprintf (stderr, "Failed to get column values, error %d.\n", status);
		return status;
	}

	return status;	
}

//------------------------------------------------------------------------------
// Get the slack of the rows
int cGetSlack(int numRows, double *cSlack) {

	int status = 0;
		
	status = CPXgetslack(env, lp, cSlack, 0, numRows - 1);
	if ( status ) {
		fprintf (stderr, "Failed to get slack values, error %d.\n", status);
		return status;
	}

	return status;	
}


//------------------------------------------------------------------------------

int cReadCopyProb(char *cFileName, char *cFileType) {

	int status = 0;
		
	status = CPXreadcopyprob(env, lp, cFileName, cFileType);
	if ( status ) {
		fprintf (stderr, "Failed to read and copy the problem data, error %d.\n", status);
		return status;
	}

	return status;	
}

//------------------------------------------------------------------------------

int cWriteProb(char *cFileName, char *cFileType) {

	int status = 0;
		
	status = CPXwriteprob(env, lp, cFileName, cFileType);
	if ( status ) {
		fprintf(stderr, "Failed to write problem, error %d.\n", status);
		return status;
	}

	return status;	
}

//------------------------------------------------------------------------------

int cSolWrite(char *cFileName) {
	
	int status = 0;
	
	status = CPXsolwrite(env, lp, cFileName);
	if ( status ) {
		fprintf (stderr, "Failed to write solution file.\n");
		return status;
	}
	
	return status;

}

//------------------------------------------------------------------------------
// Clean up and terminate CPLEX. This should be called from the Go functions
// and not from the C functions if an error condition occurs.
int cCloseCplex() {

	int status = 0;
	
// Free up the solution (more stuff goes here later).

	if (lp != NULL) {
		status = CPXfreeprob (env, &lp);
		if (status) {
			fprintf(stderr, "CPXfreeprob failed with error %d.\n", status);	
		}
	}

	if (env != NULL) {
		status = CPXcloseCPLEX(&env);
		if (status) {
			fprintf(stderr, "CPXcloseCPLEX failed with error %d.\n", status);
		}	
	}

	return status;
	
} // End cTerminate


 */
import "C"

import (
//	"fmt"
	"github.com/pkg/errors"
	"unsafe"
)


var plInfySmall = 1.0e10          // Default infinity in lpo package
var plInfyLarge = 1.0e20          // Default CPX_INFBOUND in Cplex

// Input data structures passed to gpx to define a problem.

// InputRow defines a data structure passed as an
// input argument to functions when creating the rows of the problem in Cplex.
type InputRow struct {
	Name   string      // Name of the row (constraint)
	Sense  string      // Sense (L, E, G, R) of the row as supported by Cplex
	Rhs    float64     // Value of the RHS, or lower boundary of the range
	RngVal float64     // For ranges, the range is defined as (Rhs to [Rhs + RngVal])	
}

// InputCol defines a data structure passed as an input
// argument to functions when creating the columns of the problem in Cplex.
type InputCol struct {
	Name   string      // Name of the column (variable)
	Type   string      // Type of the column as supported by Cplex
	BndLo  float64     // Lower bound of the column
	BndUp  float64 	   // Upper bound of the column
}

// InputElem defines a data structure passed as an
// input argument to functions when changing non-zero coefficients of the problem
// in Cplex.
type InputElem struct {
	RowIndex  int      // Row index for this element (coefficient)
	ColIndex  int      // Column index for this element
	Value     float64  // Value of this element
}

// InputObjCoef defines a data structure of all coefficients
// present in the objective function. It is passed as an input argument to functions
// which create the columns of the problem in Cplex.
type InputObjCoef struct {
	ColIndex  int      // Column index of this coefficient in the objective function
	Value     float64  // Value of the coefficient in the objective function
}

// Output data structures which contain the solution provided by Cplex.

// SolnRow defines a data structure of the solved rows returned from Cplex.
type SolnRow struct {
	Name    string     // Name of the row in the solution data structure
	Slack   float64    // Slack for this row as calculated by Cplex
	Pi      float64    // Pi for this row as calculated by Cplex
}

// SolnCol defines a data structure of the solved columns returned from Cplex.
type SolnCol struct {
	Name    string     // Name of the column in the solution data structure
	Value   float64    // Value for this column as calculated by Cplex
	RedCost float64    // Reduced cost for this column as calculated by Cplex
}

//==============================================================================
// FUNCTIONS FOR CREATING THE PROBLEM
//==============================================================================

// NewProblem initializes the Cplex environment and creates a problem with the
// name passed into this function. 
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXopenCPLEX, CPXsetintparam, and CPXcreateprob.
func CreateProb(Name string) error {
	var status C.int

	if Name == "" {
		Name = "DefaultName"
	}
			
	status = C.cOpenCplex()
	if status != 0 {
		return errors.Errorf("Cplex open failed with error %d", status)	
	}

	status = C.cCheckData()
	if status != 0 {
		return errors.Errorf("Enabling data checking failed with error %d", status)
	}	

	cString := C.CString(Name)
	defer C.free(unsafe.Pointer(cString))
	status = C.cCreateProb(cString)
	if status != 0 {
		return errors.Errorf("Creating problem failed with error %d", status)
	}	
			
	return nil
}

//==============================================================================

// OutputToScreen turns Cplex output to screen on or off. 
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXsetintparam with CPXPARAM_ScreenOutput set to CPX_ON if the value
// passed to this function is "true", or CPX_OFF if false. By default, output is not
// printed to the screen.
func OutputToScreen(echoOn bool) error {

	var status, cEchoState C.int

	// CPX_ON = 1, CPX_OFF = 0, but unfortunately we can't call these constants
	// here, and to be on the safe side, we will re-map on the C side in case Cplex changes.		
	if echoOn {
		cEchoState = 1	
	} else {
		cEchoState = 0
	}
	
	status = C.cOutputToScreen(cEchoState)
	if status != 0 {
		return errors.Errorf("Cplex failed to turn on screen output with error %d", status)
	}	

	return nil	
}

//==============================================================================

// NewProblem initializes the Cplex environment and creates a problem with the
// name passed into this function.
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXchgprobname.
func ChgProbName(Name string) error {
	var status C.int

	if Name == "" {
		Name = "ChangedDefaultName"
	}
			
	cString := C.CString(Name)
	defer C.free(unsafe.Pointer(cString))
	status = C.cChgProbName(cString)
	if status != 0 {
		return errors.Errorf("Changing problem name failed with error %d", status)
	}	
			
	return nil
}

//==============================================================================

// ChgObjSen changes the sense of the objective function depending on the value
// specified by the user.
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXchgobjsen. By default, the sense is set to "minimize".
// 	Supported values of sense are:
//		 1 - minimize
//		-1 - maximize
func ChgObjSen(sense int) error {

	var status C.int

	// CPX_ON = 1, CPX_OFF = 0, but unfortunately we can't call these constants
	// here, and to be on the safe side, we will re-map on the C side in case Cplex changes.
	
	switch sense {
		case -1, 1:
			status = C.cChgObjSen(C.int(sense))
			if status != 0 {
				return errors.Errorf("Failed to change objective sense, error %d", status)
			}	
		
		default:
			return errors.Errorf("Unexpected objective function sense %d", sense)
	}
	
	return nil	
}

//==============================================================================

// NewRows creates the new rows specified for this problem. In case of failure,
// it returns an error including the error code it received from Cplex. 
// This function uses CPXnewrows.
//	The values supported in the Sense field of the rList structure are:
//		L - less than or equal
//		E - equal
//		G - greater than or equal
//		R - range
//	If a constraint has rList[i].Sense = "R", the value of constraint i can be
//	between rList[i].Rhs and (rList[i].Rhs + rList[i].RngVal). For all other cases,
//	RngVal is set to zero.
func NewRows(rList []InputRow) error {
	var nameArray  []string	
	var cChar        C.char
	var cCharArray []C.char
	var status       C.int
	var cRhs       []C.double
	var cRngVal    []C.double


	if len(rList) < 1 {
		return errors.Errorf("NewRows expected more than %d rows", len(rList))	
	}

	// Build C array of row names
	for i := 0; i < len(rList); i++ {
		nameArray = append(nameArray, rList[i].Name)
	}			

	cNameArray := C.makeCharArray(C.int(len(rList)))
	defer C.freeCharArray(cNameArray, C.int(len(rList)))

	for i, s := range nameArray {
		cString := C.CString(s)
        C.setArrayString(cNameArray, cString, C.int(i))
		// The cString pointers are freed as part of freeCharArray function, not here.
	}

	// Construct the lists of the other C parameters	
	for i := 0; i < len(rList); i++ {
		runes     := []rune(rList[i].Sense)
		cChar      = C.char(runes[0])
		cCharArray = append(cCharArray, cChar)		
		cRhs       = append(cRhs, C.double(rList[i].Rhs))
		cRngVal    = append(cRngVal, C.double(rList[i].RngVal))
	} // End for list of rows
		
	status = C.cCreateRows(C.int(len(cRhs)), &cCharArray[0], cNameArray, &cRhs[0], &cRngVal[0])
	if status != 0 {
		return errors.Errorf("Creating rows failed with error %d", status)
	}	
	
	return nil
	
} // End NewRows

//==============================================================================

// NewCols creates the new columns specified for this problem. In case of failure,
// it returns an error including the error code it received from Cplex. 
// This function uses CPXnewcols.
//	The Type field in the cList data structure supports the following values:
//		C - continuous variable (CPX_CONTINUOUS); only value supported by CPXlpopt
//		B - binary variable (CPX_BINARY)
//		I - general integer variable (CPX_INTEGER)
//		S - semi-continuous variable (CPX_SEMICONT)
//		N - semi-integer variable (CPX_SEMIINT)
//	Any value other than 'C' (continuous) is interpreted by Cplex as a MIP. If
//	function CPXlpopt is called for a problem containing anything other than
//	continuous variables, it will fail with a CPXERR_NOT_FOR_MIP error. 
func NewCols(objList []InputObjCoef, cList []InputCol) error {

	var nameArray  []string	
	var cChar        C.char
	var status       C.int
	var cCharArray []C.char
	var lb, ub     []C.double
	var isMip        C.int 

	// The column list must be provided. 
	if len(cList) < 1 {
		return errors.Errorf("NewCols expected more than %d columns", len(cList))			
	}

	// Build C array of objective function coefficients.
	obj := make([]C.double, len(cList))

	for i := 0; i < len(objList); i++ {
		obj[objList[i].ColIndex] = C.double(objList[i].Value)
	}	

	// Build C array of column names				
	for i := 0; i < len(cList); i++ {
		nameArray = append(nameArray, cList[i].Name)
	}			

	cNameArray := C.makeCharArray(C.int(len(nameArray)))
	defer C.freeCharArray(cNameArray, C.int(len(nameArray)))

	for i, s := range nameArray {
		cString := C.CString(s)
        C.setArrayString(cNameArray, cString, C.int(i))
		// The cString pointers are freed as part of freeCharArray function, not here.
	}

	
	// Build C array of upper and lower bounds of the variables.
	// Due to implementation of CPXnewrows, we also need to check whether there are
	// any non-contiguous variables so entire array (as opposed to NULL) can be passed
	// to Cplex.
	
	isMip = 0
				
	for i := 0; i < len(cList); i++ {
		runes     := []rune(cList[i].Type)
		cChar      = C.char(runes[0])
		cCharArray = append(cCharArray, cChar)		
		lb         = append(lb, C.double(cList[i].BndLo))
		ub         = append(ub, C.double(cList[i].BndUp))
		
		if cList[i].Type != "C" {
			isMip = 1
		}										
	}	

	// Call the C function which passes the arrays to Cplex.	
	status = C.cCreateCols(isMip, C.int(len(cList)), &obj[0], cNameArray, &cCharArray[0], &lb[0], &ub[0])
	if status != 0 {
		return errors.Errorf("Creating columns failed with error %d", status)
	}	
			
	return nil
}

//==============================================================================

// ChgCoefList modifies the non-zero coefficients specified for the problem. 
// The rows and columns, created by other functions, are assumed to exist. 
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXchgcoeflist.
func ChgCoefList(eList []InputElem) error {

	var rowlist, collist []C.int
	var vallist []C.double
	var status C.int
		
	if len(eList) < 1 {
		return errors.Errorf("ChgCoefList expected more than %d elements", len(eList))			
	}
	
	for i := 0; i < len(eList); i++ {
		
		rowlist = append(rowlist, C.int(eList[i].RowIndex))
		collist = append(collist, C.int(eList[i].ColIndex))
		vallist = append(vallist, C.double(eList[i].Value))
		
	} // End for all rows	

    status = C.cChgCoefList(C.int(len(rowlist)), &rowlist[0], &collist[0], &vallist[0])
	if status != 0 {
		return errors.Errorf("Changing coefficients failed with error %d", status)
	}	
	
	return nil
}

//==============================================================================
// FUNCTIONS FOR SOLVING THE PROBLEM AND GETTING SOLUTION (Get series functions)
//==============================================================================

// LpOpt solves the LP, which must is assumed to have been defined by
// other functions. 
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXlpopt.
//
// The model can contain only continuous ('C') variables. The presence
// of any other variable type will cause this function to fail with a CPXERR_NOT_FOR_MIP
// error.
func LpOpt() error {

	var status C.int
	
	status = C.cLpOpt()
	if status != 0 {
		return errors.Errorf("Error %d received from cLpOpt", status)
	}	
	
	return nil	
}

//==============================================================================

// MipOpt solves the mixed integer problem, which is assumed to have been defined by
// other functions. 
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXmipopt.
func MipOpt() error {

	var status C.int
	
	status = C.cMipOpt()
	if status != 0 {
		return errors.Errorf("Error %d received from cMipOpt", status)
	}	
	
	return nil	
}

//==============================================================================

// GetSolution obtains the solution for a linear problem (LP) from Cplex. 
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXgetsolution as well as all other auxiliary
// functions, such as CPXgetrowname, CPXgetcolname, and others that may be needed
// in order to populate the solution data structures passed back to the user.
func GetSolution(objVal *float64, sRows *[]SolnRow, sCols *[]SolnCol) error {

	var cObjVal  C.double
	var status   C.int
	var numRows  C.int
	var numCols  C.int
	
	var curCol SolnCol
	var curRow SolnRow
	var err error

	// Get actual number of rows and columns and allocate memory for solution.
	_ = C.cGetNumRows(&numRows)
	_ = C.cGetNumCols(&numCols)

	cXval  := make([]C.double, numCols)
	cRcost := make([]C.double, numCols)

	cPi    := make([]C.double, numRows)
	cSlack := make([]C.double, numRows)
	

	// Get the solution using the C data structures.				
	status = C.cGetSolution(&cObjVal, &cXval[0], &cRcost[0], &cPi[0], &cSlack[0])	
	if status != 0 {
		return errors.Errorf("Error %d received from cGetSolution", status)
	}	

	*objVal = float64(cObjVal)

	for i := 0; i < int(numCols); i++ {
		curCol.Value   = float64(cXval[i])
		curCol.RedCost = float64(cRcost[i])
		*sCols = append(*sCols, curCol)	
	}
	
	if err = GetColName(*sCols); err != nil {
		return errors.Wrap(err, "GetSolution failed to get column names")
	}

	for i := 0; i < int(numRows); i++ {
		curRow.Pi   = float64(cPi[i])
		curRow.Slack = float64(cSlack[i])
		*sRows = append(*sRows, curRow)
	}	

	if err = GetRowName(*sRows); err != nil {
		return errors.Wrap(err, "GetSolution failed to get row names")
	}
				
	return nil

}

//==============================================================================

// GetMipSolution obtains the solution for a mixed integer problem (MIP) from Cplex. 
// In case of failure, it returns an error including the error code it received from Cplex.
// This function creates the slices needed for the result and populates all
// applicable fields with values provided by Cplex.
// This function uses other gpx functions and indirectly call CPXgetnumrows, 
// CPXgetnumcols, CPXgetrowname, CPXgetcolname, CPXgetobjval, CPXgetslack, and CPXgetx.
func GetMipSolution(objVal *float64, sRows *[]SolnRow, sCols *[]SolnCol) error {

	var numRows, numCols  int
	var err error

	// Initialize the return values
	*sRows  = nil
	*sCols  = nil
	*objVal = 0.0
	
	// Get the number of rows and columns, and allocate the memory for them.
	if err = GetNumRows(&numRows); err != nil {
		return errors.Wrap(err, "GetMipSolution failed to get number of rows") 
	}

	if err = GetNumCols(&numCols); err != nil {
		return errors.Wrap(err, "GetMipSolution failed to get number of columns") 
	}
	
	*sRows = make([]SolnRow, numRows)
	*sCols = make([]SolnCol, numCols)
	
	if err = GetObjVal(objVal); err != nil {
		return errors.Wrap(err, "GetMipSolution failed to get objective function value") 		
	}	


	if err = GetRowName(*sRows); err != nil {
		return errors.Wrap(err, "GetMipSolution failed to get row names") 		
	}	

	if err = GetSlack(*sRows); err != nil {
		return errors.Wrap(err, "GetMipSolution failed to get row slack") 		
	}	

	
	if err = GetColName(*sCols); err != nil {
		return errors.Wrap(err, "GetMipSolution failed to get column names") 		
	}	

	if err = GetX(*sCols); err != nil {
		return errors.Wrap(err, "GetMipSolution failed to get column values") 		
	}	
	
	return nil

}

//==============================================================================

// GetNumRows obtains the number of rows in the current problem, or 0 if none exist
// or the problem has not yet been defined. At this time the function always returns
// nil (success).
// This function uses CPXgetnumrows. 
func GetNumRows(numRows *int) error {
	var cNumRows C.int
	
	_ = C.cGetNumRows(&cNumRows)

	*numRows = int(cNumRows)	

	return nil	
}

//==============================================================================

// GetNumCols obtains the number of columns in the current problem, or 0 if none exist
// or the problem has not yet been defined. At this time the function always returns
// nil (success).
// This function uses CPXgetnumcols. 
func GetNumCols(numCols *int) error {
	var cNumCols C.int
	
	_ = C.cGetNumCols(&cNumCols)

	*numCols = int(cNumCols)	

	return nil	
}

//==============================================================================

// GetObjVal obtains the value of the objective function. 
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXgetobjval.
func GetObjVal(objVal *float64) error {
	var cObjVal C.double
	var status  C.int

	*objVal = 0
		
	status = C.cGetObjVal(&cObjVal)
	if status != 0 {
		return errors.Errorf("GetObjVal failed with error %d", status)		
	}
	*objVal = float64(cObjVal)	

	return nil	
}

//==============================================================================

// GetColName populates the solution column slice passed to the function with
// the column names being used by Cplex. It does not populate any other fields in
// the data structures, and assumes that the slice is large enough to contain all
// of the column names (i.e. it does not modify the size in case of mismatch).
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXgetnumcols and CPXgetcolname.
func GetColName(sCols []SolnCol) error {
	var numCols, status, surplus, colSpace C.int

	// Get actual number of columns and allocate memory for solution.
	_ = C.cGetNumCols(&numCols)

	if len(sCols) != int(numCols) {
		return errors.Errorf("Number of cols in target is %d, expected %d", 
			len(sCols), numCols)		
	}
	
    status = C.cGetColNameSurplus(numCols, &surplus)
	if status != 0 {
		return errors.Errorf("GetColNameSurplus failed with error %d", status)
	}	

	colSpace = -surplus	

	// Create memory for the name array.
	cColName := C.makeCharArray(numCols)
	defer C.free(unsafe.Pointer(cColName))
	cColNameStore := C.makeNameStore(colSpace)
	defer C.free(unsafe.Pointer(cColNameStore))

	status = C.cGetColNames(numCols, cColName, cColNameStore, colSpace)
	if status != 0 {
		return errors.Errorf("Get col names failed with error %d", status)
	}	
	
	for i := 0; i < int(numCols); i++ {
		cString := C.cGetArrayItem(cColName, C.int(i))	
		sCols[i].Name = C.GoString(cString)
		C.free(unsafe.Pointer(cString))			
	}

	return nil
}


//==============================================================================

// GetRowName populates the solution row slice passed to the function with
// the row names being used by Cplex. It does not populate any other fields in
// the data structures, and assumes that the slice is large enough to contain all
// of the row names (i.e. it does not modify the size in case of mismatch).
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXgetnumrows and CPXgetrowname.
func GetRowName(sRows []SolnRow) error {
	var numRows, status, surplus, rowSpace C.int

	// Get actual number of rows and allocate memory for solution.
	_ = C.cGetNumRows(&numRows)

	if len(sRows) != int(numRows) {
		return errors.Errorf("Number of rows in target is %d, expected %d", 
			len(sRows), numRows)		
	}
	
    status = C.cGetRowNameSurplus(numRows, &surplus)
	if status != 0 {
		return errors.Errorf("GetRowNameSurplus failed with error %d", status)
	}	

	rowSpace = -surplus
	
	// Create memory for the name array.
	cRowName := C.makeCharArray(numRows)
	defer C.free(unsafe.Pointer(cRowName))
	cRowNameStore := C.makeNameStore(rowSpace)
	defer C.free(unsafe.Pointer(cRowNameStore))

	status = C.cGetRowNames(numRows, cRowName, cRowNameStore, rowSpace)
	if status != 0 {
		return errors.Errorf("Get row names failed with error %d", status)
	}	
	
	for i := 0; i < int(numRows); i++ {
		cString := C.cGetArrayItem(cRowName, C.int(i))	
		sRows[i].Name = C.GoString(cString)
		C.free(unsafe.Pointer(cString))			
	}
		
	return nil
	
}

//==============================================================================

// GetX populates the solution column slice passed to the function with
// the optimal value of the variable as calculated by Cplex. It does not populate 
// any other fields in the data structures, and assumes that the slice is large 
// enough to contain all data (i.e. it does not modify the size in case of mismatch).
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXgetnumcols and CPXgetx.
func GetX(sCols []SolnCol) error {
	var cNumCols, status C.int
	var numCols int

	// Get actual number of columns and allocate memory for solution.
	_ = C.cGetNumCols(&cNumCols)

	if len(sCols) < int(cNumCols) {
		return errors.Errorf("Number of cols in target is %d, expected %d", 
			len(sCols), cNumCols)		
	}

	cXval  := make([]C.double, cNumCols)

	// Get the solution using the C data structures.				
	status = C.cGetX(cNumCols, &cXval[0])	
	if status != 0 {
		return errors.Errorf("Error %d received from cGetX", status)
	}	

	// Transfer the column value from the C structure to the slice passed to us.
	numCols = int(cNumCols)
	for i := 0; i < numCols; i++ {
		sCols[i].Value = float64(cXval[C.int(i)])
	}

	return nil
}

//==============================================================================

// GetSlack populates the solution row slice passed to the function with
// the slack value of the constraint as calculated by Cplex. It does not populate 
// any other fields in the data structures, and assumes that the slice is large 
// enough to contain all data (i.e. it does not modify the size in case of mismatch).
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXgetnumrows and CPXgetslack.
func GetSlack(sRows []SolnRow) error {
	var cNumRows, status C.int
	var numRows int

	// Get actual number of rows and allocate memory for solution.
	_ = C.cGetNumRows(&cNumRows)

	if len(sRows) < int(cNumRows) {
		return errors.Errorf("Number of rows in target is %d, expected %d", 
			len(sRows), cNumRows)		
	}

	cSlack  := make([]C.double, cNumRows)

	// Get the solution using the C data structures.				
	status = C.cGetSlack(cNumRows, &cSlack[0])	
	if status != 0 {
		return errors.Errorf("Error %d received from cGetSlack", status)
	}	

	// Transfer the row slack from the C structure to the slice passed to us.
	numRows = int(cNumRows)
	for i := 0; i < numRows; i++ {
		sRows[i].Slack = float64(cSlack[C.int(i)])
	}

	return nil
}

//==============================================================================
// FUNCTIONS FOR PROCESSING FILES AND MISCELANEOUS FUNCTIONALITY
//==============================================================================

// CloseCplex closes and cleans up the Cplex environment. 
// In case of failure, it returns an error including the error code it received from Cplex. 
// This function uses CPXfreeprob and CPXcloseCPLEX.
func CloseCplex() error {
	var status C.int
		
	status = C.cCloseCplex()
	if status != 0 {
		return errors.Errorf("Close Cplex failed with error %d", status)	
	}
	
	return nil
}

//==============================================================================

// ReadCopyProb reads the data file specified by its name and type, and populates
// the problem from information contained in this file.
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXreadcopyprob.
//	The following fileType values are supported:
//		SAV - binary format
//		MPS - MPS format
//		LP  - LP format
func ReadCopyProb(fileName string, fileType string) error {

	var status C.int

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))
	cFileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cFileType))
	
	status = C.cReadCopyProb(cFileName, cFileType)
	if status != 0 {
		return errors.Errorf("Read file failed with error %d", status)
	}	

	return nil	
}

//==============================================================================

// WriteProb writes current Cplex problem to the data file specified by its 
// name and type.
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXwriteprob.
//	The following fileType values are supported:
//		SAV - binary matrix and basis file
//		MPS - MPS format
//		LP  - CPLEX LP format with names modified to conform to LP format
//		REW - MPS format with all names changed to generic names
//		ALP - LP format with generic name, type, and bound of each variable
//	If the file name ends with one of the following extensions, a compressed
//	file is written:
//		.bz2 - files compressed with BZip2
//		.gz  - files compressed with GNU Zip
func WriteProb(fileName string, fileType string) error {

	var status C.int

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))
	cFileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cFileType))
	
	status = C.cWriteProb(cFileName, cFileType)
	if status != 0 {
		return errors.Errorf("Write problem file failed with error %d", status)
	}	

	return nil	
}

//==============================================================================

// WriteProb writes the current Cplex problem to the data file specified by its 
// name and type.
// In case of failure, it returns an error including the error code it receives
// from Cplex.
// This function uses CPXwriteprob.
func SolWrite(fileName string) error {
	
	var status C.int

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	status = C.cSolWrite(cFileName)
	if status != 0 {
		return errors.Errorf("Writing solution file failed with error %d", status)
	}	

	return nil	
		
}

//============================ END OF FILE =====================================
