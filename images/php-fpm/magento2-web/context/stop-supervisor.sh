#!/bin/sh

printf "READY\n";

while read -r line; do
  echo "Processing Event: $line" >&2;
  kill -3 "$(cat "/run/supervisord.pid")"
done < /dev/stdin
