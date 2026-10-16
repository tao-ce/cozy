#!/bin/bash

force=false

while getopts "m:n:p:f" opt; do
  case $opt in
    m) method=$OPTARG ;;
    n) image_name=$OPTARG ;;
    p) image_path=$OPTARG ;;
    f) force=true ;;
  esac
done

[ -z "$method" ] && echo "option -m method is required" && exit 1
[ -z "$image_name" ] && echo "option -n image_name is required" && exit 1
[ -z "$image_path" ] && echo "option -p image_path is required" && exit 1

platform=$(podman info -f json | jq -r .version.OsArch)
dbus_path="/com/taotesting/cozy"
dbus_interface="com.taotesting.cozy.Installer"

image_exists() {
    podman image exists "$1"
}

remote_pull_image() {
  crane pull \
    --platform "$platform" \
    "$1" \
    /dev/stdout
}

local_pull_image() {
    cat "$1"
}

remote_image_size() {
  crane manifest \
    --platform "$platform" \
    "$1" \
    | jq '.config.size + ([.layers[].size] | add)'
}

local_image_size() {
    wc -c <"$1"
}

progress() {
  pv -i 0.5 -s "$1" -bent \
    --format '%t;%e;%b;%r;%{progress-amount-only};'
}

dbus_forward() {
    while IFS=";" read -r elapsed eta bytes rate percentage _; do
        busctl emit "$dbus_path" "$dbus_interface" "ImageProgress" \
            "a{sv}" 8 \
            method s "$method" \
            name s "$image_name" \
            path s "$image_path" \
            elapsed d "$elapsed" \
            eta s "$eta" \
            bytes t "$bytes" \
            rate d "$rate" \
            percentage d "$percentage"
    done
}

[ "$force" = "false" ] && image_exists "$image_name" && { 
    echo "Image $image_name already exists"
    exit 0
}

coproc dbus_forward { dbus_forward; }
{ cat & } <&${dbus_forward[0]}

case "$method" in
  remote)
    remote_pull_image "$image_path" \
      | progress "$(remote_image_size "$image_path")"  \
      | podman load -q
    ;;
  local)
    local_pull_image "$image_path" \
      | progress "$(local_image_size "$image_path")"  \
      | podman load -q
    ;;
esac 2>&${dbus_forward[1]}

