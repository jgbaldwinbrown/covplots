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

	print("reading")
	cov = read_bed_cov_named_facetted(cov_path, FALSE)
	print("finished reading")

	print("plotting")
	plot_cov_multi_facet(cov, out_path, 20, 8, 300, calc_chrom_labels_string(cov), ymin, ymax)
	print("finished plotting")
}

main()
