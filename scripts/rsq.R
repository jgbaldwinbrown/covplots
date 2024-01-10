#!/usr/bin/env Rscript

library(data.table)

main = function() {
	args = commandArgs(trailingOnly = TRUE)
	datapath = args[1]
	depcol = as.numeric(args[2])
	indepcol = as.numeric(args[3])
	data = as.data.frame(fread(datapath, sep="\t"))
	mdata = na.omit(data[, c(depcol, indepcol)])
	colnames(mdata) = c("dep", "indep")

	r2 = cor(mdata$dep, mdata$indep)
	print(r2)
}

main()
