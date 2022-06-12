# Simple Notes/CMS/Wiki/Pages System.

Codename: designator / peasy

This is a simple CMS-like system based on simple Go code.

There is no attempt to fix HTML or CSS. HTML is created in Go templates. 
CSS is Bootstrap 5.

It can be used for knowledge storage and sharing.



# Config

config.toml contains the servers portnumber.


# Dependencies

Go
Bootstrap
fswatch


CSS / Layout

Page
    Topline
        Text
        Link
    Header
        Menu
        Logo
        Search
    Intro
        Hero
        Title
        Text
    Sections
        Section
            Text
            Grid
                Cell
                    Text
                    Image                    
    Footer
    LowLine



# Development

./dev.sh script watches for changes of .go , .html and .js files.
If change is detected, the go server is rebuild and restarted.


