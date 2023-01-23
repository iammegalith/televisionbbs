# televisionbbs
This is the home of TeleVision Bulletin Board System.

building this thing:

```
go mod init televisionbbs
go mod tidy
go build -o bbs main.go
go build -o hello externals/hello/main
go build -o guess externals/guess/main.go
```
then just run bbs and telnet to it on port 8080

## As of Jan 23 2023: STILL INCOMPLETE 
This is incomplete - it kind of works, but there are a lot of fixes required.

If you are playing with this codebase and run into issues or areas of improvement - please create an issue.  Or find me on Discord https://discord.gg/42FXyAU8MN


## Works:  
* Login
* New User Registration
* Runs on Windows, Linux, and Mac so far.
* Menu system basically works
* Show Text Files
* prelogin message/file/text thing. Shows a message prior to login.

## Sort of works:  
* Chatroom - Teleconference - it kind of works. cant do a lot in it.. but.. its sort of there.

## Does not work - but getting there:
* Message Bases - A lot of the code is in place. I just need to refine it.
* File Areas - same here.  A lot of code needs fixing, but the core is there now.
* Doors - I am working on this with extra care. I am building the system to handle redirecting STDIO, so that the doors can be written in any language.  I am about a week out on this, I think.

### This is irritating: 
I am still fiddling with how to display ANSI properly.  Either I get total trash or I get a display with some weird half-height character space between lines.  


Also - still super happy to have anyone who wants to throw some code at functions do the thing ;)    


If I don't have a working admin console from a volunteer by the time I finish the core functionality, I'll ask ChatGPT to make one for me ;)  

# CHANGELOG

## Date: Jan 20 2023
### Updates
* Started adding chatroom in
* Put a bounty on the door system. Its crazy how something so basic is so much of a PITA. 

## Date: Jan 19 2023
### Updates
* Menu system now works mostly - im sure there is something weird...
* Users now checked to see if they are already logged in
* New struct for user - active - this is the bool to see if the user is logged in
* New struct for user - clearscreen - checks to see if the user wants a clearscreen for menus or start of text files.
* Started looking into the message bases - this is a massive undertaking.

## Date: Jan 15 2023
### Updates:
* Started some hefty work on the menu subsystem.
* Menu still does not work - but pretty close.

## Date: Jan 13 2023
### Updates:   
* Moved the source home to [https://github.com/iammegalith/televisionbbs](https://github.com/iammegalith/televisionbbs)
* Changed the name to TeleVision BBS
* SQLite3 for the User and Message Database ( who knew it is good enough for bbs backending without getting all locked up? not me.. )
* INI files for configurations - should make it easy to configure each bit of the system.


Learned a LOT about golang. In a lot of ways, its more forgiving than C.. but in others.. not so much.  So far I am having fun writing things, breaking things, making mistakes, and finding out about really cool go libraries(modules..packages..whatever they are called)  


