package manager

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/shirokurostone/hosts-manager/hosts"
)

func editTempFile(fileBody string) (string, error) {
	f, err := ioutil.TempFile("", "hosts-manager")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := ioutil.WriteFile(f.Name(), []byte(fileBody), 644); err != nil {
		return "", err
	}

	editor := os.Getenv("EDITOR")
	proc := exec.Command(editor, f.Name())
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	err = proc.Run()
	if err != nil {
		return "", err
	}
	proc.Wait()

	bytes, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func NewGroup(config *Config, name string) error {

	body, err := editTempFile("")
	if err != nil {
		return err
	}
	if body == "" {
		return nil
	}

	result, err := hosts.Parse(strings.NewReader(body))
	if err != nil {
		return err
	}
	if !result.CheckSyntax() {
		return errors.New("syntax error")
	}

	config.SetHostsGroup(HostsGroup{Name: name, Body: body, IsActive: false})
	return nil
}

func EditGroup(config *Config, name string) error {

	if !config.ExistsHostsGroup(name) {
		return fmt.Errorf("group '%s' not found", name)
	}

	f, err := ioutil.TempFile("", "hosts-manager")
	if err != nil {
		return err
	}
	defer f.Close()

	group := config.GetHostsGroup(name)
	body, err := editTempFile(group.Body)
	if err != nil {
		return err
	}

	if group.Body == body {
		return nil
	} else if body == "" {
		config.RemoveHostsGroup(name)
		return nil
	}

	result, err := hosts.Parse(strings.NewReader(body))
	if err != nil {
		return err
	}
	if !result.CheckSyntax() {
		return errors.New("syntax error")
	}

	group.Body = body
	config.SetHostsGroup(*group)
	return nil
}

func RemoveGroup(config *Config, names []string) error {

	for _, n := range names {
		if !config.ExistsHostsGroup(n) {
			return fmt.Errorf("group '%s' not found", n)
		}
	}

	for _, n := range names {
		config.RemoveHostsGroup(n)
	}

	return nil
}

func ShowGroup(config *Config, names []string) error {

	for _, n := range names {
		if !config.ExistsHostsGroup(n) {
			return fmt.Errorf("group '%s' not found", n)
		}
	}

	if len(names) == 0 {
		for _, group := range config.Groups {
			if group.IsActive {
				fmt.Print(group.Body)
			}
		}
	} else {
		for _, n := range names {
			group := config.GetHostsGroup(n)
			fmt.Print(group.Body)
		}
	}

	return nil
}

func ListGroup(config *Config) error {

	for _, group := range config.Groups {
		if group.IsActive {
			fmt.Printf("%s\t%s\n", group.Name, "active")
		} else {
			fmt.Printf("%s\t%s\n", group.Name, "inactive")
		}
	}
	return nil
}

func ActivateGroup(config *Config, names []string) error {

	for _, n := range names {
		if !config.ExistsHostsGroup(n) {
			return fmt.Errorf("group '%s' not found", n)
		}
	}

	for _, n := range names {
		group := config.GetHostsGroup(n)
		group.IsActive = true
		config.SetHostsGroup(*group)
	}

	return nil
}

func DeactivateGroup(config *Config, names []string) error {

	for _, n := range names {
		if !config.ExistsHostsGroup(n) {
			return fmt.Errorf("group '%s' not found", n)
		}
	}

	for _, n := range names {
		group := config.GetHostsGroup(n)
		group.IsActive = false
		config.SetHostsGroup(*group)
	}

	return nil
}

func ApplyGroup(config *Config, hostsfile string, isDryRun bool) error {

	target := []HostsGroup{}
	for _, group := range config.Groups {
		if group.IsActive {
			target = append(target, group)
		}
	}

	f, err := ioutil.TempFile("", "hosts-manager")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadFile(hostsfile)
	if err != nil {
		return err
	}
	origin := string(data)
	sb := strings.Builder{}
	if err := ApplyHosts(strings.NewReader(origin), io.MultiWriter(f, &sb), target); err != nil {
		return err
	}

	proc := exec.Command("diff", "-u", hostsfile, f.Name())
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	proc.Run()
	proc.Wait()

	if !isDryRun {
		if err := ioutil.WriteFile(hostsfile, []byte(sb.String()), 644); err != nil {
			return err
		}
	}

	return nil
}

var header = "\n##### managed by hosts-manager : start #####\n"
var footer = "\n##### managed by hosts-manager : end #####\n"

func ApplyHosts(reader io.Reader, writer io.Writer, groups []HostsGroup) error {

	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	data := string(bs)
	headerIndex := strings.Index(data, header)
	footerIndex := strings.Index(data, footer)
	if headerIndex == -1 && footerIndex == -1 {
		headerIndex = len(data)
		footerIndex = len(data)
	} else if headerIndex == -1 || footerIndex == -1 || footerIndex < headerIndex {
		return errors.New("hostsfile format error")
	} else {
		footerIndex += len(footer)
	}

	w := bufio.NewWriter(writer)
	defer w.Flush()

	w.WriteString(data[0:headerIndex])
	w.WriteString(header)

	for _, g := range groups {
		w.WriteString("\n")
		w.WriteString(g.Body)
		w.WriteString("\n")
	}

	w.WriteString(footer)

	w.WriteString(data[footerIndex:])
	return nil

}
