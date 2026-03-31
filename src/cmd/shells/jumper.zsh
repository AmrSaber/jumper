eval "$(jumper completion zsh)"

function jump() {
  if [[ $# -eq 0 ]]; then
    jumper list
    return
  elif [[ $# -gt 1 ]]; then
    echo "usage: jump <bookmark>" >&2
    return 1
  fi

  if [[ "$1" == "-" ]]; then
    cd - || return 1
    return
  fi

  local dir
  dir=$(jumper resolve "$@") || return 1
  cd "$dir" || {
    echo "error: failed to cd into '$dir'" >&2
    return 1
  }
}

function _jumper_complete() {
  local candidates
  candidates=($(jumper complete "${words[$CURRENT]}"))
  compadd -U -S '' -- "${candidates[@]}"
  compstate[insert]='menu'
}

compdef _jumper_complete jump
