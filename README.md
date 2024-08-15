# go-hrx

This package implements a decoder for [hrx file format](https://github.com/google/hrx),
examples from the spec repo serve as test cases.

## API

`hrx` package exports a single function - `OpenReader` to open and parse a file. Returned
object implements a number of interfaces from `io/fs` package that allow to work with
it as with a folder.

```
// open an hrx file
reader, err := hrx.OpenReader(fullHrxPath)

if err != nil {
	panic(err) // always handle your errors properly
}

// print contents of parsed virtual filesystem
err = fs.WalkDir(reader, ".", func(p string, d fs.DirEntry, err2 error) error {
	fmt.Println(p)
})

if err != nil {
	panic(err) // panic is not a proper way
}
```

## Author

Dmitrii Petrov / dpetroff@gmail.com

## License

Apache license
