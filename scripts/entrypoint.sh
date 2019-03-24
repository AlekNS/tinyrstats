#!/bin/sh

set -eu

for env in $(env); do
  key=${env%=*}
  [ -n "${key#*_FROM_FILE}" ] && continue
  val=${env#*=}
  key=${key%_FROM_FILE}
  val=$(cat $val)
  export "${key}"="${val}"
done

exec $@
