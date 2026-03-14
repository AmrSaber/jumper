eval "$(jumper completion zsh)"

function jump() {
  local dir
  dir=$(jumper get "$@") || return 1
  cd "$dir" || { echo "error: failed to cd into '$dir'" >&2; return 1; }
}

function _jumper_complete() {
  local candidates
  candidates=($(jumper complete "${words[$CURRENT]}"))
  compadd -U -S '' -a candidates
}

compdef _jumper_complete jump
