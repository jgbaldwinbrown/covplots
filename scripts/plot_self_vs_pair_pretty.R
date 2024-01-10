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
	xmin = as.numeric(args[5])
	xmax = as.numeric(args[6])
	ylab = args[7]
	xlab = args[8]
	width = as.numeric(args[9])
	height = as.numeric(args[10])
	resscale = as.numeric(args[11])
	textsize = as.numeric(args[12])

	data = read_bed_2val_named(cov_path, FALSE)

	plot_self_vs_pair_pretty(data, out_path, width, height, resscale, ymin, ymax, xmin, xmax, ylab, xlab, textsize)
}

main()
