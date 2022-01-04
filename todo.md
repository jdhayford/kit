- When run using `kit`, or `kit [cmd]` (replace `kit` with normal go run for now)
    - Searches cwd (then up the file directory), for .kit file
    - Parses the .kit file for kit commands
    - Provides a Arrow Key navigable interface for those commands
    - Selecting one will execute that given commands

- .kit file
    - YAML file
    - structure
        - commands (arr)
          - command-alias
            - command itself
            - [optional] description


- Steps
  - [x] Go program that when run, prints out file path of .kit file in cwd or parent..., if found, else prints no .kit
  - [x] Read kit file, parse yaml structure into some sort of go struct
  - [x] Execute kit command, cancel, etc
  - [x] Reorganize code structure
  - [x] Support kit command arguments in yml
  - [] Prompt user for kit argument values (text for now), template into command