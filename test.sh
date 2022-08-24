#!/bin/bash

ROOTDIR=$(dirname $0)

for m in $(ls -1d $ROOTDIR/*/ | grep -Eo '\w.+\w') ; do
  coverage=$(go test -race -cover ./$m | grep -Eo '[.0-9]+%' | grep -Eo '[.0-9]+')
  if (( $(echo 1 | awk "{print ($coverage >= $COVERAGE_THRESHOLD)}") )); then
    printf '%-50s: %-6s (%3.2f%%)\n' $m Ok $coverage
  else
    modules="$modules $m"
    printf '%-50s: %-6s (%3.2f%%)\n' $m Failed $coverage
  fi
done

if [ "$modules" ]; then
  echo "Current test coverage is below threshold ($COVERAGE_THRESHOLD%)."
  echo "Please add more unit tests to the following modules:"
  echo $modules
  exit 1
fi
