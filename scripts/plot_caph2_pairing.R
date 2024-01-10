#!/usr/bin/env Rscript

library(ggplot2)
library(dplyr)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=FALSE))
	crosslevels = c("Pure D. mel (IxA)", "Pure D. sim (WXW)", "Hybrid (IxW)", "Hybrid (AxW)", "Hybrid Rescue (HxW)", "Hybrid Rescue (IxL)")
	colnames(data) = c("Chr", "Start", "End", "Pairing_Rate", "Cross", "CapH2")
	data$Cross = factor(data$Cross, levels=crosslevels)
	levels(data$Cross) = crosslevels

	data$CapH2 = factor(data$CapH2, levels=c("Peak", "Non-Peak"))
	levels(data$CapH2) = c("Peak", "Non-Peak")

	df.summary2 <- data %>%
		group_by(CapH2, Cross) %>%
		summarise(
			sd = sd(Pairing_Rate),
			Pairing_Rate = mean(Pairing_Rate)
		)

	g = ggplot(data, aes(CapH2, Pairing_Rate, fill = Cross), color = "black") +
		geom_col(data = df.summary2, position=position_dodge(.8), color = "black", width = .7) +
		theme_classic() +
		labs(x = "CapH2 Peak Type", y = "Pairing_Rate") +
		ggtitle("Pairing rate within or not within 1kb of CapH2 Peaks")

	pdf(args[2], width=4, height=3)
		print(g)
	dev.off()
		# scale_fill_manual(values=c("#000000", "#888888", "#ffffff"), breaks = crosslevels)

	g2 = g + geom_errorbar(data = df.summary2, aes(ymin = Pairing_Rate - sd, ymax = Pairing_Rate + sd), position = position_dodge(.8), color = "black", width = .2)

	pdf(args[3], width=4, height=3)
		print(g2)
	dev.off()
}

main()
