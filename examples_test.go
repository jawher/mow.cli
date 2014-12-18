package cli_test

import (
	"fmt"
	"os"
)

func Example_greet() {
	app := App("greet", "Greet")
	app.Spec = "[NAME]"
	name := app.StringArg("NAME", "stranger", "Your name", &OptExtra{EnvVar: "USER"})
	app.Action = func() {
		fmt.Printf("Hello %s\n", *name)
	}
	app.Run(os.Args)
}

func Example_cp() {
	cp := App("cp", "Copy files around")
	cp.Spec = "[-R [-H | -L | -P]] [-fi | -n] SRC... DST"

	var (
		recursive = cp.BoolOpt("R", false, "copy src files recursively", nil)

		followSymbolicCL   = cp.BoolOpt("H", false, "If the -R option is specified, symbolic links on the command line are followed.  (Symbolic links encountered in the tree traversal are not followed.)", nil)
		followSymbolicTree = cp.BoolOpt("L", false, "If the -R option is specified, all symbolic links are followed.", nil)
		followSymbolicNo   = cp.BoolOpt("P", true, "If the -R option is specified, no symbolic links are followed.  This is the default.", nil)

		force       = cp.BoolOpt("f", false, "If the destination file cannot be opened, remove it and create a new file, without prompting for confirmation regardless of its permissions.  (The -f option overrides any previous -n option.)", nil)
		interactive = cp.BoolOpt("i", false, "Cause cp to write a prompt to the standard error output before copying a file that would overwrite an existing file.  If the response from the standard input begins with the character `y' or `Y', the file copy is attempted.  (The -i option overrides any previous -n option.)", nil)
		noOverwrite = cp.BoolOpt("f", false, "Do not overwrite an existing file.  (The -n option overrides any previous -f or -i options.)", nil)
	)

	var (
		src = cp.StringsArg("SRC", nil, "The source files to copy", nil)
		dst = cp.StringsArg("DST", nil, "The destination directory", nil)
	)

	cp.Action = func() {
		fmt.Printf(`copy:
	SRC: %v
	DST: %v
	recursive: %v
	follow links (CL, Tree, No): %v %v %v
	force: %v
	interactive: %v
	no overwrite: %v`,
			*src, *dst, *recursive,
			*followSymbolicCL, *followSymbolicTree, *followSymbolicNo,
			*force,
			*interactive,
			*noOverwrite)
	}

	cp.Run(os.Args)
}

func Example_docker() {
	docker := App("docker", "A self-sufficient runtime for linux containers")

	docker.Command("run", "Run a command in a new container", func(cmd *Cmd) {
		cmd.Spec = "[-d|--rm] IMAGE [COMMAND [ARG...]]"

		var (
			detached = cmd.BoolOpt("d detach", false, "Detached mode: run the container in the background and print the new container ID", nil)
			rm       = cmd.BoolOpt("rm", false, "Automatically remove the container when it exits (incompatible with -d)", nil)
			memory   = cmd.StringOpt("m memory", "", "Memory limit (format: <number><optional unit>, where unit = b, k, m or g)", nil)
		)

		var (
			image   = cmd.StringArg("IMAGE", "", "", nil)
			command = cmd.StringArg("COMMAND", "", "The command to run", nil)
			args    = cmd.StringsArg("ARG", nil, "The command arguments", nil)
		)

		cmd.Action = func() {
			how := ""
			switch {
			case *detached:
				how = "detached"
			case *rm:
				how = "rm after"
			default:
				how = "--"
			}
			fmt.Printf("Run image %s, command %s, args %v, how? %v, mem %s", *image, *command, *args, how, *memory)
		}
	})

	docker.Command("pull", "Pull an image or a repository from the registry", func(cmd *Cmd) {
		cmd.Spec = "[-a] NAME"

		all = cmd.BoolOpt("a all-tags", false, "Download all tagged images in the repository", nil)

		name = cmd.StringArg("NAME", "", "Image name (optionally NAME:TAG)", nil)

		cmd.Action = func() {
			if *all {
				fmt.Printf("Download all tags for image %s", *name)
				return
			}
			fmt.Printf("Download image %s", *name)
		}
	})

	docker.Run(os.Args)
}
