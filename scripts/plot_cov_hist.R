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
	binwidth = as.numeric(args[9])
	width = as.numeric(args[10])
	height = as.numeric(args[11])
	res_scale = as.numeric(args[12])

	data = read_bed_cov_named(cov_path, FALSE)

	plot_cov_hist(data, out_path, width, height, res_scale, ymin, ymax, xmin, xmax, xlab, ylab, binwidth)

}

main()
