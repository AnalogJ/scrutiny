package shell

func Create() Interface {
	return new(localShell)
}
