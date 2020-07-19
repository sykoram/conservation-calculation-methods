# Conservation Calculation Methods

Implementation of several known methods of calculation the conservation score of protein residues from MSA, and evaluation of their impact on [P2Rank] pocket prediction.

This project is from [Matfyz Summer of Code 2020 (cz)](https://d3s.mff.cuni.cz/msoc/) organized by [Faculty of Mathematics and Physics, Charles University, Czech Republic (en)][MFF].



## Background

The function of proteins depends on their interactions with other molecules. Really important is an interaction between proteins and small molecules (=ligands). Lots of today's drugs are small molecules, which inhibit the function of some protein by binding to its active site.

Detection of these active sites is experimentally difficult, so there are various methods for active sites prediction. One of the best ones is [P2Rank] developed at [MFF UK][MFF]. It is known that residues of active sites are more evolutionary conserved than others. This feature enhanced the [P2Rank] prediction, but only one conservation calculation method was used.

The main goal of this project is to implement various known conservation calculation methods, and evaluate their impact on P2Rank prediction.




[P2Rank]: https://github.com/rdk/p2rank
[MFF]: https://www.mff.cuni.cz/en