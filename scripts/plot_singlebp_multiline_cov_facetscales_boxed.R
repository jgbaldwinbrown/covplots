#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_cov_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	cov_path = args[1]
	out_path = args[2]
	scalespath = args[3]
	rectpath = args[4]

	cov = read_bed_cov_named_facetted(cov_path, FALSE)

	scales = read_scales(scalespath)

	print(paste("getting rect from path:", rectpath))
	rect = bed2rect_small(rectpath)
	print("got rect")

	if (nrow(rect) > 0) {
		plot_cov_multi_facetsc_boxed(cov, out_path, 20, 8, 300, calc_chrom_labels_string(cov), scales, rect)
	} else {
		plot_cov_multi_facetsc(cov, out_path, 20, 8, 300, calc_chrom_labels_string(cov), scales)
	}
	print("plotted")
}

main()

# bed2rect <- function(path) {
# 	# bed = read_bed_noval(path)
# 	bed = read_bed_postsub(path)
# 	rect = data.frame(ymin = rep(-Inf, nrow(bed)),
# 		ymax = rep(Inf, nrow(bed)),
# 		xmin = bed$cumsum.tmp,
# 		xmax = bed$cumsum.tmp2)
# 	return(rect)
# }
# 
