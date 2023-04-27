## Building TeleVisionBBS from Source

### Recommended:

TeleVisionBBS is using the following environment, and it is suggested that to minimize issues building and deploying it, you might consider the same.

* OS: Ubuntu 22.04LTS
* Go: Always current stable
* DB: sqlite3

### Prepare for the build: 
I use go-sqlite, which requires CGO. If you want to build this, you will need to install gcc.  To avoid dependancy heck, it's easier to just install the build-essential package.

```
sudo apt install build-essential
```

Install golang (visit https://go.dev/dl/)  - latest stable version

In the TeleVisionBBS source directory, run:

```
go mod tidy
```

Build the entire thing:

```
./buildit.sh all
```

You should now find everything you need in the dist/ folder.

---

### Using the build script

Usage: **./build.sh** [*clean*|*all*|*dirs*|*tools*|*modules*|*populate*|*initdb*|*bbs*]

* **clean:** removes all files from the ./dist directory
* **all:** runs all functions (default)
* **dirs:** creates the required directories
* **tools:** builds the tools
* **modules:** builds the modules
* **populate:** copies ansi, ascii, exampledoors, and menus to dist
* **initdb:** initializes the database and creates users
* **bbs:** builds the BBS


The clean command removes all files from the ./dist directory. The all command runs all functions, which include creating the required directories, building the tools, building the modules, copying ansi, ascii, exampledoors, and menus to dist, initializing the database and creating users, and building the BBS. The dirs command creates the required directories. The tools command builds the tools. The modules command builds the modules. The populate command copies ansi, ascii, exampledoors, and menus to ./dist. The initdb command initializes the database and creates users. The bbs command builds the BBS.

You can run the script by passing one of these commands as an argument. If you don't pass an argument, the all command will be run by default.

