package inspec

func buildInspecCommand(profiles []string) []string {
	posArgs := append([]string{"exec"}, profiles...)
	posArgs = append(posArgs, "--json-config=-")

	var cmdargs []string
	cmdargs = append([]string{"inspec"}, posArgs...)
	return cmdargs
}
