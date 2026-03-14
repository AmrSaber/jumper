eval "$(jumper completion bash)"

function jump() {
  local dir
  dir=$(jumper get "$@") || return 1
  cd "$dir" || { echo "error: failed to cd into '$dir'" >&2; return 1; }
}

function _jumper_complete() {
  local prefix="${COMP_WORDS[$COMP_CWORD]}"
  local candidates
  mapfile -t candidates < <(jumper complete "$prefix")
  COMPREPLY=($(compgen -W "${candidates[*]}" -- ""))
  compopt -o nospace 2>/dev/null
}

complete -F _jumper_complete jump
