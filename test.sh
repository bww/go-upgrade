#!/usr/bin/env bash

set -eo pipefail

# where am i?
me="$0"
me_home=$(dirname "$0")
me_home=$(cd "$me_home" && pwd)

# deps
DLV="dlv"

# parse arguments
args=$(getopt dcvS $*)
set -- $args
for i; do
  case "$i"
  in
    -d)
      debug="true";
      shift;;
    -c)
      other_flags="$other_flags -cover";
      shift;;
    -v)
      other_flags="$other_flags -v";
      shift;;
    -S)
      other_flags="$other_flags -tags sqlite3";
      shift;;
    --)
      shift; break;;
  esac
done

if [ ! -z "$debug" ]; then
  GO_UPGRADE_TEST_RESOURCES="${me_home}/test" "$DLV" test $* -- $other_flags
else
  GO_UPGRADE_TEST_RESOURCES="${me_home}/test" go test$other_flags $*
fi
