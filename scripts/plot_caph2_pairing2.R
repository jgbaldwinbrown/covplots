#!/usr/bin/env Rscript

library(ggplot2)
library(dplyr)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=FALSE))
	colnames(data) = c("Chr", "Start", "End", "Pairing_Rate", "Cross", "CapH2")
	data$Cross = factor(data$Cross, levels = c("Pure D. mel (IxA)", "Hybrid (IxW)", "Hybrid (AxW)"))
	pngpath = args[2]
	pdfpath = args[3]
	width = 8
	height = 6
	res_scale = 300
	textsize = 18

	ch2breaks = c("Peak", "Non-Peak")
	ch2labels = c("Binding site", "No binding site")

	a = ggplot(data = data) +
		geom_boxplot(aes(Cross, Pairing_Rate, fill = CapH2)) +
		labs(x = "Cross", y = "Pairing_Rate") +
		ggtitle("Pairing rate at CapH2 binding sites") +
		scale_fill_manual(
			name = "Presence of CapH2\nsite in 100kb window",
			breaks = ch2breaks,
			labels = ch2labels,
			values = grey.colors(2)
		) +
		theme_bw() +
		theme(text = element_text(size=textsize)) +
		theme(axis.text.x = element_text(angle = 22.5, hjust=1))

	pdf(pdfpath, width = width, height = height)
		print(a)
	dev.off()

	png(pngpath, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a)
	dev.off()

	a2 = a + lims(y = c(0, 0.5))

	print("hi1")
	pdfpath2 = paste(pdfpath, ".lim.pdf", sep="")
	print("hi2")
	pdf(pdfpath2, width = width, height = height)
		print(a2)
	dev.off()
	print("hi3")

	pngpath2 = paste(pngpath, ".lim.png", sep="")
	print("hi4")
	png(pngpath2, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a2)
	dev.off()
	print("hi5")

	data$treat = paste(data$Cross, data$CapH2)
	model = aov(Pairing_Rate~treat, data=data)
	print(TukeyHSD(model))
}

main()
