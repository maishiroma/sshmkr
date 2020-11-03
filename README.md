# sshmkr: A pretty, yet powerful SSH Config Manager

> This README is a WIP! 

## Project Goals
- Create a SSH Manager that not only helps streamline adding config to your `~/.ssh/config`, but also can utilize templates within your `~/.ssh/config` to help streamline adding new config options!
- Doing the above, BUT also support formatting your `~/.ssh/config` to abid by your comments!
- Dabble into `go`
- Something to work on while I'm at work!

## Brainstorm
This binary will be able to do the following:
- add: Adds in a new host block to ~/.ssh/config
    - flag 1: --source = specifies a template to go off from; defaults to blank
- delete: Deletes a specific host block
    - flag 1: --source = specifies the config to remove
- copy: Copies a specific hostblock config to a new one.
    - flag 1: --source = specifies the config to use as a base copy or a templated one
- show: Displays a specific host block config
    - flag 1: --source = specifies the host block to show
- list: Prints out all of the host block configs that the ssh config file has


Should create a new template file, named after the default config file (config_templates)

Workflow:

Idea: Organize ssh_config to haev the following hierarchy
```
### Header

## Sub Header

# Comments (ignored)

```