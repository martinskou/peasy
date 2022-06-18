

# setup ctrl-c to terminate the script by killing the go program.

trap ctrl_c INT

function ctrl_c() {
    printf "Exiting\n"
    kill $PID
    exit 1
}

# make sure we have a place for the binary file.

mkdir -p bin

EXEC="bin/server"

# loop which is executed each time change to the source files are detected.

while true; do
    rm $EXEC 2>/dev/null
    
    go build -o $EXEC src/* # build the server

    if [ -f "$EXEC" ]; then
        # the new executeable was found, build was a success.
        ./bin/server &
        PID=$!
        echo "Started with PID $PID in $PWD\n"
    else
        # no exec was found, something failed.
        echo "Build failed"
    fi

    # wait for new changes to the source files.
    fswatch -1 -r -e ".*" -i "\\.go$" -i "\\.html$"  -i "\\.js$" .

    # change detected, stop/wait for the process to terminate.
    kill $PID
    wait $PID 2>/dev/null

done
