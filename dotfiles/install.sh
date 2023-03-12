#!/bin/bash

set -eu

XDG_CONFIG_HOME=${XDG_CONFIG_HOME:-${HOME}/.config}

# `set -e' requires that following environment variables are set.
dirs_env="HOME XDG_CONFIG_HOME"

for d in ${dirs_env}; do
	for f in $(find ${d} -type f); do
		target=$(realpath ${f})
		link_name=$(echo ${f} | sed -e "s:${d}:${!d}:g")

		mkdir -p $(dirname ${link_name})
		ln -sf ${target} ${link_name}

		echo "link: ${target} -> ${link_name}"
	done
done
