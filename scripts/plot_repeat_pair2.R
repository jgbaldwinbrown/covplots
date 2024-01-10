#!/usr/bin/env Rscript

library(ggplot2)
library(dplyr)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=FALSE))
	colnames(data) = c("Chr", "Start", "End", "Pairing_Rate", "Cross", "Repeat_Type")
	pngpath = args[2]
	pdfpath = args[3]
	width = 12
	height = 9
	res_scale = 300
	textsize = 12


	a = ggplot(data = data) +
		geom_boxplot(aes(Repeat_Type, Pairing_Rate, fill = Cross)) +
		labs(x = "Repeat Type", y = "Pairing Rate") +
		ggtitle("Pairing rate at repeats") +
		scale_fill_discrete(name = "Cross") +
		theme_bw() +
		theme(text = element_text(size=textsize)) +
		theme(axis.text.x = element_text(angle = 22.5, hjust=1))

	pdf(pdfpath, width = width, height = height)
		print(a)
	dev.off()

	png(pngpath, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a)
	dev.off()

	a2 = a + lims(y=c(0, 1))

	pdfpath2 = paste(pdfpath, ".lim.pdf", sep="")
	pdf(pdfpath2, width = width, height = height)
		print(a2)
	dev.off()

	pngpath2 = paste(pngpath, ".lim.png", sep="")
	png(pngpath2, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a2)
	dev.off()

	data$treat = paste(data$Cross, data$Repeat_Type)
	model = aov(Pairing_Rate~treat, data=data)
	print(TukeyHSD(model))
}

main()
