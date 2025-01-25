#!/bin/sh

#
# NOT COMMANDS PER SE, JUST A HANDY COLLECTION OF BASH UTILS
#

# taken from: https://gist.github.com/judy2k/7656bfe3b322d669ef75364a46327836#gistcomment-2882842
read_env_var() {
  VAR=$(grep "^$1=" "${PWD}/.env" | xargs)
  IFS="=" read -ra VAR <<< "$VAR"
  IFS=" "
  echo ${VAR[1]}
}
