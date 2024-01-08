# BCell-VaccineInference

This is a collection of scripts for the identification of vaccine specific B cells from immune repertoire sequencing.
It seeks to classify B cells as either vaccine-specific, non-specific or background.

R scripts are dependent on each other; sourcing Mixture_EM.R will sequentially source all other scripts.

Required packages are:VGAM, ismev, MASS, and their dependencies.


## Input
The EM function, is the main function for model fitting, and has the following inputs:

- **Data** a 3 dimensional array of clonal abundance over samples and time points: clone x sample x time.
- **max_iter** the maximum number of iterations to perform. Algorithm will finish either when convergence or the maximum number of iterations is achieved, whichever is soonest.
- **p_bg** default =.85; the initial probability of being assigned to the background class.
- **p_ind** default=.1; the initial probability of being assigned to the vaccine-specific class.
- **w_bg** default=c(.99,.009,.001); the initial mixture weights for a clone assigned to the background class.
- **w_ind** default=c(.1,.2,.7); the initial mixture weights for a clone assigned to the vaccine-specific class.
- **Params** default=NULL; if NULL initial parameters for the distributions are calculated based on the entire data set.

## Output

Returns a list containing the following elements:

1. Assignment of a clone to the background (labelled 1), vaccine-specific (labelled 2) or non-specific (labelled 3) categories.
2. Final mixture weights for a clone assigned to the background class.
3. Final probability of being assigned to the background class.
4. Final probability of a clone assigned to the background class being present in the sample.
5. Final mixture weights for a clone assigned to the vaccine-specific class.
6. Final probability of being assigned to the vaccine-specific class.
7. Final probability of a clone assigned to the vaccine-specific class being present in the sample.
8. Log-likelihood of the final model.
9. Log-likelihood of all models evaluated.
10. Number of iterations performed.

## Example

```
load('ExampleData.RData')
source('Mixture_EM.R')
results<-EM(X,max_iter=10)

#number of sequences allocated to vaccine specific class
sum(results[[1]]==2)

#number of sequences allocated to background class
sum(results[[1]]==1)
```
