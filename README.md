# scrub

Recursively set specific struct fields to their zero values

Possible use cases:
  - scrubbing sensitive data from structs before logging
  - comparing structs with noisy fields (timestamps, random values, etc) for testing or diagnostic purposes

## Examples

### Using struct tags

```go
package main

import (
  "fmt"

  "github.com/acj/scrub"
)

type User struct {
  Name string
  Age  int    `scrub:"true"`
}

func main() {
  user := User{
    Name: "Wall-E",
    Age:  22,
  }
  scrub.TaggedFields(&user)
  fmt.Printf("%+v\n", user) // {Name:Wall-E Age:0}
}
```

### Using named fields (blocklist)

```go
package main

import (
  "fmt"

  "github.com/acj/scrub"
)

type User struct {
  Name string
  Age  int
}

func main() {
  user := User{
    Name: "Wall-E",
    Age:  22,
  }
  scrub.NamedFields(&user, "Age")
  fmt.Printf("%+v\n", user) // {Name:Wall-E Age:0}
}
```

## License

MIT