#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_cov_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	cov_path = args[1]
	out_path = args[2]
	ymin = as.numeric(args[3])
	ymax = as.numeric(args[4])

	data = read_bed_2val_named(cov_path, FALSE)

	plot_self_vs_pair_lims(data, out_path, 4, 3, 300)
}

main()
