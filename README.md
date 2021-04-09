![golangci-lint](https://github.com/godzie44/go-xls/workflows/golangci-lint/badge.svg) 
![tests](https://github.com/godzie44/go-xls/workflows/tests/badge.svg) 
[![Coverage Status](https://coveralls.io/repos/github/godzie44/go-xls/badge.svg?branch=master)](https://coveralls.io/github/godzie44/go-xls?branch=master)


# GO-xls
GO wrapper for **[libxls](https://github.com/libxls)**. Read old Excel files (.xls) from GO.

## Install

1) install **[libxls](https://github.com/libxls)**
2) add CGO_LDFLAGS=-lxlsreader to go build / go test command.

See /example/Dockerfile as example of installing required dependencies.

## Example

```go
	wb, err := xls.OpenFile(someFile, "UTF-8")
	if err != nil {
		log.Fatal(err)
	}
	defer wb.Close()

	sheet, err := wb.OpenWorkSheet(0)
	if err != nil {
		log.Fatal(err)
	}
	defer sheet.Close()

	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
			fmt.Println(cell.Value.String())
		}
	}
```

Or more specialized cell parsing:

```go
	//...

	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
                        switch v := cell.Value.(type) {
                        case *BlankValue:
                            fmt.Printf("%T no val \n", v)
                        case *FloatValue:
                            fmt.Printf("%T %f \n", v, v.Val)
                        case *BoolValue:
                            fmt.Printf("%T %t \n", v, v.Val)
                        case *ErrValue:
                            fmt.Printf("%T %d \n", v, v.Code)
                        case *StringValue:
                            fmt.Printf("%T %s \n", v, v.Val)
                        case *UnknownValue:
                            fmt.Printf("%T %s \n", v, v.Val)
                        }
		    }
		}
	}
```

Also, use `make example` command for run example project (converting xls table to html).