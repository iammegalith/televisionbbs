# <p style="text-align: center;">TeleVision Bulletin Board System<br>SysOps Guidebook</p>
<p style="text-align: center;">Version 1Q2023.1</p>

Note that all configuration files, *.asc files, and any source code you may edit can be done in any text editor of your choosing.  I typically use vim - but you can use Notepad, TextEdit, VS Code, etc.

## BBS Files, what they do, and how to edit them.

### TeleVision BBS Folder Hierarchy
```
TelevisionBBS
├── doors
│   ├── hello
│   │   ├── hello(.exe)
│   ├── guess
│       ├── guess(.exe)
├── config
│   ├── actions.config
│   ├── bulletins.config
│   ├── doors.config
│   ├── main.config
│   ├── messages.config
│   ├── strings.config
├── tvbbs(.exe)
├── README.md
├── LICENSE 
└── bbs.config
```

The bbs.config file is the base configuration for TeleVision BBS.
It resides in the root director of the BBS hierarchy.

```ini
[mainconfig]
port = 8080
bbsname = "Television BBS"
sysopname = "Elliot Gould"
prelogin = true
bulletins = true
newregistration = true
defaultlevel = 1
configpath = "config/"
ansipath = "ansi/"
asciipath = "ascii/"
modulepath = "modules/"
datapath = "data/"
filespath = "files/"
stringsfile = "strings.config"
```
port - this is the TCP port you would like your BBS to listen on.  I like using a high port (greater than 1024).  Default is port 8080.

bbaname - this is the name of the bbs.  The name you want it to be called.  For instance, if I want to call my BBS "The SysOps Lair", i would change this parameter to:

```
bbsname = "The SysOps Lair"
```

sysopname - this is you name/handle (screen name).

prelogin - this shows an ascii text file called prelogin.asc.  You can skip it by setting this to false.  However, this is genrally a good way to show a banner about the BBS, or warnings in case you run a BBS that may have content not suitable for all.

bulletins - you can have several bulletin files.  These are usually things like the BBS rules, Info about upcoming events, or general news.

newregistration - this is set to true by default.  It can be turned off by setting this value to false if you want to run a private invite only BBS.

defaultlevel - this sets the default Level of the user. Access levels are any whole number from 0 to 255, and can be used to restrict access to certain parts of the bulletin board system.

configpath - by default this is in the config/ folder. You may wish to change this for some reason, so - here's your rope.

ansipath - by default, all ansi files are stored in the ansi/ folder.  You may wish to change this to something else - maybe "text/ansi" .. or not.

asciipath - by default, all ascii text files are stored in the ascii/ folder.  You may wish to change this to something else - maybe "text/ascii"

modulepath - this is where 3rd party applications live.  Doors, utilities, etc. By default: modules/

datapath - this is where databases live. Message Bases, User Base, etc. By default: data/

filespath - files/ is the parent directory for your file upload/download area.  You may wish to point this to a larger storage option.

stringsfile - all responses given by the BBS are configurable here.  The file config/strings.config are where you will find them.  
