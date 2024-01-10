#!/usr/bin/env Rscript

library(ggplot2)
library(dplyr)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=FALSE, sep = "\t"))
	colnames(data) = c("Chr", "Start", "End", "Structural_difference", "Cross", "Repeat_type")
	data$Cross = factor(data$Cross, levels=c("Pure D. mel", "Hybrid"))
	levels(data$Cross) = c("Pure D. mel", "Hybrid")

	# data$Repeat_type = factor(data$Repeat_type, levels=c("Peak", "Non-Peak"))
	# levels(data$Repeat_type) = c("Peak", "Non-Peak")

	print(head(data))

	df.summary2 <- data %>%
		group_by(Repeat_type, Cross) %>%
		summarise(
			sd = sd(Structural_difference),
			Structural_difference = mean(Structural_difference)
		)

	g = ggplot(data, aes(Repeat_type, Structural_difference, fill = Cross), color = "black") +
		geom_col(data = df.summary2, position=position_dodge(.8), color = "black", width = .7) +
		theme_classic() +
		labs(x = "Repeat type", y = "Structural_difference") +
		theme(axis.text.x = element_text(angle = 90, vjust = 0.5, hjust=1)) +
		ggtitle("Rate of structural difference by repeat type") +
		scale_fill_manual(values=c("#000000", "#ffffff"), breaks = c("Pure D. mel", "Hybrid"))

		# scale_y_log10() +

	rescale = 1.5

	pdf(args[2], width=4 * rescale, height=3 * rescale)
		print(g)
	dev.off()

	g2 = g + geom_errorbar(data = df.summary2, aes(ymin = Structural_difference - sd, ymax = Structural_difference + sd), position = position_dodge(.8), color = "black", width = .2)

	pdf(args[3], width=4 * rescale, height=3 * rescale)
		print(g2)
	dev.off()

	g3 = g + scale_y_continuous(limits = c(-1000000, 1000000))

	pdf(args[4], width=4 * rescale, height=3 * rescale)
		print(g3)
	dev.off()
}

main()
