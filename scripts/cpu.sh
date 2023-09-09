#!/bin/bash

top -bn 1 \
  | awk '/Cpu\(s\)/ {print 100 - $8}' \
  | awk '{printf "%.2f", $1}'
