#!/bin/bash

# generate $nfiles large files in a dataset directory
nfiles=$1

# File size in megabytes
file_size_mb=5

mkdir -p /tmp/dataset

for i in $(seq 1 $nfiles); do
    # Generate a large file with the specified size
    dd if=/dev/urandom of=/tmp/dataset/file.$i bs=1M count=$file_size_mb status=progress
done