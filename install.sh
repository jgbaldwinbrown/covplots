#!/bin/bash
set -e

(cd cmd && (
	ls *.go | while read i ; do
		go build $i
	done
))

(cd scripts && (
	ls *.go | while read i ; do
		go build $i
	done
))

cp scripts/full_single_cov_plot ~/mybin/full_single_cov_plot
chmod +x ~/mybin/full_single_cov_plot
cp cmd/all_singlebp_multiline ~/mybin/all_singlebp_multiline
chmod +x ~/mybin/all_singlebp_multiline
cp scripts/plot_single_cov.R ~/mybin/plot_single_cov
chmod +x ~/mybin/plot_single_cov
cp scripts/plot_sub_single_cov.R ~/mybin/plot_sub_single_cov
chmod +x ~/mybin/plot_sub_single_cov

cp scripts/plot_singlebp_multiline_cov.R ~/mybin/plot_singlebp_multiline_cov
chmod +x ~/mybin/plot_singlebp_multiline_cov

cp scripts/plot_singlebp_multiline_cov_facet.R ~/mybin/plot_singlebp_multiline_cov_facet
chmod +x ~/mybin/plot_singlebp_multiline_cov_facet

cp scripts/plot_cov_helpers.R ~/rlibs

cp scripts/plot_cov_vs_pair.R ~/mybin/plot_cov_vs_pair
chmod +x ~/mybin/plot_cov_vs_pair

cp scripts/plot_cov_vs_pair_minimal.R ~/mybin/plot_cov_vs_pair_minimal
chmod +x ~/mybin/plot_cov_vs_pair_minimal
