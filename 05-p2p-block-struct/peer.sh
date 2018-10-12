#!/bin/bash

option="local"

if [ $# -eq 1 ]
  then option="$1"
fi

go run *.go -v -b $option