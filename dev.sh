
trap ctrl_c INT

function ctrl_c() {
    printf "Exiting\n"
    kill $PID
    exit 1
}

mkdir -p bin

while true; do
    rm ./assets/styledist.css
    cd dependencies
    npx tailwindcss -i ../assets/style.css -o ../assets/styledist.css
    cd ..
    go build -o bin/server src/* 
    ./bin/server &
    PID=$!
    echo "Started with PID $PID in $PWD\n"

    fswatch -1 -r -e ".*" -i "\\.go$" -i "\\.html$"  -i "\\.js$" .

    kill $PID
    wait $PID 2>/dev/null
done
