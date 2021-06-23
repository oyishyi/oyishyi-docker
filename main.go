package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var app = &cli.App{
		Description: "时值 golang 战国年代，冉冉升起的一颗巨星，其名为...\n         ::::::::   :::   :::  :::::::::::   ::::::::   :::    :::  :::   :::  ::::::::::: \n        :+:    :+:  :+:   :+:      :+:      :+:    :+:  :+:    :+:  :+:   :+:      :+:     \n        +:+    +:+   +:+ +:+       +:+      +:+         +:+    +:+   +:+ +:+       +:+     \n        +#+    +:+    +#++:        +#+      +#++:++#++  +#++:++#++    +#++:        +#+     \n        +#+    +#+     +#+         +#+             +#+  +#+    +#+     +#+         +#+     \n        #+#    #+#     #+#         #+#      #+#    #+#  #+#    #+#     #+#         #+#     \n         ########      ###     ###########   ########   ###    ###     ###     ########### ",
		Commands: []*cli.Command{
			&initCommand,
			&runCommand,
			&commitCommand,
			&psCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
