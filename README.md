# televisionbbs
This is the home of TeleVision Bulletin Board System.

## INCOMPLETE 13-01-2023
This is incomplete - it kind of works, but there are a lot of fixes required.

# CHANGELOG

## Date: Jan 13 2023
### Updates:   
* Moved the source home to [https://github.com/iammegalith/televisionbbs](https://github.com/iammegalith/televisionbbs)
* Changed the name to TeleVision BBS
* SQLite3 for the User and Message Database ( who knew it is good enough for bbs backending without getting all locked up? not me.. )
* INI files for configurations - should make it easy to configure each bit of the system.


Learned a LOT about golang. In a lot of ways, its more forgiving than C.. but in others.. not so much.  So far I am having fun writing things, breaking things, making mistakes, and finding out about really cool go libraries(modules..packages..whatever they are called)  


Works:  
* Basic Login
* Basic Registration (Hashed passwords in the database)
* Runs on Windows, Linux, and Mac so far.


Sort of works:  
* I have a basic "door" system that kind of works.  I am not happy enough to call it "ok for now" 
   
I am still fiddling with how to display ANSI properly.  Either I get total trash or I get a display with some weird half-height character space between lines.  


Also - still super happy to have anyone who wants to throw some code at functions do the thing ;)    


If I don't have a working admin console from a volunteer by the time I finish the core functionality, I'll ask ChatGPT to make one for me ;)  
