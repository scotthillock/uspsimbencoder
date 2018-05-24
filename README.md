## Introduction ##

This is a package for GO which can be used to encode a numeric IMB string to an ADFT version

## Example ##

```go
package main

import (
	"fmt"
	"github.com/scotthillock/uspsimbencoder"
)

func main() {
	s := "5337977723499454492851135759461"
	encoded := imbencode.Encode(s)
	fmt.Println(encoded)
}
```

## Credits ##
Original conversion from the TCPDF project to a standalone PHP class was done by Ian Simpson

see http://www.tcpdf.org and https://gist.github.com/IanSimpson/96355bb98d396c4c49cb