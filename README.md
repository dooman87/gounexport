# Gounexport #

Gounexport is a tool for finding exported symbols that is not used outside of the package and unexporting them by renaming to lowercase.

By default, it's working in safe mode and only printing out result without renaming. Use -rename option to do actual renaming.

```
Usage: gounexport [OPTIONS] package
  -exclude string
        File with exlude patterns for objects that shouldn't be unexported. Each pattern should be started at new line. Default pattern is Test* to exclude tests methods.
  -out string
        Output file. If not set then stdout will be used
  -rename
        If set, then all defenitions that will be determined as unused will be renamed in files
  -verbose
        Turning on verbose mode
```

# History #

The app was originally developed as part of fifth [golang-challenge](http://golang-challenge.com/go-challenge5).
Surprisingly, I found that the solution is doing the job quite sufficient. However, the code was not looking good and
required a lot of refactoring, cause it was my first go app and I was running out of time. Finally, I did it and proud of it :)
