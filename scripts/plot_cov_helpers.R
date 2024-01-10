#!/usr/bin/env Rscript

library(dplyr)
library(data.table)
library(magrittr)
library(ggplot2)
library(facetscales)

read_scales <- function(inpath) {
	print("reading scales from:")
	print(inpath)
	giant = as.data.frame(fread(inpath, header=TRUE))
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("FACET", "MIN", "MAX")
	out = apply(giant[,2:3], 1, function(x){scale_y_continuous(lim=c(x[1], x[2]))})
	names(out) = giant$FACET
	return(out)
}

read_combined_pvals_precomputed <- function(inpath) {
	giant = as.data.frame(fread(inpath), header=TRUE)
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("chrom", "BP1", "BP", "PFST", "CHISQ", "WINDOW_P", "THRESH", "WINDOW_FDR_P", "WINDOW_FDR_NLOGP", "BONF_THRESH", "CHR", "cumsum.tmp")
	giant$VAL = -log10(giant$WINDOW_FDR_P)
	return(giant)
}

read_combined_pvals <- function(inpath) {
	giant = as.data.frame(fread(inpath), header=TRUE)
	colnames(giant) = c("chrom", "BP1", "BP", "PFST", "CHISQ", "WINDOW_P", "CHR", "cumsum.tmp")
	giant$VAL = -log10(giant$WINDOW_P)
	return(giant)
}

read_pvals_nowin <- function(inpath) {
	giant = as.data.frame(fread(inpath), header=TRUE)
	colnames(giant) = c("chrom", "BP", "PFST", "CHR", "cumsum.tmp")
	giant$VAL = -log10(giant$PFST)
	return(giant)
}

# FST win is in bed format
read_bed <- function(inpath, nlog) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	colnames(giant) = c("chrom", "BP1", "BP", "VAL", "CHR", "cumsum.tmp")
	if (nlog) {
		giant$VAL = -log10(giant$VAL)
	}
	return(giant)
}

# FST win is in bed format
read_bed_cov <- function(inpath, nlog) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("chrom", "BP1", "BP", "VAL", "CHR", "cumsum.tmp", "cumsum.tmp2")
	if (nlog) {
		giant$VAL = -log10(giant$VAL)
	}
	return(giant)
}

# FST win is in bed format
read_bed_cov_named <- function(inpath, nlog) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			numeric(),
			character(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("chrom", "BP1", "BP", "VAL", "NAME", "CHR", "cumsum.tmp", "cumsum.tmp2")
	if (nlog) {
		giant$VAL = -log10(giant$VAL)
	}
	return(giant)
}

# FST win is in bed format with extra FACET column after NAME
read_bed_cov_named_facetted <- function(inpath, nlog) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	print(head(giant))
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			numeric(),
			character(),
			character(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	print("giant could be empty")
	print(head(giant))
	colnames(giant) = c("chrom", "BP1", "BP", "VAL", "FACET", "NAME", "CHR", "cumsum.tmp", "cumsum.tmp2")
	if (nlog) {
		giant$VAL = -log10(giant$VAL)
	}
	print("final giant")
	print(head(giant))
	return(giant)
}

# FST win is in bed format with extra FACET column after NAME
read_bed_cov_named_labelled <- function(inpath, nlog) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	print(head(giant))
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			numeric(),
			character(),
			character(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	print("giant could be empty")
	print(head(giant))
	colnames(giant) = c("chrom", "BP1", "BP", "VAL", "LABEL", "NAME", "CHR", "cumsum.tmp", "cumsum.tmp2")
	if (nlog) {
		giant$VAL = -log10(giant$VAL)
	}
	print("final giant")
	print(head(giant))
	return(giant)
}

# FST win is in bed format
read_bed_2val_named <- function(inpath, nlog) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	print(head(giant))
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			numeric(),
			character(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("chrom", "BP1", "BP", "VAL1", "VAL2", "NAME", "CHR", "cumsum.tmp", "cumsum.tmp2")
	if (nlog) {
		giant$VAL = -log10(giant$VAL)
	}
	return(giant)
}

