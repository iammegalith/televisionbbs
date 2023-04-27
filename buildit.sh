#!/bin/bash

# Define the directories to create
outdir="./dist"
bindir="./dist/bin"
datadir="./dist/data"
usermaker="$bindir/usermaker"
export CGO_ENABLED=1

# Function to clean the outdir directory
function clean() {
    if ! rm -rf "$outdir"/*; then
        echo "Error: Failed to clean directory $outdir"
        exit 1
    fi
}

# Function to create the required directories
function create_dirs() {
    if ! mkdir -p "$outdir"/{bin,ansi,ascii,data,files,menus,modules}; then
        echo "Error: Failed to create directory $outdir folders"
        exit 1
    fi
}

# Function to build the tools
function build_tools() {
    # Find all the directories that contain a main.go file
    DIRS=$(find Tools -name "main.go" -type f -exec dirname {} \; | sort)

    # Loop through the directories and build the executables
    for DIR in $DIRS; do
        NAME=$(basename "$DIR")
        echo "Building $NAME..."
        if ! go build -o "$bindir/$NAME" "$DIR/main.go"; then
            echo "Error: Failed to build $NAME"
            exit 1
        fi
    done
}

function build_modules() {
    DIRS=$(find modules -name "main.go" -type f -exec dirname {} \; | sort)
    for DIR in $DIRS; do
        NAME=$(basename "$DIR")
        echo "Building $NAME..."
        mkdir -p $outdir/modules/$NAME
        cp ./modules/$NAME/$NAME.config $outdir/modules/$NAME/
        if ! go build -buildmode=plugin -o "$outdir/modules/$NAME/$NAME" "$DIR/main.go"; then
            echo "Error: Failed to build $NAME"
            exit 1
        fi
    done
}

# Copy ansi, ascii, and menus to dist
function populate_directories() {
  for dir in ansi ascii menus; do
    if ! cp -R "$dir" "$outdir/"; then
        echo "Error: Failed to copy directory $dir"
        exit 1
    fi
  done
  cp ./bbs.config $outdir/
}

# Function to initialize the database and create users
function initialize_database() {
    if ! "$bindir/initdb" "$datadir/userdata.db" "sql_schema.sql"; then
        echo "Error: Failed to initialize database"
        exit 1
    fi

    while IFS=, read -r username level password; do
        cmd="$usermaker -d $datadir -u $username -l $level"
        if [ -n "$password" ]; then
            cmd+=" -p $password"
        fi
        echo "$cmd"
        if ! $cmd; then
            echo "Error: Failed to create user $username"
            exit 1
        fi
    done < "users.csv"
}

# Function to build the BBS
function build_bbs() {
    echo "Building .:: TeleVision BBS ::."
    if ! go build -o "$outdir/television_bbs" -tags=sqlite3 "main.go"; then
        echo "Error: Failed to build bbs"
        exit 1
    fi
}

# Parse command line arguments
case "$1" in
    -h|--help)
        echo "Usage: $0 [clean|all|create_dirs|build_tools|build_modules|populate_directories|initialize_database|build_bbs]"
        echo ""
        echo "clean: removes all files from the $outdir directory"
        echo "all: runs all functions (default)"
        echo "dirs: creates the required directories"
        echo "tools: builds the tools"
        echo "modules: builds the modules"
        echo "populate: copies ansi, ascii, and menus to dist"
        echo "initdb: initializes the database and creates users"
        echo "bbs: builds the BBS"
        exit 0
        ;;
    clean)
        clean
        ;;
    all|"")
        create_dirs
        build_tools
        build_modules
        populate_directories
        initialize_database
        build_bbs
        ;;
    dirs)
        create_dirs
        ;;
    tools)
        build_tools
        ;;
    modules)
        build_modules
        ;;
    populate)
        populate_directories
        ;;
    initdb)
        initialize_database
        ;;
    bbs)
        build_bbs
        ;;
    *)
        echo "Error: Invalid argument"
        exit 1
        ;;
esac
