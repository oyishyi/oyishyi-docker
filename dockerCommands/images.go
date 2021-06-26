package dockerCommands

import (
	"fmt"
	"github.com/oyishyi/docker/container"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"text/tabwriter"
)

func GetImages() {
	files, err := os.ReadDir(container.ImagesURL)
	if err != nil {
		logrus.Errorf("get files fails: %v", err)
		return
	}
	var images []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tar") {
			images = append(images, strings.TrimSuffix(file.Name(), ".tar"))
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tDIR\n")
	for _, image := range images {
		fmt.Fprintf(w, "%s\t%s\n", image, container.ImagesURL)
	}
	if err := w.Flush(); err != nil {
		logrus.Errorf("flush error: %v", err)
		return
	}
}