### Package gpx

Package gpx provides a Go language suite of tools for using a subset of callable C functions available for the
Cplex solver. 

This package is intended to be used in conjunction with the separate LPO package 
which provides tools for modelling and solving Linear Programming (LP) and Mixed-Integer 
Linear Programming (MILP) problems. Package gpx assumes that the Cplex header files
and object files as well as a compatible C compiler have been installed and configured
separately on the computer where this package is to be used.

Once gpx has been downloaded and Cplex has been installed and configured, you will need to modify the
current placeholders for the Cplex object file (or dll) and header file with the current location.
The two lines which must be changed in gpx.go are:

#cgo LDFLAGS: -LD:/pk_cplex/include -lcplex1271

#include <D:/pk_cplex/include/ilcplex/cplex.h>

### Executable gpxrun

The subdirectory gpxrun contains the executable which illustrates the functionality of the
gpx package and how it can be used to solve LP and MILP problems via Cplex. The executable
also includes a wrapper which can be used to independently exercise each function exported
by gpx.
