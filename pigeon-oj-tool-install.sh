#!/usr/bin/env bash

# This is the bootstrap Unix installer served by `https://get.pigeon-oj.cloud`.
# Its responsibility is to query the system to determine what OS the system
# has, fetch and install the appropriate build of pigeon-oj-tool.

# NOTE: to use an internal company repo, change how this determines the latest version
get_latest_release() {
  curl --silent "https://get.pigeon-oj.cloud/latest-version"
}

release_url() {
  echo "https://github.com/Pigeon-Developer/pigeon-oj-tool/releases"
}

download_release_from_repo() {
  local version="$1"
  local os_info="$2"
  local tmpdir="$3"

  local filename="pojt-$version-$os_info.tar.gz"
  local download_file="$tmpdir/$filename"
  local archive_url="$(release_url)/download/v$version/$filename"

  curl --progress-bar --show-error --location --fail "$archive_url" --output "$download_file" --write-out "$download_file"
}

usage() {
    cat >&2 <<END_USAGE
pigeon-oj-tool-install: The installer for pigeon-oj-tool

USAGE:
    pigeon-oj-tool-install [FLAGS] [OPTIONS]

FLAGS:
    -h, --help                  Prints help information

OPTIONS:
        --version <version>     Install a specific release version of pigeon-oj-tool
END_USAGE
}

info() {
  local action="$1"
  local details="$2"
  command printf '\033[1;32m%12s\033[0m %s\n' "$action" "$details" 1>&2
}

error() {
  command printf '\033[1;31mError\033[0m: %s\n\n' "$1" 1>&2
}

warning() {
  command printf '\033[1;33mWarning\033[0m: %s\n\n' "$1" 1>&2
}

request() {
  command printf '\033[1m%s\033[0m\n' "$1" 1>&2
}

eprintf() {
  command printf '%s\n' "$1" 1>&2
}

bold() {
  command printf '\033[1m%s\033[0m' "$1"
}

# check for issue with POJT_HOME
# if it is set, and exists, but is not a directory, the install will fail
pojt_home_is_ok() {
  if [ -n "${POJT_HOME-}" ] && [ -e "$POJT_HOME" ] && ! [ -d "$POJT_HOME" ]; then
    error "\$POJT_HOME is set but is not a directory ($POJT_HOME)."
    eprintf "Please check your profile scripts and environment."
    return 1
  fi
  return 0
}


# returns the os name to be used in the packaged release
parse_os_info() {
  local uname_str="$1"
  local arch="$(uname -m)"

  case "$uname_str" in
    Linux)
      if [ "$arch" == "i386" ]; then
        echo "linux-i386"
      elif [ "$arch" == "i486" ]; then
        echo "linux-i386"
      elif [ "$arch" == "i586" ]; then
        echo "linux-i386"
      elif [ "$arch" == "i686" ]; then
        echo "linux-i386"
      elif [ "$arch" == "x86_64" ]; then
        echo "linux-amd64"
      elif [ "$arch" == "aarch64" ]; then
        echo "linux-arm64"
      elif [ "$arch" == "loongarch64" ]; then
        echo "linux-loong64"
      else
        error "Releases for architectures other than x64 and arm are not currently supported."
        return 1
      fi
      ;;
    *)
      return 1
  esac
  return 0
}

parse_os_pretty() {
  local uname_str="$1"

  case "$uname_str" in
    Linux)
      echo "Linux"
      ;;
    *)
      echo "$uname_str"
  esac
}

# return true(0) if the element is contained in the input arguments
# called like:
#  if element_in "foo" "${array[@]}"; then ...
element_in() {
  local match="$1";
  shift

  local element;
  # loop over the input arguments and return when a match is found
  for element in "$@"; do
    [ "$element" == "$match" ] && return 0
  done
  return 1
}

create_tree() {
  local install_dir="$1"

  info 'Creating' "directory layout"

  # pigeon-oj-tool/
  #     bin/

  mkdir -p "$install_dir" && mkdir -p "$install_dir"/bin
  if [ "$?" != 0 ]
  then
    error "Could not create directory layout. Please make sure the target directory is writeable: $install_dir"
    exit 1
  fi
}

install_version() {
  local version_to_install="$1"
  local install_dir="$2"

  if ! pojt_home_is_ok; then
    exit 1
  fi

  case "$version_to_install" in
    latest)
      local latest_version="$(get_latest_release)"
      info 'Installing' "latest version of pigeon-oj-tool ($latest_version)"
      install_release "$latest_version" "$install_dir"
      ;;
    *)
      # assume anything else is a specific version
      info 'Installing' "pigeon-oj-tool version $version_to_install"
      install_release "$version_to_install" "$install_dir"
      ;;
  esac

  if [ "$?" == 0 ]
  then
    info 'Finished' "installation. Updating bin file link."
    ln -s "$install_dir"/bin/pigeon-oj-tool /usr/local/bin/pjot
    "$install_dir"/bin/pigeon-oj-tool setup
  fi
}

install_release() {
  local version="$1"
  local install_dir="$2"

  download_archive="$(download_release "$version"; exit "$?")"
  exit_status="$?"
  if [ "$exit_status" != 0 ]
  then
    error "Could not download pigeon-oj-tool version '$version'. See $(release_url) for a list of available releases"
    return "$exit_status"
  fi

  install_from_file "$download_archive" "$install_dir"
}

download_release() {
  local version="$1"

  local uname_str="$(uname -s)"
  local os_info
  os_info="$(parse_os_info "$uname_str")"
  if [ "$?" != 0 ]; then
    error "The current operating system ($uname_str) does not appear to be supported by pigeon-oj-tool."
    return 1
  fi
  local pretty_os_name="$(parse_os_pretty "$uname_str")"

  info 'Fetching' "archive for $pretty_os_name, version $version"
  # store the downloaded archive in a temporary directory
  local download_dir="$(mktemp -d)"
  download_release_from_repo "$version" "$os_info" "$download_dir"
}

install_from_file() {
  local archive="$1"
  local install_dir="$2"

  create_tree "$install_dir"

  info 'Extracting' "pigeon-oj-tool binaries"
  # extract the files to the specified directory
  tar -xf "$archive" -C "$install_dir"/bin
}

# return if sourced (for testing the functions above)
return 0 2>/dev/null

# default to installing the latest available version
version_to_install="latest"

# install to POJT_HOME, defaulting to /etc/pigeon-oj-tool
install_dir="${POJT_HOME:-"/etc/pigeon-oj-tool"}"

# parse command line options
while [ $# -gt 0 ]
do
  arg="$1"

  case "$arg" in
    -h|--help)
      usage
      exit 0
      ;;
    --version)
      shift # shift off the argument
      version_to_install="$1"
      shift # shift off the value
      ;;
    *)
      error "unknown option: '$arg'"
      usage
      exit 1
      ;;
  esac
done

install_version "$version_to_install" "$install_dir"
