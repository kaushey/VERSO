# Define the array of options
options=("init" "add" "commit" "log" "status" "cat-file")

# Function to handle autocomplete
_verso_autocomplete() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    COMPREPLY=( $(compgen -W "${options[*]}" -- "$cur") )
}

# Register the autocomplete function for the 'verso' command
complete -F _verso_autocomplete verso
