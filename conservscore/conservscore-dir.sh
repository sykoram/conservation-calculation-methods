#!/usr/bin/env bash

# This script runs conservscore program for each file in a directory

#set -o errexit  # Used to exit upon error, avoiding cascading errors

#WORK_DIR="$(pwd)"
SCRIPT_DIR="${0%/*}"
CONSERVSCORE="$SCRIPT_DIR/conservscore"

# handle flags
while getopts :i:o:m: option
do
case "${option}"
in
i) IN_DIR=${OPTARG};;
o) OUT_DIR=${OPTARG};;
m) METHOD=${OPTARG};;
*) echo "Invalid flag" ; exit 1;;
esac
done

# check flags
if [[ -z $IN_DIR ]]; then
  echo "Input directory missing: -i" ; EXIT=true
fi
if [[ -z $OUT_DIR ]]; then
  echo "Output directory missing: -o" ; EXIT=true
fi
if [[ -z $METHOD ]]; then
  echo "Method missing: -m" ; EXIT=true
fi
if [[ $EXIT ]]; then
  exit 1
fi

# main loop
for IN_FILEPATH in "$IN_DIR"/*.gz; do
  FILE="${IN_FILEPATH##*/}"  # filepath -> filename (./path/to/___.pdb.seq.fasta.hom.gz  ->  ___.pdb.seq.fasta.hom.gz)
  echo "Processing $FILE..."
  $CONSERVSCORE -i "$IN_DIR/$FILE" -o "$OUT_DIR/$FILE" -m "$METHOD" # the output file will be gzipped
done

echo "Done!"
