#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_cov_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	cov_path = args[1]
	cov = read_bed_cov_named(cov_path, FALSE)

	out_path = args[2]

	ymin = as.numeric(args[3])
	ymax = as.numeric(args[4])

	xlab = args[5]
	ylab = args[6]

	width = as.numeric(args[7])
	height = as.numeric(args[8])
	res = as.numeric(args[9])
	textsize = as.numeric(args[10])

	plot_cov_multi_pretty(cov, out_path, width, height, res, calc_chrom_labels_string(cov), ymin, ymax, xlab, ylab, textsize)
}

main()
