## Generate Go struct from SPSS Sav file

Given an input .sav file. an output file, package name and struct name generates a Go file containing a struct than can 
be used by the dataset import function.

Options:
``````
  -input string
        input file name
  -output string
        output file name
  -package string
        package name (default "lfs")
  -struct string
        structure name (default "SpssDataItem")
