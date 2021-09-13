package IOUtils

import (
	"fmt"
	"io"
	"log"
)

func Close(closeable io.Closer)  {
	err := closeable.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to close %#v, error is: %s", closeable, err.Error()))
	}
}

func Write(writeable io.Writer, data []byte)  {
	size, err := writeable.Write(data)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to write %#v, error is: %s", writeable, err.Error()))
	} else {
		log.Println(fmt.Sprintf("Write %d bytes to %#v", size, writeable))
	}
}