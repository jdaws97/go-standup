# STANDUP

Standup is a note taking tool for your standup meetings!

It's a CLI tool made for Linux with GO

Clone the repository, and build the executable with GO:

```
go build standup.go
```

To add it as an alias (add the path to the executable):

```
export standup='$PATH_TO_EXECUTABLE/./standup'
```

EXAMPLE USAGE:

```
standup -c in-progress -s "Working on ticket #1234"
standup --open -d 7
standup --open 
standup --config
standup -s "Need to message Jim" --> this defaults to being added to the notes category
```

List of all possible arguments:

--open -> this will open the standup notes

--config -> this will open the config file to change it

-c -> options: DONE, IN-PROGRESS, BLOCKERS, NOTES (you can change your categories in the config)

-s -> just write your sentence in double quotes
-d -> you can open your notes from however many days ago (if it exists)
