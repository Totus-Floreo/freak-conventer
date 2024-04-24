[![Tests](https://github.com/Totus-Floreo/freak-conventer/actions/workflows/tests.yml/badge.svg)](https://github.com/Totus-Floreo/freak-conventer/actions/workflows/tests.yml)

# Project Title: Unix Time Converter

This project is a Go module that provides functionality to convert Go struct fields of type `time.Time` to Unix time (the number of seconds elapsed since January 1, 1970 UTC). The module uses reflection to traverse the struct fields and perform the conversion.

## Getting Started

To use this module, you need to import it in your Go project:

```go
import "github.com/Totus-Floreo/freak-conventer"
```

## Usage

The main function of this module is `ConvertToUnixTime(v interface{})`. This function takes a struct as an argument and returns a map where the keys are the struct field names and the values are the corresponding field values. If a field is of type `time.Time`, the value will be the Unix time equivalent.

Here is an example of how to use it:

```go
type MyStruct struct {
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

s := MyStruct{
    Name:      "Test",
    CreatedAt: time.Now(),
}

result, err := freak_conventer.ConvertToUnixTime(s)
if err != nil {
    log.Fatal(err)
}

fmt.Println(result) // map[created_at:1713840688 name:Test]
```

In the output map, the `created_at` field will be represented as Unix time.

## Contributing

Contributions are welcome. Please feel free to submit a pull request.

## License

This project is licensed under the MIT License.
