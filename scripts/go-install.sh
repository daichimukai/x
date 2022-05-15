#!/bin/bash

# download and install Go
# link: https://go.dev/doc/install

set -euo pipefail

if [ "$(id -u)" != "0" ]; then
	echo "Error: run as root" >&2
	exit 1
fi

if [ $# != 1 ]; then
	echo "usage: $0 <version>" >&2
	exit 1
fi
version="$1"

tmp=$(mktemp /tmp/go_install.${version}.XXXXXXX.tar.gz)
trap "rm -r ${tmp}" EXIT SIGHUP SIGINT SIGQUIT SIGTERM

url="https://go.dev/dl/go${version}.linux-amd64.tar.gz"
curl -o ${tmp} -sSLf ${url}

target_dir="/usr/local"
go_dir="${target_dir}/go"

rm -rf ${go_dir}
tar -C ${target_dir} -xzf ${tmp}
