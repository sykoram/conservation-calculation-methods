/*
This file contains functions for score calculations.

Some parts of this code were implemented according to the following paper and its supplementary data:
Capra JA and Singh M.
Predicting functionally important residues from sequence conservation.
Bioinformatics. 23(15): 1875-1882, 2007.
 */

package main

import (
	"math"
	"strings"
)

type MsaColumn = string
type Aa = rune
type SimilarityMatrix = [][]int
type MethodFunc func(col MsaColumn, simMatrix SimilarityMatrix, bgDistr []float64, seqWeights []float64) (conservScore float64)

var Methods = map[string]MethodFunc{
	"zero": Zero,
	"shannon-entropy": ShannonEntropy,
	"property-entropy": ShannonPropertyEntropy,
}

var AminoAcids = []Aa{'A', 'R', 'N', 'D', 'C', 'Q', 'E', 'G', 'H', 'I', 'L', 'K', 'M', 'F', 'P', 'S', 'T', 'W', 'Y', 'V', '-'}
var Blosum62BgDistr = []float64{0.078, 0.051, 0.041, 0.052, 0.024, 0.034, 0.059, 0.083, 0.025, 0.062, 0.092, 0.056, 0.024, 0.044, 0.043, 0.059, 0.055, 0.014, 0.034, 0.072}
var AaToIndex = getAaToIndexMap()
//var UseSeqWeights = true
var Pseudocount = 0.000001  // 10e-6
//var Pseudocount = 0.0000001  // xxx 10e-7
var MaxGapRatio = .30
var UseGapPenalty = true
//var WindowSize = 3  // 0 for no window scores calculation
var WindowLam = 0.5
/*
Matrix made by matblas from blosum62.iij,
* column uses minimum score,
BLOSUM Clustered Scoring Matrix in 1/2 Bit Units,
Blocks Database = /data/blocks_5.0/blocks.dat,
Cluster Percentage: >= 62,
Entropy =   0.6979, Expected =  -0.5209
*/
var Blosum62SimMatrix = SimilarityMatrix{
	//A   R   N   D   C   Q   E   G   H   I   L   K   M   F   P   S   T   W   Y   V   B   Z   X   *
	{ 4, -1, -2, -2,  0, -1, -1,  0, -2, -1, -1, -1, -1, -2, -1,  1,  0, -3, -2,  0, -2, -1,  0, -4},
	{-1,  5,  0, -2, -3,  1,  0, -2,  0, -3, -2,  2, -1, -3, -2, -1, -1, -3, -2, -3, -1,  0, -1, -4},
	{-2,  0,  6,  1, -3,  0,  0,  0,  1, -3, -3,  0, -2, -3, -2,  1,  0, -4, -2, -3,  3,  0, -1, -4},
	{-2, -2,  1,  6, -3,  0,  2, -1, -1, -3, -4, -1, -3, -3, -1,  0, -1, -4, -3, -3,  4,  1, -1, -4},
	{ 0, -3, -3, -3,  9, -3, -4, -3, -3, -1, -1, -3, -1, -2, -3, -1, -1, -2, -2, -1, -3, -3, -2, -4},
	{-1,  1,  0,  0, -3,  5,  2, -2,  0, -3, -2,  1,  0, -3, -1,  0, -1, -2, -1, -2,  0,  3, -1, -4},
	{-1,  0,  0,  2, -4,  2,  5, -2,  0, -3, -3,  1, -2, -3, -1,  0, -1, -3, -2, -2,  1,  4, -1, -4},
	{ 0, -2,  0, -1, -3, -2, -2,  6, -2, -4, -4, -2, -3, -3, -2,  0, -2, -2, -3, -3, -1, -2, -1, -4},
	{-2,  0,  1, -1, -3,  0,  0, -2,  8, -3, -3, -1, -2, -1, -2, -1, -2, -2,  2, -3,  0,  0, -1, -4},
	{-1, -3, -3, -3, -1, -3, -3, -4, -3,  4,  2, -3,  1,  0, -3, -2, -1, -3, -1,  3, -3, -3, -1, -4},
	{-1, -2, -3, -4, -1, -2, -3, -4, -3,  2,  4, -2,  2,  0, -3, -2, -1, -2, -1,  1, -4, -3, -1, -4},
	{-1,  2,  0, -1, -3,  1,  1, -2, -1, -3, -2,  5, -1, -3, -1,  0, -1, -3, -2, -2,  0,  1, -1, -4},
	{-1, -1, -2, -3, -1,  0, -2, -3, -2,  1,  2, -1,  5,  0, -2, -1, -1, -1, -1,  1, -3, -1, -1, -4},
	{-2, -3, -3, -3, -2, -3, -3, -3, -1,  0,  0, -3,  0,  6, -4, -2, -2,  1,  3, -1, -3, -3, -1, -4},
	{-1, -2, -2, -1, -3, -1, -1, -2, -2, -3, -3, -1, -2, -4,  7, -1, -1, -4, -3, -2, -2, -1, -2, -4},
	{ 1, -1,  1,  0, -1,  0,  0,  0, -1, -2, -2,  0, -1, -2, -1,  4,  1, -3, -2, -2,  0,  0,  0, -4},
	{ 0, -1,  0, -1, -1, -1, -1, -2, -2, -1, -1, -1, -1, -2, -1,  1,  5, -2, -2,  0, -1, -1,  0, -4},
	{-3, -3, -4, -4, -2, -2, -3, -2, -2, -3, -2, -3, -1,  1, -4, -3, -2, 11,  2, -3, -4, -3, -2, -4},
	{-2, -2, -2, -3, -2, -1, -2, -3,  2, -1, -1, -2, -1,  3, -3, -2, -2,  2,  7, -1, -3, -2, -1, -4},
	{ 0, -3, -3, -3, -1, -2, -2, -3, -3,  3,  1, -2,  1, -1, -2, -2,  0, -3, -1,  4, -3, -2, -1, -4},
	{-2, -1,  3,  4, -3,  0,  1, -1,  0, -3, -4,  0, -3, -3, -2,  0, -1, -4, -3, -3,  4,  0, -1, -4},
	{-1,  0,  0,  1, -3,  3,  4, -2,  0, -3, -3,  1, -1, -3, -1,  0, -1, -3, -2, -2,  0,  4, -1, -4},
	{ 0, -1, -1, -1, -2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -2,  0,  0, -2, -1, -1, -1, -1, -1, -4},
	{-4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4, -4,  1},
}


