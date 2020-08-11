#!/usr/bin/env bash

# This script helps to evaluate how the specified conservation-calculation method affects results of P2Rank predictions.
# 1. prepare conservation-score files - recalculate them using the specified method
# 2. collect files of the two datasets used in 'prank traineval' command
# 3. train and evaluate new model(s) (chen11 is used for training, coach420(mlig) for evaluation, loop=10)

set -o errexit  # Used to exit upon error, avoiding cascading errors

TRAIN_DS="chen11"
EVAL_DS="coach420"  # specify without '(mlig)', it will be added later
RUNS=10  # each with different seed

# handle flags
METHOD=""
WINDOW=0
while getopts :m:w: option
do
  case "${option}"
    in
    m) METHOD=${OPTARG};;
    w) WINDOW=${OPTARG};;
    *) echo "Invalid flag" ; exit 1;;
  esac
done
if [[ -z $METHOD ]]; then
  echo "Method missing: -m" ; exit 1
fi


# print info
echo "TRAIN DATASET: $TRAIN_DS"
echo "EVAL DATASET: $EVAL_DS"
echo "METHOD: $METHOD (WINDOW: $WINDOW)"
echo "$RUNS RUNS(S)"


# define paths
SCRIPT_DIR="${0%/*}"
CONSERVSCORE="$SCRIPT_DIR/conservscore/conservscore"
CONSERVSCORE_DIR_SCRIPT="$SCRIPT_DIR/conservscore/conservscore-dir.sh"
P2RANK_DIR="$SCRIPT_DIR/../p2rank_2.1"

DATASETS="$P2RANK_DIR/datasets"
DS_SCORES_DIR="conservation/e5i1/scores" # relative to the dataset
ROOT_CONSERVATION_DIR="$P2RANK_DIR/conservation"


echo;echo "##############################################"
echo "  PREPARING FILES WITH CONSERVATION SCORE..."
echo "##############################################"

echo "BUILDING conservscore..."
(cd "${CONSERVSCORE%/*}" && go build -o conservscore ./conservscore.go ./methods.go || exit 1)

for DS in "$TRAIN_DS" "$EVAL_DS"; do
  echo "PROCESSING DATASET $DS"
  IN_DIR="$DATASETS/$DS/$DS_SCORES_DIR"
  DS_METHOD_DIR="$ROOT_CONSERVATION_DIR/$DS/$METHOD-w$WINDOW"
  $CONSERVSCORE_DIR_SCRIPT -i "$IN_DIR" -o "$DS_METHOD_DIR" -m "$METHOD" -w "$WINDOW"
done


echo;echo "########################"
echo "  COLLECTING FILES..."     # copy files of the two datasets used for traineval into one directory
echo "########################"
CONSERVDIR="$ROOT_CONSERVATION_DIR/${TRAIN_DS}_${EVAL_DS}/$METHOD-w$WINDOW"
mkdir -p "$CONSERVDIR"
cp -R "$ROOT_CONSERVATION_DIR/$TRAIN_DS/$METHOD-w$WINDOW/." "$CONSERVDIR"
cp -R "$ROOT_CONSERVATION_DIR/$EVAL_DS/$METHOD-w$WINDOW/." "$CONSERVDIR"


echo;echo "##############################"
echo "  TRAINING AND EVALUATING..."
echo "##############################"
cd "$P2RANK_DIR" # to prevent problems with P2Rank
OUT_SUBDIR="CONSERV-CALC-METHODS"
./prank traineval \
  -t "./datasets/$TRAIN_DS.ds" \
  -e "./datasets/$EVAL_DS(mlig).ds" \
  -threads 8 -rf_trees 200 -delete_models 0 -loop $RUNS -seed 42 -c "./config/conservation" \
  -fail_fast 1 \
  -label "conservation_$METHOD-w$WINDOW" \
  -conservation_dir "../conservation/${TRAIN_DS}_${EVAL_DS}/$METHOD-w$WINDOW" \
  -out_subdir "$OUT_SUBDIR"
  # The conservation_dir path is relative to the dataset definition file.

echo "Done!"