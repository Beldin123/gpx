# Package gpx

Package gpx provides a Go language suite of tools for using a subset of callable C functions available for the
Cplex solver. 

This package is intended to be used in conjunction with the separate LPO package 
which provides tools for modelling and solving Linear Programming (LP) and Mixed-Integer 
Linear Programming (MILP) problems. Package gpx assumes that the Cplex header files
and object files as well as a compatible C compiler have been installed and configured
separately on the computer where this package is to be used.

# Installation and Configuration

The gpx package requires the errors package and Cplex to be installed. To install the packages and gpxrun executable 
on a Windows platform, go to the cmd.exe window and enter the command:
```
  go get -u github.com/pkg/errors
  go get -u github.com/go-opt/gpx
```

Download, install, and configure Cplex (https://www.ibm.com/ca-en/marketplace/ibm-ilog-cplex), and confirm
that it is working correctly. 

You will then need to modify the current placeholders for the Cplex object file (or dll) and header file with 
the correct location on your computer. The two lines which must be changed in gpx.go are:

```
#cgo LDFLAGS: -LD:/pk_cplex/include -lcplex1271
#include <D:/pk_cplex/include/ilcplex/cplex.h>
```
The gpx package was developed with Cplex version 12.7.1. It may not be compatible with earlier versions
of Cplex.

## Tips for Configuring Cplex

The installation and configuration of Cplex is outside the scope of this package, and the appropriate instructions
accompanying Cplex must be followed. During development and testing with Cplex version 12.7.1 on a Windows 7
PC, the following "tweaks" were needed in order for the callable C libraries to be compiled and linked correctly.

In file cpxconst.h, comment out reference to CPXSIZE_BITS_TEST_DISABLE:
```
  #ifndef CPXSIZE_BITS_TEST_DISABLE
  //typedef int CPXSIZE_BITS_TEST1[1 + (int)sizeof(CPXSIZE) - (int)sizeof(size_t)];
  typedef int CPXSIZE_BITS_TEST2[1 + (int)sizeof(size_t) - (int)sizeof(CPXSIZE)];
  #endif /* !CPXSIZE_BITS_TEST_DISABLE */
```

In file cplex.h, comment out all references to CPXDEPRECATEDAPI
```
/*
CPXDEPRECATEDAPI(12700100)
int CPXPUBLIC
   CPXfreeparenv (CPXENVptr env, CPXENVptr *child_p);
*/
```
Please note that the "tweaks" used to get Cplex to compile while developing gpx may not be the "correct way",
and they may differ from tweaks you might need to perform on your computer and with your version of Cplex.

# Executable gpxrun

The subdirectory gpxrun contains the executable which illustrates the functionality of the
gpx package and how it can be used to solve LP and MILP problems via Cplex.

The gpxrun directory has several text files which contain sample LP and MILP problems that are used
by the gpxrun executable. The sample files are hard-coded in the executable to minimize user input and keep the demo
program as simple as possible. Users who may be interested in a more comprehensive and complex test tool may look
https://github.com/go-opt/runopt, which was used during development and testing.

# Development and Testing

The lpo and gpx packages were developed and tested using the following software

* Windows 7 Service Pack 1
* golang, version 1.8.3
* LiteIDE X32.2, version 5.6.2 (32-bit)
* MinGW-w64 64/32-bit C compiler, version 5.1.0
* Cplex, version 12.7.1
* Coin-OR, version 1.7.4-win32-msvc11

Testing was performed using LP problems provided by netlib (http://www.netlib.org/lp/data/index.html) and some (not all)
MILP problems provided by miplib (http://miplib.zib.de/).
