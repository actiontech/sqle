package conf

import (
	"bytes"
	"io"
	"os"
	"sort"
)

// WriteConfigFile saves the configuration representation to a file.
// The desired file permissions must be passed as in os.Open.
// The header is a string that is saved as a comment in the first line of the file.
func (c *ConfigFile) WriteConfigFile(fname string, perm uint32, header string, firstSections []string) (err error) {
	var file *os.File

	if file, err = os.Create(fname); err != nil {
		return err
	}
	if err = c.Write(file, header, firstSections); err != nil {
		return err
	}

	return file.Close()
}

// WriteConfigBytes returns the configuration file.
func (c *ConfigFile) WriteConfigBytes(header string) (config []byte) {
	buf := bytes.NewBuffer(nil)

	c.Write(buf, header, []string{})

	return buf.Bytes()
}

// Writes the configuration file to the io.Writer.
func (c *ConfigFile) Write(writer io.Writer, header string, firstSections []string) (err error) {
	buf := bytes.NewBuffer(nil)

	if header != "" {
		if _, err = buf.WriteString("# " + header + "\n"); err != nil {
			return err
		}
	}

	//sort sections
	sectionSort := make([]string, 0)
	for section := range c.data {
		found := false
		for _, first := range firstSections {
			if first == section {
				found = true
			}
		}
		if !found {
			sectionSort = append(sectionSort, section)
		}
	}
	sort.Strings(sectionSort)
	sectionSort = append(firstSections, sectionSort...)

	for _, section := range sectionSort {
		sectionmap, exist := c.data[section]
		if !exist {
			continue
		}
		if section == DefaultSection && len(sectionmap) == 0 {
			continue // skip default section if empty
		}
		if _, err = buf.WriteString("[" + section + "]\n"); err != nil {
			return err
		}

		optionSort := make([]string, 0)
		for option := range sectionmap {
			optionSort = append(optionSort, option)
		}

		sort.Strings(optionSort)

		for _, option := range optionSort {
			value, exist := sectionmap[option]
			if !exist {
				continue
			}
			if _, err = buf.WriteString(option + "=" + value + "\n"); err != nil {
				return err
			}
		}
		if _, err = buf.WriteString("\n"); err != nil {
			return err
		}
	}

	buf.WriteTo(writer)

	return nil
}
