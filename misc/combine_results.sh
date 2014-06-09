output=data/combined.csv
echo "seed,trans_probability,emitions,states" > $output
cat data/*.dat | sort >> $output
