#!/usr/bin/env bash
# shellcheck disable=SC2124

_echo() {
  local to="$1"
  case "$to" in
    stdout)
      echo "${@:2}" >&1
      ;;
    stderr)
      echo "${@:2}" >&2
      ;;
    *)
      echo "$0 <stdout|stderr> ..." >&2
      exit 1;
      ;;
  esac
}

_echo "$@"
