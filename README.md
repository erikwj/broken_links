# Brokenlinks

This application can be used to validate markdown files for brokenlinks. Most IDEs also have the option to do this, but these will not validate web urls. Weblinks are captured and printed out in such a way that they can be opened via a terminal. 

Usage

Use the `-h` flag for printing the documentation
``` shell
./brokenlinks -h
```



```
./brokenlinks --dir /dir/with/markdowns
./brokenlinks --dir /dir/bla | sh
```

```
go run main.go --dir /path/to/rstfiles --ext .rst
```