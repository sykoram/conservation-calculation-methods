# Conservation Calculation Methods

Implementation of several known methods used for calculation of conservation score of protein residues from MSA, and evaluation of their impact on [P2Rank] pocket prediction.

This repo contains a program [conservscore] written in Go that can calculate conservation scores using a specified method, and supplementary Bash scripts for easier usage and evaluation.

This project is for [Matfyz Summer of Code 2020 (cz)](https://d3s.mff.cuni.cz/msoc/) organized by [Faculty of Mathematics and Physics, Charles University, Czech Republic][MFF].


## Table of Contents

- [Background](#background)
- [How it works](#how-it-works)
- [Setup](#setup)
- [Usage](#usage)
- [Evaluation](#evaluation)
- [Results](#results)
- [Sources](#sources)


## Background

The function of proteins depends on their interactions with other molecules. Really important is an interaction between proteins and small molecules (=ligands). Lots of today's drugs are small molecules, which inhibit the function of some protein by binding to its active site.

Detection of these active sites is experimentally difficult, so there are various methods for active sites prediction. One of the best ones is [P2Rank] developed at [MFF UK][MFF]. It is known that residues of active sites are more evolutionary conserved than others. This feature enhanced the [P2Rank] prediction, but only one conservation calculation method was used.

The main goal of this project is to implement various known conservation calculation methods, and evaluate their impact on P2Rank prediction.

[COMMENT]: # (Rewritten: Funkce proteinů je odvozena od jejich interakce s ostatními molekulami. Velice důležitý typ vazby je mezi proteiny a malými molekulami, tzv. ligandy. Např. naprostá většina současných léčiv jsou právě malé molekuly, které inhibují funkci některého proteinu tím, že se váží do jeho aktivního místa a zabraňují tak šíření informace. Detekce těchto aktivních míst je ovšem experimentálně velice náročná a proto existují počítačové metody pro predikci aktivních míst. Jedna z nejlepších metod pro predikci protein-ligand aktivních míst z proteinové struktury, pojmenovaná P2Rank, byla vyvinuta na MFF UK. Je známo, že aminokyseliny aktivních míst proteinu jsou evolučně konzervovaná více než ostatní aminokyseliny a proto byla do P2Ranku přidána možnost měření evoluční konzervovanosti a využití této informace v rámci predikce. Tento přístup vskutku vedl ke zlepšení schopnosti predikce, nicméně byla použita pouze jedna metoda výpočtu konzervovanosti. Cílem projektu je tedy implementovat různé známě přístupy k výpočtu evoluční konzervovanosti a evaluovat jejich vliv na predikční schopnosti metody P2Rank. Stávající verze algoritmu je přístupná i jako webový portál na adrese [www.prankweb.cz].)


## How it works

Each [P2Rank dataset](https://github.com/rdk/p2rank-datasets) should have files containing conservation scores in `<dataset>/conservation/e5i1/scores`. The files look like this:

```
# /tmp/file4i81hc -- js_divergence - window_size: 3 - background: blosum62 - seq. weighting: True - gap penalty: 1 - normalized: False
# align_column_number	score	column

0	-1000.00000	-------------------M----M--M-------------M-------------M--
1	0.57284	-M-----------MMMMMMAMMMMVMMTMMMMMMMMMMMM-TMMMMMMMMMMMMMPMM
2	0.52533	-TMMMMM--MMMMQQQQQQQLLFFQYLFFLLLLFYLLLLL-LLLLLLLLLLLLLLLLL
3	0.62155	DVSSSSSSSSSSSDDDDDDDDDDDDDDDDDDDDDDDDDDD-DDDDDDDDDDDDDDDDD
4	0.54686	ASIIIIVIIIILFAAAAAAAAAIAAAAAVAAAAAAAAAAA-VAAAAAAAAAAAAAAAA
```

In this project, these files are used as input files for [conservscore]. This program extracts the MSA columns, calculates their conservation scores using a specified method and window size (both are passed as arguments to conservscore), and saves the scores to an output file that looks similar to the input file:

```
0	0.00000	-------------------M----M--M-------------M-------------M--
1	0.51554	-M-----------MMMMMMAMMMMVMMTMMMMMMMMMMMM-TMMMMMMMMMMMMMPMM
2	0.46502	-TMMMMM--MMMMQQQQQQQLLFFQYLFFLLLLFYLLLLL-LLLLLLLLLLLLLLLLL
3	0.71309	DVSSSSSSSSSSSDDDDDDDDDDDDDDDDDDDDDDDDDDD-DDDDDDDDDDDDDDDDD
4	0.59215	ASIIIIVIIIILFAAAAAAAAAIAAAAAVAAAAAAAAAAA-VAAAAAAAAAAAAAAAA
```

Next, [P2Rank] loads these files together with its regular input files, and they are used for training of new models and their evaluation.


## Setup

You should have [P2Rank] installed and relevant [datasets](https://github.com/rdk/p2rank-datasets) downloaded (chen11 and coach420 are used by default).

[Go](https://golang.org/) has to be installed in order to build the [conservscore] program.

Since the scripts are in Bash, it should run on Linux without a problem. On Windows, you may need to either [build conservscore](conservscore/README.md#setup) and execute commands manually, or have, for example, Git Bash or Windows Subsystem for Linux installed in order to run the Bash scripts.

You might have to change some paths of commands and directories in p2rank-eval-method script.


## Usage

The p2rank-eval-method script helps with the whole process:
1. It prepares the files with conservation score: builds conservscore, recalculates them using the specified method and collects them.
2. It runs [P2Rank] that trains and evaluates 10 new models.

```sh
./p2rank-eval-method.sh -m METHOD [-w WINDOW]
```

Available methods: 
- Shannon entropy of residues: `shannon-entropy`
- Shannon entropy of residue properties: `property-entropy`
- relative entropy (Kullback–Leibler divergence): `relative-entropy`
- Jensen-Shannon divergence: `jensen-shannon-divergence`
- sum-of-pair measure: `sum-of-pairs`

If the window is greater than 0, a score of a column is affected by nearby column scores. (the value is a number of residues on either side included in the window)

By default, the window size is 0; dataset chen11 is used for training and coach420(mlig) for evaluation. You can change some options in the scripts or in the conservscore program.

[Usage for conservscore program](conservscore/README.md#usage)

Tip: to measure the duration of the execution of a command/script/program and to print its output to both the terminal and a log file, use:

```sh
time SOME_COMMAND 2>&1 | tee log/test.log
```


## Evaluation

To evaluate the impact of different methods (or configurations) on P2Rank prediction, the following tests were made:
1. without conservation files
2. with default conservation files
3. with custom conservation files recalculated using a specified method and the window size of 0 and 3. \
   methods:
   - [Shannon entropy][SE] of residues: The most simple and commonly used method.
   - [Shannon entropy][SE] of residue properties: Amino acids are grouped by biochemical similarities.
   - [relative entropy (Kullback–Leibler divergence)][RE]: Takes into account background distribution (how rare an amino acid is).
   - [Jensen-Shannon divergence][JSD]: Based on the relative entropy, symmetric and has finite value.
   - sum-of-pairs measure: Uses a similarity matrix (takes into account how similar amino acids are).

This configuration was used:
- [P2Rank]:
  - config file: config/conservation (if with conservation files)
  - traineval training dataset: chen11.ds
  - traineval evaluation dataset: coach420(mlig).ds
  - loop: 10 (trains 10 models, each with a different seed)
  - seed: 42 (first seed)
  - rf_trees: 200
- [conservscore]: (only for the custom conservation files)
  - pseudocount: 10e-6 (0.000001)
  - use sequence weights: true
  - max. gap percentage: 30%
  - use gap penalty: true
  - similarity matrix and background distribution: BLOSUM62
  - window lam: 0.5
  - replace negative scores with 0: true

The only difference between the configurations of default and custom conservation files is that the default conservation files were probably calculated using Jensen-Shannon divergence, window size of 3 and pseudocount of 10e-7.


## Results

The following table shows the results of the tests. 

DCA(4.0) metric is used: the pocket prediction is considered successful if the distance between the pocket center and any ligand atom is less then (or equal to) 4 Å. \
[0] denotes Top-(n+0) and [2] Top-(n+2) success rate where n is the number of ligands of a protein. In order to get Top-(n+0) and Top-(n+2) success rate of 100%, the predicted pockets have to be within the first n and n+2 positions, respectively. \
These numbers are always an average score of 10 runs (each with a different seed), and were produced by P2Rank.

Shannon entropy means Shannon entropy of residues, and Property entropy means Shannon entropy of residue properties. \
By default, the window size is 0 (no window), the w3 denotes window size of 3 residues on each side.

Imp. is the improvement over using the default conservation files.

| Model / Method                 | DCA(4.0) [0] | Imp.      | DCA(4.0) [2] | Imp.      |
| ------------------------------ | ------------ | --------- | ------------ | --------- |
| **Without conservation**       | **68.3**     |           | **72.6**     |           |
| **Default conservation files** | **72.0**     | **0,0%**  | **75.4**     | **0,0%**  |
| Shannon entropy                | 72.2         | +0,3%     | 76.5         | +1,5%     |
| Shannon entropy (w3)           | 71.6         | -0,6%     | 75.1         | -0,4%     |
| property entropy               | 72.7         | +1,0%     | 76.3         | +1,2%     |
| property entropy (w3)          | 71.9         | -0,1%     | 75.3         | -0,1%     |
| **relative entropy**           | **73.7**     | **+2,4%** | **77.8**     | **+3,2%** |
| relative entropy (w3)          | 72.0         | 0,0%      | 76.0         | +0,8%     |
| **Jensen-Shannon divergence**  | **73.4**     | **+1,9%** | **77.5**     | **+2,8%** |
| Jensen-Shannon divergence (w3) | 72.2         | +0,3%     | 75.6         | +0,3%     |
| sum-of-pairs measure           | 72.4         | +0,6%     | 76.4         | +1,3%     |
| sum-of-pairs measure (w3)      | 71.9         | -0,1%     | 75.2         | -0,3%     |

The results show that P2Rank prediction is generally improved by including the conservation scores.

Both relative entropy and Jensen-Shannon divergence performed better than other tested methods, but most importantly better than using the default conservation files. Relative entropy has a slightly better score than Jensen-Shannon divergence.

Using the window definitely lowers the score (at least window size of 3). In my opinion, the reason might be that this window is sequential, however P2Rank uses its own spacial window that is more important for the pocket prediction.

**In conclusion, relative entropy (or Jensen-Shannon divergence) without a window can be used to improve P2Rank prediction.**


## Sources

Capra JA and Singh M. \
[Predicting functionally important residues from sequence conservation](https://doi.org/10.1093/bioinformatics/btm270). \
Bioinformatics. 23(15): 1875-1882, 2007.

Radoslav Krivák and David Hoksza. \
[P2Rank: machine learning based tool for rapid and accurate prediction of ligand binding sites from protein structure](https://doi.org/10.1186/s13321-018-0285-8). \
Journal of Cheminformatics. Aug 2018

Lukáš Jendele and Radoslav Krivák and Petr Škoda and Marian Novotný and David Hoksza. \
[PrankWeb: a web server for ligand binding site prediction and visualization](https://doi.org/10.1093/nar/gkz424). \
Nucleic Acids Research. May 2019





[P2Rank]: https://github.com/rdk/p2rank
[MFF]: https://www.mff.cuni.cz/en
[conservscore]: ./conservscore

[SE]: https://en.wikipedia.org/wiki/Entropy_(information_theory)
[RE]: https://en.wikipedia.org/wiki/Kullback%E2%80%93Leibler_divergence
[JSD]: https://en.wikipedia.org/wiki/Jensen%E2%80%93Shannon_divergence
