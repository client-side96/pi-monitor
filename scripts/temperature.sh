#!/bin/bash

/usr/bin/vcgencmd measure_temp \
  | cut -d "=" -f2 \
  | cut -d "'" -f1 \
  | tr -d '\n'