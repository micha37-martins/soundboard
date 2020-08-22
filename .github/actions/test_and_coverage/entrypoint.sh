#!/bin/sh

# Using the "checkout@v2" action in the worflow allows to access all files
# of the repo. This allows to access the Makefiles functions.

make test
make coverage
