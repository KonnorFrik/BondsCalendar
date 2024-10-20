### Bonds Calendar ###
Created for calculate coupon payment dates

Dependencies:
    github.com/gbin/goncurses

Features:
    - Use terminal for manage your bonds
        - Creating a multiply bonds and see all payments at graph
        - Delete any bonds if you don't like it
        - Save them into json file
        - Load bonds from previously saved json file
        - See info about all appended bonds

Movement:
    - 'h' - Open a info window with keys for this 
    - 'q' - Exit from programm or close opened sub-window
    - ':' - Start typing a command in the terminal
    In main window(graph):
        - '>' - Show payment graph for next year
        - '<' - Show payment graph for previous year
    In any scrollable window:
        - 'w' - Scroll down
        - 's' - Scroll up

Terminal commands:
    Pattern: [COMMAND] [ARGUMENTS]
        COMMAND - always required
        ARGUMENTS - in most cases optional (programm will ask for them)

    - "help [command]" - Show all commands and their info
    - "new" - Start filling out the form for create a new bond
    - "list" - Show all bonds with indices and some info
    - "delete [index]" - Delete bond by it index
    - "save [filename]" - Save bonds info into json file
    - "load [filename]" - Load bonds info from json file
