#!/bin/bash
# shellcheck disable=SC2154
# shellcheck disable=SC1068
# shellcheck disable=SC2037
FILES="/home/daniil/Documents/projects/pet/go/src/huffman/resources/*"
for file in $FILES; do
    if [ -f "$file" ]; then
        /home/daniil/Documents/projects/pet/go/src/huffman/bin/run_encoder --mode encode -i "$file"
        /home/daniil/Documents/projects/pet/go/src/huffman/bin/run_encoder --mode decode -i "$file.out"
        sdiff -s "$file" "$file.out.out"
    fi
done