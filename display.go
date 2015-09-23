package doit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/digitalocean/godo"
)

const (
	// NSRoot is a configuration key that signifies this value is at the root.
	NSRoot = "doit"
)

// DisplayOutput displays an object or group of objects to a user. It
// checks to see what the output type should be.
func DisplayOutput(item interface{}, out io.Writer) error {
	output := DoitConfig.GetString(NSRoot, "output")
	if output == "" {
		output = "text"
	}

	switch output {
	case "json":
		return writeJSON(item, out)
	case "text":
		return writeText(item, out)
	default:
		return fmt.Errorf("unknown output type")
	}
}

func writeJSON(item interface{}, w io.Writer) error {
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	_, err = out.WriteTo(w)
	return err
}

func writeText(item interface{}, w io.Writer) error {
	switch item.(type) {
	case *godo.Account:
		a := item.(*godo.Account)
		outputAccount(a, w)
	case *godo.Action:
		i := item.(*godo.Action)
		outputActions([]godo.Action{*i}, w)
	case []godo.Action:
		outputActions(item.([]godo.Action), w)
	case *godo.Domain:
		outputZone(item.(*godo.Domain), w)
	case []godo.Domain:
		outputDomains(item.([]godo.Domain), w)
	case *godo.DomainRecord:
		i := item.(*godo.DomainRecord)
		outputRecords([]godo.DomainRecord{*i}, w)
	case []godo.DomainRecord:
		outputRecords(item.([]godo.DomainRecord), w)
	case *godo.Droplet:
		d := item.(*godo.Droplet)
		outputDroplets([]godo.Droplet{*d}, w)
	case []godo.Droplet:
		outputDroplets(item.([]godo.Droplet), w)
	case *godo.Image:
		i := item.(*godo.Image)
		outputImages([]godo.Image{*i}, w)
	case []godo.Image:
		outputImages(item.([]godo.Image), w)
	case *godo.Kernel:
		i := item.(*godo.Kernel)
		outputKernels([]godo.Kernel{*i}, w)
	case []godo.Kernel:
		outputKernels(item.([]godo.Kernel), w)
	case *godo.Key:
		i := item.(*godo.Key)
		outputSSHKeys([]godo.Key{*i}, w)
	case []godo.Key:
		outputSSHKeys(item.([]godo.Key), w)

	case []godo.Region:
		outputRegions(item.([]godo.Region), w)
	case []godo.Size:
		outputSizes(item.([]godo.Size), w)
	}

	return nil
}

func outputAccount(account *godo.Account, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "Email\tDroplet Limit\tEmail Verified\tUUID\t")
	fmt.Fprintf(w, "")
	fmt.Fprintf(w, "%s\t%d\t%t\t%s\t\n", account.Email, account.DropletLimit, account.EmailVerified, account.UUID)
	fmt.Fprintln(w)
	w.Flush()
}

func outputActions(list []godo.Action, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "ID\tStatus\tType\tStarted At\tCompleted At\tResource ID\tResource Type\tRegion")

	for _, a := range list {
		fmt.Fprintf(w, "")
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
			a.ID, a.Status, a.Type, a.StartedAt, a.CompletedAt, a.ResourceID, a.ResourceType, a.RegionSlug)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func outputDroplets(list []godo.Droplet, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "ID\tName\tPublic IPv4\tMemory\tVCPUs\tDisk\tRegion\tImage\tStatus")

	for _, d := range list {
		ip := extractDropletPublicIP(&d)
		image := fmt.Sprintf("%s %s", d.Image.Distribution, d.Image.Name)
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t%d\t%s\t%s\t%s\n",
			d.ID, d.Name, ip, d.Memory, d.Vcpus, d.Disk, d.Region.Slug, image, d.Status)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func outputImages(list []godo.Image, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "ID\tName\tType\tDistribution\tSlug\tPublic\tMin Disk")

	for _, i := range list {
		publicStatus := false
		if i.Public {
			publicStatus = true
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%t\t%d\n",
			i.ID, i.Name, i.Type, i.Distribution, i.Slug, publicStatus, i.MinDiskSize)

	}
	fmt.Fprintln(w)
	w.Flush()
}

func outputKernels(list []godo.Kernel, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "ID\tName\tVersion")

	for _, k := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\n", k.ID, k.Name, k.Version)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func outputRegions(list []godo.Region, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "Slug\tName\tAvailable")

	for _, r := range list {
		fmt.Fprintf(w, "%s\t%s\t%t\n", r.Slug, r.Name, r.Available)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func outputSizes(list []godo.Size, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "Slug\tMemory\tVcpus\tDisk\tPrice Monthly\tPrice Hourly")

	for _, s := range list {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%0.2f\t%f\n",
			s.Slug, s.Memory, s.Vcpus, s.Disk, s.PriceMonthly, s.PriceHourly)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func outputSSHKeys(list []godo.Key, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "ID\tName\tFingerprint")

	for _, s := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\n",
			s.ID, s.Name, s.Fingerprint)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func outputZone(domain *godo.Domain, out io.Writer) {
	fmt.Fprintln(out, domain.ZoneFile)
}

func outputDomains(list []godo.Domain, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "Name")

	for _, d := range list {
		fmt.Fprintf(w, "%s\n", d.Name)
	}

	fmt.Fprintln(w)
	w.Flush()
}

func outputRecords(list []godo.DomainRecord, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "ID\tType\tName\tData\tPriority\tPort\tWeight")

	for _, d := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%d\t%d\n", d.ID, d.Type, d.Name, d.Data,
			d.Priority, d.Port, d.Weight)
	}

	fmt.Fprintln(w)
	w.Flush()
}
