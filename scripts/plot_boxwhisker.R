#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_cov_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)
	print(paste("args:", args))

	cov_path = args[1]
	pdf_path = args[2]
	png_path = args[3]
	ymin = as.numeric(args[4])
	ymax = as.numeric(args[5])
	xmin = as.numeric(args[6])
	xmax = as.numeric(args[7])
	ylab = args[8]
	xlab = args[9]
	width = as.numeric(args[10])
	height = as.numeric(args[11])
	resscale = as.numeric(args[12])
	textsize = as.numeric(args[13])
	fillname = args[14]

	data = read_bed_2val_named(cov_path, FALSE)
	data = na.omit(data)
	data = data[data$VAL2 != "",]

	print(paste("debug fillname:", fillname))

	plot_box_and_whisker(data, pdf_path, png_path, width, height, resscale, ymin, ymax, xmin, xmax, xlab, ylab, textsize, fillname)
}

main()
