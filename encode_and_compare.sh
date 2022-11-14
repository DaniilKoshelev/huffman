#!/bin/bash
# shellcheck disable=SC2154
# shellcheck disable=SC1068
# shellcheck disable=SC2037
sum=0
FILES="/home/daniil/Documents/projects/pet/go/src/huffman/resources/*"
for file in $FILES; do
    if [ -f "$file" ]; then
        echo $(basename -- "$file")
        /home/daniil/Documents/projects/pet/go/src/huffman/bin/run_encoder --mode encode -i "$file"
        /home/daniil/Documents/projects/pet/go/src/huffman/bin/run_encoder --mode decode -i "$file.out"
        sdiff -s "$file" "$file.out.out"
        size=$(wc -c < "$file.out")
        sum=$((sum+size))
        echo "$size bytes"
        /home/daniil/Documents/projects/pet/go/src/huffman/bin/run_entropy --mode entropy -i "$file"
        echo "------------"
    fi
done

echo "$sum bytes"