// ----- METHODS -----

/*
Dummy method. Returns zero.
 */
func Zero(col MsaColumn, simMatrix SimilarityMatrix, bgDistr []float64, seqWeights []float64) float64 {
	return 0
}

/*
Calculates the Shannon entropy (of residues) of the column. simMatrix and bgDistr are ignored.
If UseGapPenalty is true, then gaps are penalized.
The entropy will be between zero and one because of its base. See p.13 of Valdar 02 for details.
The information score 1 - h is returned for the sake of consistency with other scores.
*/
func ShannonEntropy(col MsaColumn, simMatrix SimilarityMatrix, bgDistr []float64, seqWeights []float64) float64 {
	wfc := GetWeightedFrequencyCount(col, seqWeights, Pseudocount)

	h := 0.0
	for _, c := range wfc {
		h += c * math.Log(c)  // xxx: what base should be used? 2 / e / 10
	}
	h /= math.Log(math.Min(float64(len(wfc)), float64(len(col))))
	h *= -1

	score := 1 - h
	if UseGapPenalty {
		return score * GetWeightedGapPenalty(col, seqWeights)
	}
	return score
}

/*
Calculates the Shannon entropy (of residue properties) of a column col relative to a partition of the amino acids.
Similar to Mirny '99. sim_matrix and bg_distr are ignored.
The information score 1 - h is returned for the sake of consistency with other scores.
*/
func ShannonPropertyEntropy(col MsaColumn, simMatrix SimilarityMatrix, bgDistr []float64, seqWeights []float64) float64 {
	propertyPartition := [][]Aa{{'A','V','L','I','M','C'}, {'F','W','Y','H'}, {'S','T','N','Q'}, {'K','R'}, {'D', 'E'}, {'G', 'P'}, {'-'}}
	wfc := GetWeightedFrequencyCount(col, seqWeights, Pseudocount)

	// sum AA frequencies -> property frequencies
	propertyWfc := make([]float64, len(propertyPartition))
	for p := range propertyWfc {
		for _, aa := range propertyPartition[p] {
			propertyWfc[p] += wfc[AaToIndex[aa]]
		}
	}

	h := 0.0
	for _, c := range propertyWfc {
		h += c * math.Log(c)
	}
	h /= math.Log(math.Min(float64(len(propertyWfc)), float64(len(col))))
	h *= -1

	score := 1 - h
	if UseGapPenalty {
		return score * GetWeightedGapPenalty(col, seqWeights)
	}
	return score
}

