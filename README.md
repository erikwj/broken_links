# Brokenlinks

This application can be used to validate markdown files for brokenlinks. Most IDEs also have the option to do this, but these will not validate web urls. Weblinks are captured and printed out in such a way that they can be opened via a terminal. 

Usage

Use the `-h` flag for printing the documentation
``` shell
./brokenlinks -h
```



```
# Running on linux
./brokenlinks --dir /dir/with/markdowns
# Running on Mac intel
./brokenlinks-amd64 --dir . | sh # to open directly in your local browser

# Running on Mac M1+
./brokenlinks-arm64 --dir . | sh # to open directly in your local browser
```

```
go run main.go --dir /path/to/rstfiles --ext .rst
```