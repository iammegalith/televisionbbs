# What is a Bulletin Board System?

Before the internet, noone had smartphones, there were only a few large networked systems - but they cost too much to run, so the price was too high for the average person to use.

some very clever computer hobbyists discovered that they could host "server" type systems on their home computer and allow others to connect and interact.

eventually this grew from just a simple online Bulletin Board (like at the supermarket or library where people post notices).

Some systems allowed people to post messages (like web forums), upload and download files, and the really advanced systems had several modems and phone lines and could handle several users at once.

TeleVision Bulletin Board System aims to be a modern version of the BBS's of the old days.

A person should be able to telnet to the TeleVision BBS, create a new account, read and post messages in various message bases, upload and download files, play text based games, and chat with one another.

# How does TeleVision BBS work?

Internally, a telnet server and listener are set up.  The BBS listens for incoming connections.

When a connection is made, the main function hands off to the handleconnection function.

This is where the user map is create and populated with telnet.Connection data.

I wonder if handleConnection should be a go routine?

the user logs into the system, and is sent to the menu loop.  everything happens here. 

When the user logs off, the session is closed and the users temporary data is purged - anything changing in the in-memory data is written to the database.

If the users connection drops unexpectedly, the checkDisconnectedClients function should handle closing the session.

User should be able to send messages to each other, see who is online, and interact via channels.



# Directory Listing

-----------------

```
├───config
|      ├──television.conf
|      └──strings.conf
├───data
|      └──bbs.db
├───textfiles
|      ├──info.txt
|      └──prelogin.txt
├───util
|      └──main.go
└─main.go
```

# config folder
television.conf: this file contains the configuration for the system.  It sets the port number to listen on, name of the system, file paths, and sets boolean flags for the global system.

strings.conf: this allows you to change the menu prompts, and give a level of customization to the system.

# data folder
bbs_schema.sql: copy the contents of this to sqlite3 to create the bbs.db database.  This is a sqlite3 database used to store user data, session data (maybe - i dont know if it is needed here or not) 

``` bash
name@host (/televisionbbs/) $ sqlite3 data/bbs.db
SQLite version 3.40.1 2022-12-28 14:03:47
Enter ".help" for usage hints.
sqlite>  CREATE TABLE users (
   ...>     id INTEGER PRIMARY KEY,
   ...>     username TEXT NOT NULL,
   ...>     password TEXT NOT NULL,
   ...>     level INTEGER NOT NULL,
   ...>     linefeeds INTEGER NOT NULL,
   ...>     hotkeys INTEGER NOT NULL,
   ...>     active INTEGER NOT NULL,
   ...>     clearscreen INTEGER NOT NULL
   ...> );
sqlite> 
sqlite> CREATE TABLE sessions (
   ...>   id INT AUTO_INCREMENT PRIMARY KEY,
   ...>   username VARCHAR(255) NOT NULL,
   ...>   active TINYINT(1) DEFAULT 1,
   ...>   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   ...> );
sqlite> .q
name@host (/televisionbbs/) $
```

# textfiles folder
This directory contains textfiles used to display to the user. Right now there are only the two files and they are for testing pagination, and making sure things display properly with varying user terminals.

info.txt: Eventually will have a little information about the system itself.

prelogin.txt: this file is displayed prior to a user logging into the system.  It could show some rules or a news bulletin.

## util folder
This is the util package - it contains some things that can be accessed by other go modules that will be added later, once the core functionality is in the main program.

main.go: this is the actual util package file. it contains global constants (ANSI_*), small functions that can be used from other programs (GenRandomString() )

# the main application

to build this, make sure we have the data base initialized.

``` bash
go build -o bbs main.go
```

You may need to setup the go mod stuff:  
  
``` bash
go mod init televisionbbs
go mod tidy
```

The following is the flow of the main package:

func main()
  - calls getConfig() and populates the system configuration
  - call getStrings() and populates the system strings 
  - runs a go routine to watch to see if a user disconnects
  - sets up a telnet server using the handleConnection function when a new session is established.
  - starts the telnet listener

When the application starts, main() loads its configs, creates a watcher go routine to monitor all sessions for a disconnect, then when a user telnets to the server, it creates a new connection and hands the session off to handleConnection.

func checkDisconnectedClients()
  - go routine that watches the system for disconnects. 
  - if a connection is dropped, this should remove the session information so that the user no longer shows up in the [w]ho menu function, and so that the system does not try to send data to that users connection.
  
func handleConnection(conn *telnet.Connection, TheUser map[string]CurrentUser)
  - opens database (should this be here or in main?)
  - sends connection string (Welcome to TeleVision BBS)
  - calls login(conn, db, TheUser)
  - if login is successful, calls the menu loop
  - when a user logs off, locks maps, removes user data from the maps, unlocks maps
  - exits
  
  the job of this is to get the user logged in to the system, and run the menu loop.  I think this is where the unique user map should be initialized with just the telnet.Connection info and maybe a uuid or generated unique id of some kind.  maybe i am over thinking it.
  
func checkAnsi
  - checks to see if the user can use ANSI color in their terminal
  - returns boolean
 
func showTextFiles
  - displays a text file to the user

func menu
  - this is the main logic loop
  - presents a menu to the user
  - user interacts with the menu
  - user can toggle ANSI mode on and off
  - user can list users with the [w]ho command, which _SHOULD_ return the list of logged in users and what they are doing.
  - user can view a text file showing system info [i]nfo
  - user can gracefully exit the system via [g]oodbye
  
func alreadyLoggedIn
  - this should check to see if a user is already logged in
  - this function is called by login() to ensure that a user is only logging in on one connection, and cannot have more than one simultaneous connection.
  - if a user enters the correct username and password, they are allowed to kick the other connection offline so they can log in.
    - this is useful in the case where the connection watcher (checkDisconnectedClients()) fails, or someone else has logged in to the users account.

func login
  - this is the login function.
  - looks for a match of username/passwd in the database
  - sets the users ephemeral variables for the duration of the session
  - if the user is new, send them to newUser()
  - once the user is logged in, send them to the menu loop
  - if a password fails 3 times, kick the user off
  
func checkAttempts
  - if attempts > 3 then send message and kick the user off

func newUser
  - checks the username in the database
  - if the numrows = 0 then allow the user to sign up
  - otherwise tell them to try another name
  - called by the login function
  - returns to login function where the users ephemeral variables should be set for the duration of the session
  
func hashPassword
  - uses bcrypt to create a password hash to store in the database so we arent storing plain text.
  - called by login and newuser

func pressKey
  - displays some text and tells the user to press a key to continue
  - used for pagination

func logout
  - gracefully disconnects the user
  - asks if they are sure they want to logout
  - if yes, disconnect, clear ephemeral variables
  - if no, return to previous menu
  - called by [g]oodby in main menu

func askYesNo
  - asks the user a question and returns boolean 

func listUsers
  - this is supposed to show a list of users currently online and what they are doing
  
func toggleAnsi
  - this turns ANSI mode on and off, returns boolean

func writeLine
  - this prints a string to the user connection
  - this is what allows text to be sent to the user terminal
  
func readLine
  - reads from user terminal and puts the result in a string
  
func readKey
  - reads a single keypress from the user
  - useful for askYesNo and menu items

func convertBytetoString
  - converts byte[] to string

func convertError
  - converts error message to string

func getStrings
  - reads string file from strings.conf
  - applies values to variables

func getConfig
  - reads television.conf configuration file
  - applies values to variables
  - this is the primary config file for the BBS

func init()
  - this is where the maps are being setup.
  - i dont know if this is even the right place.
  