read_2_col <- function(inpath, nlog) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	if (ncol(giant) == 0) {
		giant = data.frame(
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("VAL1", "VAL2")
	if (nlog) {
		giant$VAL = -log10(giant$VAL)
	}
	return(giant)
}

# for reading bed format files with only the minimum columns
read_bed_noval <- function(inpath) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("chrom", "BP1", "BP", "CHR", "cumsum.tmp", "cumsum.tmp2")
	return(giant)
}

# for reading bed format files with only the minimum columns
read_bed_postsub <- function(inpath) {
	giant = as.data.frame(fread(inpath), header=FALSE)
	if (ncol(giant) == 0) {
		giant = data.frame(
			character(),
			numeric(),
			numeric(),
			character(),
			numeric(),
			numeric(),
			character(),
			numeric(),
			numeric(),
			numeric(),
			stringsAsFactors = FALSE
		)
	}
	colnames(giant) = c("chrom", "BP1", "BP", "dot", "len", "chrlen", "name", "CHR", "cumsum.tmp", "cumsum.tmp2")
	return(giant)
}

read_fst <- function(inpath) {
	giant = as.data.frame(fread(inpath), header=TRUE)
	colnames(giant) = c("chrom", "BP", "FST", "CHR", "cumsum.tmp")
	giant$VAL = giant$FST
	return(giant)
}

read_selec <- function(inpath) {
	giant = as.data.frame(fread(inpath), header=TRUE)
	colnames(giant) = c("chrom", "BP", "S", "S_P", "CHR", "cumsum.tmp")
	giant$VAL = giant$S
	return(giant)
}

calc_chrom_labels <- function(giant) {
	medians <- giant %>% dplyr::group_by(CHR) %>% dplyr::summarise(median.x = median(cumsum.tmp))
}

calc_chrom_labels_string <- function(giant) {
	medians <- giant %>% dplyr::group_by(chrom) %>% dplyr::summarise(median.x = median(cumsum.tmp))
}

calc_thresh <- function(data, colname, thresh, na.rm) {
	# usually use .9999 as threshold
	return(quantile(data[,colname], thresh, na.rm=na.rm))
}

pass_thresh <- function(data, colname, thresh) {
	data[,colname] > thresh
}

threshcolor <- function(data, chrcol, passcol) {
	out = factor(((data[,chrcol] %% 2) * (1-data[,passcol])) + (3 * data[,passcol]))
	return(out)
	#return(factor(((data[,chrcol] %% 2) * (1-data[,passcol])) + (3 * data[,passcol])))
}

nothreshcolor <- function(data, chrcol) {
	return(factor(data[,chrcol] %% 2))
}

join <- function(datas, thresholds, names, threshold_names) {
	for (i in 1:length(names)) {
		if (nrow(datas[[i]]) > 0) {
			datas[[i]]$NAME = names[i]
		}
	}
	outdata = as.data.frame(
		do.call("rbind",
			lapply(datas, function(x) {
				if (nrow(x) > 0) {
					return(x[,c("CHR", "BP", "cumsum.tmp", "VAL", "NAME", "color")])
				}
				
				return(data.frame(
					numeric(),
					numeric(),
					numeric(),
					character(),
					character(),
					stringsAsFactors = FALSE
				))
			})
		)
	)

	for (i in 1:length(threshold_names)) {
		thresholds[[i]]$NAME = threshold_names[i]
	}
	outthresholds = as.data.frame(do.call("rbind", thresholds))
	return(list(outdata, outthresholds))
}

plot <- function(data, valcol, path, width, height, res_scale, thresholds, medians) {
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = cumsum.tmp, y = VAL, color = color)) +
		geom_hline(data = thresholds, aes(yintercept = THRESH), linetype="dashed") +
		scale_x_continuous(breaks = medians$median.x, labels = medians$CHR) +
		guides(colour=FALSE) +
		xlab("Chromosome") +
		ylab(expression(-log[10](italic(p)))) +
		scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"))+
		theme_bw() +
		facet_grid(NAME~., scales="free_y") +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
}

