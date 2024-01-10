#!/usr/bin/env Rscript

library(data.table)

main = function() {
	args = commandArgs(trailingOnly = TRUE)
	datapath = args[1]
	depcol = as.numeric(args[2])
	indepcol = as.numeric(args[3])
	data = as.data.frame(fread(datapath, sep="\t"))
	mdata = data[, c(depcol, indepcol)]
	colnames(mdata) = c("dep", "indep")

	fit = lm(dep ~ indep, data = mdata)

	af <- anova(fit)
	afss <- af$"Sum Sq"
	print(cbind(af,PctExp=afss/sum(afss)*100))
}

main()
