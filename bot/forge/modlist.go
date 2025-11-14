package forge

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

func ParseModList(dir string) (ModList, error) {
	l, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var sli ModList
	var wg sync.WaitGroup
	for _, fi := range l {
		s := strings.Split(fi.Name(), ".")
		if s[len(s)-1] == "jar" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				f, err := os.OpenFile(fmt.Sprintf("%s/%s", dir, fi.Name()), os.O_RDONLY, 0644)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				fn, err := f.Stat()
				if err != nil {
					panic(err)
				}
				r, err := zip.NewReader(f, fn.Size())
				if err != nil {
					panic(err)
				}
				for _, v := range r.File {
					if v.Name == "mcmod.info" {
						rc, err := v.Open()
						if err != nil {
							panic(err)
						}
						buf := make([]byte, v.UncompressedSize64)
						if _, err := rc.Read(buf); err != nil {
							if err != io.EOF {
								panic(err)
							}
						}
						j, err := jsonvalue.Unmarshal(buf)
						if err != nil {
							panic(err)
						}
						// 奇葩1
						if v, err := j.GetArray("modList"); err == nil {
							j = v
						}
						for _, v := range j.ForRangeArr() {
							id, err := v.GetString("modid")
							if err != nil {
								panic(err)
							}
							ver, err := v.GetString("version")
							if err != nil {
								panic(err)
							}
							sli = append(sli, ModInfo{
								ModId:   id,
								Version: ver,
							})
						}
					}
				}
			}()
		}
	}
	wg.Wait()
	return sli, nil
}