plot_scaled_y <- function(data, valcol, path, width, height, res_scale, thresholds, medians, scales_y) {
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = cumsum.tmp, y = VAL, color = color)) +
		geom_hline(data = thresholds, aes(yintercept = THRESH), linetype="dashed") +
		scale_x_continuous(breaks = medians$median.x, labels = medians$CHR) +
		guides(colour=FALSE) +
		xlab("Chromosome") +
		ylab(expression(-log[10](italic(p)))) +
		scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"))+
		theme_bw() +
		facet_grid_sc(NAME~., scales=list(y=scales_y)) +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
}

plot_scaled_y_boxed <- function(data, valcol, path, width, height, res_scale, thresholds, medians, scales_y, rect) {
	print(rect)
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_rect(data = rect, aes(xmin = xmin, xmax = xmax, ymin = ymin, ymax = ymax), fill = "#5555DD", color = "#5555DD", alpha = 0.3) +
		geom_point(aes(x = cumsum.tmp, y = VAL, color = color)) +
		geom_hline(data = thresholds, aes(yintercept = THRESH), linetype="dashed") +
		scale_x_continuous(breaks = medians$median.x, labels = medians$CHR) +
		guides(colour=FALSE) +
		xlab("Chromosome") +
		ylab(expression(-log[10](italic(p)))) +
		scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"))+
		theme_bw() +
		facet_grid_sc(NAME~., scales=list(y=scales_y)) +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
}

plot_scaled_y_boxed_text <- function(data, valcol, path, width, height, res_scale, thresholds, medians, scales_y, rect, text) {
	print(rect)
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_rect(data = rect, aes(xmin = xmin, xmax = xmax, ymin = ymin, ymax = ymax), fill = "#5555DD", color = "#5555DD", alpha = 0.3) +
		geom_point(aes(x = cumsum.tmp, y = VAL, color = color)) +
		geom_hline(data = thresholds, aes(yintercept = THRESH), linetype="dashed") +
		geom_text(data = text, aes(x = x, y = y, label = textlabel)) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$CHR) +
		guides(colour=FALSE) +
		xlab("Chromosome") +
		ylab(expression(-log[10](italic(p)))) +
		scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"))+
		theme_bw() +
		facet_grid_sc(factor(NAME, levels=c("black", "white", "figurita", "runt"))~., scales=list(y=scales_y)) +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
}

get_vert <- function(data, threshold) {
	print("threshold")
	print(threshold)
	print("length(data$cumsum.tmp)")
	print(length(data$cumsum.tmp))
	print("length(data$cumsum.tmp[data$VAL>=threshold])")
	print(length(data$cumsum.tmp[data$VAL>=threshold$THRESH[1]]))
	return(data$cumsum.tmp[data$VAL>=threshold$THRESH[1]])
}

get_verts <- function(data, thresholds) {
	names = as.character(levels(factor(data$NAME)))
	vertslist = sapply(
		names,
		function(x){
			get_vert(data[data$NAME==x,], thresholds[thresholds$NAME==x,])
		}
	)
	verts = Reduce(c, vertslist)
	return(verts)
}

plot_scaled_y_vert <- function(data, valcol, path, width, height, res_scale, thresholds, medians, scales_y, rect) {
	verts = get_verts(data, thresholds)
	vertsd = data.frame(Verts=verts)
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_vline(data = vertsd, aes(xintercept = Verts), color = "#5555dd", alpha=0.3) +
		geom_point(aes(x = cumsum.tmp, y = VAL, color = color)) +
		geom_hline(data = thresholds, aes(yintercept = THRESH), linetype="dashed") +
		scale_x_continuous(breaks = medians$median.x, labels = medians$CHR) +
		guides(colour=FALSE) +
		xlab("Chromosome") +
		ylab(expression(-log[10](italic(p)))) +
		scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"))+
		theme_bw() +
		facet_grid_sc(NAME~., scales=list(y=scales_y)) +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
}

