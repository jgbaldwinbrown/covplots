#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_cov_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	cov_path = args[1]
	out_path = args[2]
	scalespath = args[3]

	cov = read_bed_cov_named_facetted(cov_path, FALSE)

	scales = read_scales(scalespath)

	plot_cov_multi_facetsc(cov, out_path, 20, 8, 300, calc_chrom_labels_string(cov), scales)
}

main()
