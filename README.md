# sshmkr: A pretty, yet powerful SSH Config Manager

> Picture this: Ever wanted to have a quick way to manage your ssh_config file without downloading any pre-requesites? Or better yet, do you like to organize your ssh_config to be readable? Perhaps there's some hosts that you like to reuse their configuration for other ssh hosts, since their configs are too long to remember. 

Sounds like you? Let me introduce you to this new tool, `sshmkr`! This is a lightweight Go binary that allows one to manage their ssh_configs directly from the CLI.

## Table Of Contents
- [How to Install](#How-to-Install)
- [Workflow](#Workflow)
- [Commands](#Commands)
- [Contribute](#Contribute)
- [Inspiration](#Inspiration)

## How to Install
There are two ways of installing this binary onto your system:
1. Refer to the latest releases in the release panel (as of now, the latest one is v[1.1.0](https://github.com/maishiroma/sshmkr/releases/tag/1.1.0))
2. Clone the repository and build/install the binary yourself with `go build .`.

Once you have the binary, simply place it in your `$PATH` and refer to it via `sshmkr`. To test if it is working, enter:
```
sshmkr --version
```

## Workflow
By default, `sshmkr` looks in `~/.ssh/config` for your ssh_config and `~/.ssh/config_templates` for all of the templates it can leverage. This can be changed via the `--path` flag.

Before using `sshmkr`, one needs to create the `config_templates` file, which can be as simple as the following:

```
# You can add a comment like this!
Host template_1
    Hostname something
    Port 22

# Another one!
Host cool_template_2
    Hostname another
    Port 20
    IdentityFile ~/.ssh/id_rsa
```
The functionality of this file will be described in greater detail in a future section.

One also will need to add in __headers__ to their ssh_config file. These are simple comments that are placed in the ssh_config to help organize the ssh_hosts.
```
#### This is a main header

## This is a sub header
Host host_1
    Hostname something
    Port 22

## This is another sub header
Host host_2
    Hostname another
    Port 20
    IdentityFile ~/.ssh/id_rsa

#### This is a new main Header
```

There are two types of headers that the binary supports:
- `####`: This is a main header, which are the high level organization sections. Think of these like root directories in a file system.
- `##`: This is a sub header, which corresponds to last declared main header. Think of this as a sub directory in a root directory in a file system.

> `#` are treated as comments in config files, so if additional comments need to be made, use these.

These are used to help organize what host configs correspond to specific organization levels that one may have sorted their host file. 

### Templates
`sshmkr` utilizes an external file, `config_templates`, that is located in `~/.ssh/` by default. This file has the exact same syntax as a normal ssh_config file.

This config file will be used as a basis when adding in new ssh_hosts via the `add` command. The template acts as a way to specify __default__ values for specific host configs so that you can simply reuse them to your heart's content.

### Headers
These are specialized comments that are present in the ssh_config file. They are used to organize ssh headers into specific categories for the binary to sort these in.

## Commands
Below is a list of available commands that can be utilized. 

### Add
Adds a new configuration to the ssh_config. This utilizes the `config_template` that was defined in `~/.ssh` (by default). Upon calling this function, the binary will guide the user through where to put the new host, via the headers.

Upon completion, the new ssh header will be truncated __before__ the next declared header.

Note that if the template is commented out via `#`, this command will ignore said template.

Example:

```
$ cat ~/.ssh/config_template
Host sampleTemplate
    Hostname myhost
    Port 22
    IdentityFile ~/.ssh/id_rsa
    ProxyCommand ssh -F ~/.ssh/config -W %h:%p personal_jb

$ sshmkr add --source sampleTemplate
~ Main Header Selection ~
1.)  Personal
2.)  Project 1
3.)  Project 2
Select a main header: 2

~ Sub Header Selection ~
1.)  Sites
2.)  Jumpboxes
3.)  Instances
Select a sub header: 3

~ Template ~
Enter a value for Host [ default: NewHost ]: someHost
Enter a value for Hostname [ default: myhost ]: 111.1111.111
Enter a value for Port [ default: 22 ]:
Enter a value for IdentityFile [ default: ~/.ssh/id_rsa ]:
Enter a value for ProxyCommand [ default: ssh -F ~/.ssh/config -W %h:%p personal_jb ]:

Sucessfully added host someHost to config!

$ cat ~/.ssh/config 
...
#### Project_1

## Instances
Host someHost
	Hostname 111.1111.111 
	Port 22 
	IdentityFile ~/.ssh/id_rsa 
	ProxyCommand ssh -F ~/.ssh/config -W %h:%p personal_jb 

#### Project_2
```

### Delete
Removes a specific host config that is specified when calling this command.

This removes the first host config that matches in the ssh_config starting from the top of the file. Like `add`, if the host config is commented out via `#`, this command will ignore looking at thoses host names.

Example:
```
$ sshmkr delete --source NewHost
Sucessfully removed host NewHost from ssh_config!
```

### Comment
This comments out the specified host config from the ssh_config file. This in of itself prevents that host config to be read by any of the other commands here as well as used in other standard CLI commands.

This command is smart enough to know when to comment in/out said host config. The precedence of how commenting works is that the command will look for the first hostname that it finds in the ssh_config, starting from the top.

Exampe:
```
$ sshmkr comment --source NewHost
Sucessfully commented out host NewHost

$ cat ~/.ssh/config
...
#Host NewHost
#	Hostname myhost 
#	Port 22 
#	IdentityFile ~/.ssh/id_rsa 
#	ProxyCommand ssh -F ~/.ssh/config -W %h:%p personal_jb

$ sshmkr comment --source NewHost
Sucessfully uncommented out host NewHost

$ cat ~/.ssh/config
...
Host NewHost
	Hostname myhost 
	Port 22 
	IdentityFile ~/.ssh/id_rsa 
	ProxyCommand ssh -F ~/.ssh/config -W %h:%p personal_jb 
```

### Copy
Copies an __existing__ host config that is present in the ssh_config and uses it as a template for a new config. This works by performing a front to end search on all of the hosts in a file. As such, if there are multiple iterations of a host name, this command will only take the first host that it finds. 

This is especially handy if you are reusing configs because they share similar atrributes (i.e. ports, proxys). 

Like the other commands, if the host specified is commented out, this command will ignore it in its selection

Example:
```
$ sshmkr copy --source github.com
Successfully read host config github.com from file!

~ Main Header Selection ~
1.)  Personal
2.)  Project 1
3.)  Project 2
Select a main header: 3

~ Sub Header Selection ~
1.)  Sites
2.)  Jumpboxes
3.)  Instances
Select a sub header: 3

~ Template ~
Enter a value for Host [ default: NewHost ]: AnotherOne
Enter a value for User [ default: maishiroma ]: notmaishiroma
Enter a value for IdentityFile [ default: ~/.ssh/personal_projects/github_id_rsa ]: ~/.ssh/id_rsa

Sucessfuly created new host AnotherOne from template!

$ cat ~/.ssh/config
Host github.com
    User maishiroma 
	IdentityFile ~/.ssh/id_rsa

....

Host AnotherOne
	User notmaishiroma 
	IdentityFile ~/.ssh/id_rsa 
```

### Show
Displays the specified host config out to the console. 

This us useful if you want to quickly peak at a particular host config or if you want to interpolate the contents into another CLI command.

Like the other commands, if the config is commeted out, this command will ignore it.

Example:
```
$ sshmkr show --source NewHost
Host  NewHost
 Hostname myhost
 Port 22
 IdentityFile ~/.ssh/id_rsa
 ProxyCommand ssh -F ~/.ssh/config -W %h:%p personal_jb
```

### Edit
This command allows for an in-line edit on a given ssh host. This is useful if one needs to alter a specific item in a config without having to manually create a brand new config for the same host.

Example:
```
$ sshmkr edit --source NewHost
Found host config to use for template, NewHost ...


~ Template ~
Enter a value for Host [ default: NewHost ]: EditedHost
Enter a value for Hostname [ default: myhost  ]: anotherOne
Enter a value for Port [ default: 22  ]:
Enter a value for IdentityFile [ default: ~/.ssh/id_rsa  ]:
Enter a value for ProxyCommand [ default: ssh -F ~/.ssh/config -W %h:%p personal_jb  ]:

Sucesfully edited host config, NewHost !

$ sshmkr show --source EditedHost
Host  EditedHost
 Hostname anotherOne
 Port 22
 IdentityFile ~/.ssh/id_rsa
 ProxyCommand ssh -F ~/.ssh/config -W %h:%p personal_jb
```

## Contribute
This project is free to be leveraged by whoever else finds this helpful. If one wants to request for more features and/or issues, feel free to open up new issues/forks on this repository! Just make sure to ping me in them so that I can take a look at your inquiry. 

## Inspiration
This project was made out of both wanting to get more familiar with golang, but to also satisfy a personal need that I use all the time at work. Plus, this project used the following projects as inspiration points:
- [kevinburke/ssh_config](https://github.com/kevinburke/ssh_config) for the helpful library that helped parse the SSH config file for this project
- [dbrady/ssh-config](https://github.com/dbrady/ssh-config) for the initial idea of a golang CLI binary used to interact with SSH configs
- [sarcasticadmin/sshcb](https://github.com/sarcasticadmin/sshcb) another inspiration point on the project