bed2rect <- function(path) {
	# bed = read_bed_noval(path)
	bed = read_bed_postsub(path)
	rect = data.frame(ymin = rep(-Inf, nrow(bed)),
		ymax = rep(Inf, nrow(bed)),
		xmin = bed$cumsum.tmp,
		xmax = bed$cumsum.tmp2)
	return(rect)
}

bed2rect_small <- function(path) {
	bed = read_bed_noval(path)
	rect = data.frame(ymin = rep(-Inf, nrow(bed)),
		ymax = rep(Inf, nrow(bed)),
		xmin = bed$cumsum.tmp,
		xmax = bed$cumsum.tmp2)
	return(rect)
}

plot_cov <- function(data, path, width, height, res_scale, medians) {
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = cumsum.tmp, y = VAL, color = color)) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		guides(colour=FALSE) +
		xlab("Chromosome") +
		ylab("Raw coverage") +
		scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"))+
		ylim(0,300) +
		theme_bw() +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
}

plot_cov_sub <- function(data, path, width, height, res_scale, medians) {
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = cumsum.tmp, y = VAL, color = color)) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		guides(colour=FALSE) +
		xlab("Chromosome") +
		ylab("Raw coverage") +
		scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"))+
		ylim(-300,300) +
		theme_bw() +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
}

plot_cov_multi <- function(data, path, width, height, res_scale, medians, ylimmin, ylimmax) {
	print("ylimmin:")
	print(ylimmin)
	print("ylimmax:")
	print(ylimmax)
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab("Chromosome") +
		ylab("Raw coverage") +
		scale_color_discrete(name = "Dataset")+
		ylim(ylimmin, ylimmax) +
		theme_bw() +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}
		#scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"), name = "Dataset")+

plot_cov_vs_pair <- function(data, path, width, height, res_scale) {
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = VAL1, y=VAL2)) +
		xlim(c(0, 250)) +
		xlab("Coverage") +
		ylab("Pairing proportion") +
		theme_bw()
		print(a)
	dev.off()
}

plot_self_vs_pair <- function(data, path, width, height, res_scale) {
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = VAL1, y=VAL2)) +
		xlab("Self interactions") +
		ylab("Pairing interactions") +
		theme_bw()
		print(a)
	dev.off()
}

plot_self_vs_pair_lim <- function(data, path, width, height, res_scale, ymin, ymax, xmin, xmax) {
	a = ggplot(data = data) +
		geom_point(aes(x = VAL1, y=VAL2)) +
		xlab("Self interactions") +
		ylab("Pairing interactions") +
		theme_bw() +
		lims(x = c(xmin, xmax), y = c(ymin, ymax))

	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a)
	dev.off()
}

plot_self_vs_pair_pretty <- function(data, path, width, height, res_scale, ymin, ymax, xmin, xmax, ylabel, xlabel, textsize) {
	a = ggplot(data = data) +
		geom_point(aes(x = VAL1, y=VAL2), size=0.3) +
		xlab(xlabel) +
		ylab(ylabel) +
		theme_bw() +
		theme(text = element_text(size=textsize)) +
		lims(x = c(xmin, xmax), y = c(ymin, ymax)) +
		geom_smooth(aes(x = VAL1, y=VAL2), method = 'lm', se = TRUE)

	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a)
	dev.off()
}

plot_self_vs_pair_pretty_fixed <- function(data, path, width, height, res_scale, ymin, ymax, xmin, xmax, ylabel, xlabel, textsize) {
	a = ggplot(data = data) +
		geom_point(aes(x = VAL1, y=VAL2), size = 0.3) +
		xlab(xlabel) +
		ylab(ylabel) +
		theme_bw() +
		theme(text = element_text(size=textsize)) +
		lims(x = c(xmin, xmax), y = c(ymin, ymax)) +
		geom_smooth(aes(x = VAL1, y=VAL2), method = 'lm', se = TRUE) +
		coord_fixed(ratio = 1)

	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a)
	dev.off()
}

