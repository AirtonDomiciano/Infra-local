package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// UpsertKey atualiza KEY=VALUE se existir, ou adiciona no final se não existir.
func UpsertKey(path string, key string, value string) error {
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	var (
		lines  []string
		found  bool
		prefix = key + "="
	)

	sc := bufio.NewScanner(in)
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)

		// mantém comentários/linhas vazias
		if trim == "" || strings.HasPrefix(trim, "#") {
			lines = append(lines, line)
			continue
		}

		if strings.HasPrefix(trim, prefix) {
			lines = append(lines, fmt.Sprintf("%s%s", prefix, value))
			found = true
		} else {
			lines = append(lines, line)
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s%s", prefix, value))
	}

	out := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(out), 0644)
}
