#!/bin/bash
set -e

cp scripts/plot_single_cov.R ~/mybin/plot_single_cov
chmod +x ~/mybin/plot_single_cov
cp scripts/plot_sub_single_cov.R ~/mybin/plot_sub_single_cov
chmod +x ~/mybin/plot_sub_single_cov
cp scripts/plot_singlebp_multiline_cov.R ~/mybin/plot_singlebp_multiline_cov
chmod +x ~/mybin/plot_singlebp_multiline_cov
cp scripts/plot_cov_helpers.R ~/rlibs

cp scripts/plot_cov_vs_pair.R ~/mybin/plot_cov_vs_pair
chmod +x ~/mybin/plot_cov_vs_pair

cp scripts/plot_cov_vs_pair_minimal.R ~/mybin/plot_cov_vs_pair_minimal
chmod +x ~/mybin/plot_cov_vs_pair_minimal