plot_cov_multi_facet <- function(data, path, width, height, res_scale, medians, ylimmin, ylimmax) {
	print("ylimmin:")
	print(ylimmin)
	print("ylimmax:")
	print(ylimmax)
	print("data head:")
	print(head(data))
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab("Chromosome") +
		ylab("Raw coverage") +
		scale_color_discrete(name = "Dataset")+
		ylim(ylimmin, ylimmax) +
		theme_bw() +
		theme(text = element_text(size=24)) +
		facet_grid(FACET~.)

		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}

plot_cov_multi_facetsc <- function(data, path, width, height, res_scale, medians, scales_y) {
	print("data head:")
	print(head(data))
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab("Chromosome") +
		ylab("Value") +
		scale_color_discrete(name = "Dataset")+
		theme_bw() +
		theme(text = element_text(size=24)) +
		facet_grid_sc(factor(FACET, levels=names(scales_y))~., scales=list(y=scales_y))

		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}

plot_cov_multi_facetsc_boxed <- function(data, path, width, height, res_scale, medians, scales_y, rect) {
	print("rect:")
	print(rect)
	print("data head:")
	print(head(data))
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_rect(data = rect, aes(xmin = xmin, xmax = xmax, ymin = ymin, ymax = ymax), fill = "#5555DD", color = "#5555DD", alpha = 0.3) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab("Chromosome") +
		ylab("Value") +
		scale_color_discrete(name = "Dataset")+
		theme_bw() +
		theme(text = element_text(size=24)) +
		facet_grid_sc(factor(FACET, levels=names(scales_y))~., scales=list(y=scales_y))

		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}

plot_cov_multi_facetsc_names <- function(data, path, width, height, res_scale, medians, scales_y) {
	print("data head:")
	print(head(data))
	print("scales:")
	print(scales_y)
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(LABEL))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab("Chromosome") +
		ylab("Pairing rate (regressed on Ill. cov. and self rate)") +
		scale_color_discrete(name = "Outlier set")+
		theme_bw() +
		theme(text = element_text(size=24)) +
		facet_grid_sc(factor(NAME, levels=names(scales_y))~., scales=list(y=scales_y))

		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}

plot_cov_multi_facetsc_x <- function(data, path, width, height, res_scale, medians, scales_y) {
	print("data head:")
	print(head(data))
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab("Chromosome") +
		ylab("Value") +
		scale_color_discrete(name = "Dataset")+
		theme_bw() +
		theme(text = element_text(size=24)) +
		facet_grid_sc(factor(FACET, levels=names(scales_y))~., scales=list(y=scales_y)) +
		geom_vline(xintercept=23278007)

		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}

plot_cov_multi_pretty <- function(data, path, width, height, res_scale, medians, ylimmin, ylimmax, xlab, ylab, textsize) {
	if (textsize == 0) {
		textsize = 36
	}

	ascale = 1.5
	png(path, width = width * res_scale * ascale, height = height * res_scale * ascale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab(xlab) +
		ylab(ylab) +
		scale_color_discrete(name = "Dataset")+
		ylim(ylimmin, ylimmax) +
		theme_bw() +
		theme(text = element_text(size=textsize))
		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}
		#scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"), name = "Dataset")+

getblue <- function() {
	#return("#619cff")
	return("#00bfc4")
}

plot_cov_multi_pretty_blue <- function(data, path, width, height, res_scale, medians, ylimmin, ylimmax, xlab, ylab) {
	ascale = 1.5
	png(path, width = width * res_scale * ascale, height = height * res_scale * ascale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab(xlab) +
		ylab(ylab) +
		scale_color_manual(name = "Dataset", values = c(getblue()))+
		ylim(ylimmin, ylimmax) +
		theme_bw() +
		theme(text = element_text(size=36))
		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}
		#scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"), name = "Dataset")+

plot_cov_multi_pretty_blue_nolegend <- function(data, path, width, height, res_scale, medians, ylimmin, ylimmax, xlab, ylab) {
	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		a = ggplot(data = data) +
		geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL), color = getblue()) +
		scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
		xlab(xlab) +
		ylab(ylab) +
		scale_color_discrete(name = "Dataset")+
		ylim(ylimmin, ylimmax) +
		theme_bw() +
		theme(text = element_text(size=24))
		print(a)
	dev.off()
		#geom_point(aes(x = cumsum.tmp, y = VAL, color = factor(NAME))) +
}
		#scale_color_manual(values = c(gray(0.5), gray(0), "#EE2222"), name = "Dataset")+