/*
This function takes a list of scores and a length and transforms them so that each position is a weighted average of the surrounding positions.
Positions with scores less than zero are not changed and are ignored in the calculation.
Here windowLen is interpreted to mean windowLen residues on either side of the current residue.
*/
func WindowScores(scores []float64, windowLen int, lam float64) []float64 {
	wScores := make([]float64, len(scores))
	copy(wScores, scores)

	for i := windowLen; i < len(scores) - windowLen; i++ {  // xxx few first and last elements are not included
		if scores[i] < 0 {
			continue
		}

		sum := 0.0
		nTerms := 0.0
		for j := i - windowLen; j < i + windowLen + 1; j++ {
			if i != j && scores[j] > 0 {
				nTerms++
				sum += scores[j]
			}
		}
		if nTerms > 0 {
			wScores[i] = (1 - lam) * (sum / nTerms) + lam * scores[i]
		}
	}

	return wScores
}

// ----- END OF METHODS -----

/*
Calculates the sequence weights using the Henikoff '94 method.
The parameter msaCols is a slice of MSA columns (not sequences),
but the returned slice is a slice of sequence weights.
Every column must have the same length.
Therefore, the number of returned sequence weights (length of the returned slice)
equals the length of the first (or any) column.
*/
func GetSequenceWeights(msaCols []MsaColumn) []float64 {
	nSeqs := len(msaCols[0])
	nCols := len(msaCols)

	seqWeights := make([]float64, nSeqs)
	for _, col := range msaCols {
		freqCounts := make([]int, len(AminoAcids))

		for _, aa := range col {
			if aa != '-' { // ignore gaps
				freqCounts[AaToIndex[aa]]++
			}
		}

		nObservedTypes := 0
		for _, fc := range freqCounts {
			if fc > 0 {
				nObservedTypes++
			}
		}

		for i, aa := range col {
			d := freqCounts[AaToIndex[aa]] * nObservedTypes
			if d > 0 {
				seqWeights[i] += 1.0 / float64(d)
			}
		}
	}

	for i := range seqWeights {
		seqWeights[i] /= float64(nCols)
	}

	return seqWeights
}

/*
Returns gaps in column divided by length of column. 20% gaps -> 0.2
*/
func GetGapRatio(c MsaColumn) float64 {
	nGaps := strings.Count(c, "-")
	return float64(nGaps) / float64(len(c))
}

/*
Calculates the simple gap penalty multiplier for the column.
If the sequences are weighted, the gaps, when penalized, are weighted accordingly.
*/
func GetWeightedGapPenalty(c MsaColumn, seqWeights []float64) float64 {
	gapsSum := 0.0
	for i, aa := range c {
		if aa == '-' {
			gapsSum += seqWeights[i]
		}
	}
	return 1 - (gapsSum / sumFloat64s(seqWeights))
}

/*
Returns the weighted frequency count (with pseudocount) of every amino acid for a column.
*/
func GetWeightedFrequencyCount(c MsaColumn, seqWeights []float64, pseudocount float64) []float64 {
	freqCounts := make([]float64, len(AminoAcids))
	for i := range freqCounts {
		freqCounts[i] = pseudocount
	}

	for i, aa := range AminoAcids {
		for j, caa := range c {
			if caa == aa {
				freqCounts[i] += seqWeights[j]
			}
		}
	}

	swSum := sumFloat64s(seqWeights)
	for i := range freqCounts {
		freqCounts[i] = freqCounts[i] / (swSum + float64(len(AminoAcids)) * pseudocount)
	}

	return freqCounts
}

/*
Returns the sum of the elements of the float64 slice.
 */
func sumFloat64s(slice []float64) float64 {
	sum := 0.0
	for _, f := range slice {
		sum += f
	}
	return sum
}

/*
Returns a map: keys are amino acids and values are their indexes.
 */
func getAaToIndexMap() map[Aa]int {
	m := make(map[Aa]int)
	for i, aa := range AminoAcids {
		m[aa] = i
	}
	return m
}

/*
Returns identifiers of the Methods.
 */
func GetMethodNames() []string {
	keys := make([]string, len(Methods))
	i := 0
	for k := range Methods {
		keys[i] = k
		i++
	}
	return keys
}