def "do run" [] {
    echo "TASK DO RUN..."
    do all
    echo "RUNNING..."
    ./bin/main.exe
}

def "do build" [] {
    echo "TASK DO BUILD..."
    cd bin
    go build ../cmd/app/main.go
}
    
def "undo build" [] {
    echo "TASK UNDO BUILD..."
    rm bin/main.exe
}

def "do templ" [] {
    echo "TASK DO TEMPL..."
    templ generate
}
    
def "undo templ" [] {
    echo "TASK UNDO TEMPL..."
    cd ./internal/handlers
    ls *.go | each { |x| rm $x.name }
}

def "do css" [] {
    echo "TASK DO CSS..."
    cd internal/css
    npx tailwindcss -i input.css -o output.css
}

def "undo css" [] {
    echo "TASK UNDO CSS..."
    rm ./internal/static/css/output.css
}

def "do all" [] {
    # do css
    do templ
    do build
}

def "undo all" [] {
    undo css
    undo templ
    undo build
}
