# Conservation Calculation Methods

Implementation of several known methods used for calculation of conservation score of protein residues from MSA, and evaluation of their impact on [P2Rank] pocket prediction.

This repo contains a program [conservscore] written in Go that can calculate conservation scores using a specified method, and supplementary Bash scripts for easier usage and evaluation.

This project is for [Matfyz Summer of Code 2020 (cz)](https://d3s.mff.cuni.cz/msoc/) organized by [Faculty of Mathematics and Physics, Charles University, Czech Republic][MFF].


## Table of Contents

- [Background](#background)
- [Setup](#setup)
- [Usage](#usage)
- [Evaluation](#evaluation)
- [Results](#results)
## Background

The function of proteins depends on their interactions with other molecules. Really important is an interaction between proteins and small molecules (=ligands). Lots of today's drugs are small molecules, which inhibit the function of some protein by binding to its active site.

Detection of these active sites is experimentally difficult, so there are various methods for active sites prediction. One of the best ones is [P2Rank] developed at [MFF UK][MFF]. It is known that residues of active sites are more evolutionary conserved than others. This feature enhanced the [P2Rank] prediction, but only one conservation calculation method was used.

The main goal of this project is to implement various known conservation calculation methods, and evaluate their impact on P2Rank prediction.

[COMMENT]: # (Rewritten: Funkce proteinů je odvozena od jejich interakce s ostatními molekulami. Velice důležitý typ vazby je mezi proteiny a malými molekulami, tzv. ligandy. Např. naprostá většina současných léčiv jsou právě malé molekuly, které inhibují funkci některého proteinu tím, že se váží do jeho aktivního místa a zabraňují tak šíření informace. Detekce těchto aktivních míst je ovšem experimentálně velice náročná a proto existují počítačové metody pro predikci aktivních míst. Jedna z nejlepších metod pro predikci protein-ligand aktivních míst z proteinové struktury, pojmenovaná P2Rank, byla vyvinuta na MFF UK. Je známo, že aminokyseliny aktivních míst proteinu jsou evolučně konzervovaná více než ostatní aminokyseliny a proto byla do P2Ranku přidána možnost měření evoluční konzervovanosti a využití této informace v rámci predikce. Tento přístup vskutku vedl ke zlepšení schopnosti predikce, nicméně byla použita pouze jedna metoda výpočtu konzervovanosti. Cílem projektu je tedy implementovat různé známě přístupy k výpočtu evoluční konzervovanosti a evaluovat jejich vliv na predikční schopnosti metody P2Rank. Stávající verze algoritmu je přístupná i jako webový portál na adrese [www.prankweb.cz].)


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
time YOUR_COMMAND 2>&1 | tee log/test.log
```


## Evaluation

To evaluate the impact of different methods (or configurations) on P2Rank prediction, the following tests were made:
1. without conservation files
2. with default conservation files
3. with custom conservation files recalculated using a specified method and the window size of 0 and 3. \
   methods:
   - Shannon entropy of residues
   - Shannon entropy of residue properties
   - relative entropy (Kullback–Leibler divergence)
   - Jensen-Shannon divergence
   - sum-of-pair measure

This configuration was used:
- P2Rank:
  - config file: config/conservation (if with conservation files)
  - traineval training dataset: chen11.ds
  - traineval evaluation dataset: coach420(mlig).ds
  - loop: 10 (trains 10 models, each with a different seed)
  - seed: 42 (first seed)
  - rf_trees: 200
- conservscore: (only for the custom conservation files)
  - pseudocount: 10e-6 (0.000001)
  - use sequence weights: true
  - max. gap percentage: 30%
  - use gap penalty: true
  - similarity matrix and background distribution: BLOSUM62
  - window lam: 0.5
  - replace negative scores with 0: true

The only difference between the configurations of default and custom conservation files is that the default conservation files were probably calculated using Jensen-Shannon divergence, window size of 3 and pseudocount of 10e-7.






[P2Rank]: https://github.com/rdk/p2rank
[MFF]: https://www.mff.cuni.cz/en
[conservscore]: ./conservscore