plot_cov_hist <- function(data, path, width, height, res_scale, ylimmin, ylimmax, xlimmin, xlimmax, xlab, ylab, binwidth) {
	a = ggplot(data = data) +
		geom_histogram(aes(VAL), binwidth = binwidth) +
		xlab(xlab) +
		ylab(ylab) +
		scale_color_discrete(name = "Dataset")+
		ylim(ylimmin, ylimmax) +
		xlim(xlimmin, xlimmax) +
		theme_bw() +
		theme(text = element_text(size=18)) +
		theme(axis.text.x = element_text(angle = 22.5, hjust=1))
		# theme(axis.text.x = element_text(angle = 22.5, vjust = 0.5, hjust=1))

	png(path, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a)
	dev.off()
}

plot_box_and_whisker <- function(data, pdfpath, pngpath, width, height, res_scale, ylimmin, ylimmax, xlimmin, xlimmax, xlab, ylab, textsize, fillname) {
	a = ggplot(data = data) +
		geom_boxplot(aes(NAME, VAL1, fill = VAL2)) +
		xlab(xlab) +
		ylab(ylab) +
		scale_fill_discrete(name = fillname) +
		ylim(ylimmin, ylimmax) +
		theme_bw() +
		theme(text = element_text(size=textsize)) +
		theme(axis.text.x = element_text(angle = 22.5, hjust=1))

	pdf(pdfpath, width = width, height = height)
		print(a)
	dev.off()

	png(pngpath, width = width * res_scale, height = height * res_scale, res = res_scale)
		print(a)
	dev.off()
}

# colorseries_palette = c("#00B0F6", "#A3A500", "#E76BF3", "#F8766D", "#00BF7D")
# colorseries_names = c("Pure D. melanogaster", "Lhr rescue", "Hmr rescue", "Hybrid 1", "Hybrid 2")

colorseries_palette = c("#00B0F6", "#F8766D", "#00BF7D", "#A3A500", "#E76BF3")
colorseries_names = c("Pure D. melanogaster", "Hybrid 1", "Hybrid 2", "Lhr rescue", "Hmr rescue")

plot_cov_multi_pretty_colorseries_one <- function(data, path, width, height, res_scale, medians, ylimmin, ylimmax, xlab, ylab, series_names) {
	data$NAME = factor(data$NAME, levels = series_names)

	a = ggplot(data = data) +
	geom_point(aes(x = (cumsum.tmp + cumsum.tmp2) / 2, y = VAL, color = factor(NAME))) +
	scale_x_continuous(breaks = medians$median.x, labels = medians$chrom) +
	xlab(xlab) +
	ylab(ylab) +
	scale_color_manual(name = "Dataset", values = colorseries_palette)+
	ylim(ylimmin, ylimmax) +
	theme_bw() +
	theme(text = element_text(size=36))

	ascale = 1.5
	png(path, width = width * res_scale * ascale, height = height * res_scale * ascale, res = res_scale)
		print(a)
	dev.off()
}

plot_cov_multi_pretty_colorseries <- function(data, outpre, width, height, res_scale, medians, ylimmin, ylimmax, xlab, ylab) {
	for (i in 1:5) {
		outpath = paste(outpre, as.character(i), ".png", sep="")
		series_names = colorseries_names[1:i]
		mdata = data[data$NAME %in% series_names,]
		plot_cov_multi_pretty_colorseries_one(mdata, outpath, width, height, res_scale, medians, ylimmin, ylimmax, xlab, ylab, series_names)
	}
}
