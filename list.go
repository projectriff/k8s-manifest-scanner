package scan

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pivotal/go-ape/pkg/furl"
)

func ListImages(res string, baseDir string) ([]string, error) {
	fmt.Printf("Scanning %s\n", res)
	contents, err := furl.Read(res, baseDir)
	if err != nil {
		return nil, err
	}

	imgs := []string{}

	docs := strings.Split(string(contents), "---\n")
	if runtime.GOOS == "windows" {
		// allow lines to end in LF or CRLF since either may occur
		d := strings.Split(string(contents), "---\r\n")
		if len(d) > len(docs) {
			docs = d
		}
	}
	for _, doc := range docs {
		if strings.TrimSpace(doc) != "" {
			y := make(map[string]interface{})
			err = yaml.Unmarshal([]byte(doc), &y)
			if err != nil {
				return nil, fmt.Errorf("error parsing resource file %s: %v", res, err)
			}

			visitImages(y, func(imageName string) {
				imgs = append(imgs, imageName)
			})
		}
	}

	return imgs, nil
}

func visitImages(y interface{}, visitor func(string)) {
	switch v := y.(type) {
	case map[string]interface{}:
		if val, ok := v["image"]; ok {
			if vs, ok := val.(string); ok {
				if !strings.HasPrefix(vs, "$") { // skip parameter usages
					visitor(vs)
				}
			}
		}

		if args, ok := v["args"]; ok {
			if ar, ok := args.([]interface{}); ok {
				for i, a := range ar {
					if a, ok := a.(string); ok {
						if strings.HasPrefix(a, "-") && strings.HasSuffix(a, "-image") && len(ar) > i+1 {
							if b, ok := ar[i+1].(string); ok {
								visitor(b)
							}
						}
					}
				}
			}
		}

		if val, ok := v["config"]; ok {
			if vs, ok := val.(string); ok {
				y := make(map[string]interface{})
				err := yaml.Unmarshal([]byte(vs), &y)
				if err == nil {
					visitImages(y, visitor)
				}
			}
		}

		if val, ok := v["template"]; ok {
			if vs, ok := val.(string); ok {
				// treat templates as lines each of which may contain YAML
				lines := strings.Split(vs, "\n")
				for _, line := range lines {
					y := make(map[string]interface{})
					err := yaml.Unmarshal([]byte(line), &y)
					if err == nil {
						visitImages(y, visitor)
					}
				}
			}
		}

		if parms, ok := v["parameters"]; ok {
			if pr, ok := parms.([]interface{}); ok {
				for _, p := range pr {
					if pmap, ok := p.(map[string]interface{}); ok {
						// if this parameter map has a "name" key which indicates an image and a "default" key with a
						// string value, treat the value as a possible image
						if name, ok := stringMapValue(pmap, "name"); ok {
							if strings.HasSuffix(name, "IMAGE") {
								if deflt, ok := stringMapValue(pmap, "default"); ok {
									visitor(deflt)
								}
							}
						}
					}
				}
			}
		}

		for key, val := range v {
			if strings.HasSuffix(key, "Image") || strings.HasSuffix(key, "-image") {
				if vs, ok := val.(string); ok {
					visitor(vs)
				}
			}
			visitImages(val, visitor)
		}
	case map[interface{}]interface{}:
		for _, val := range v {
			visitImages(val, visitor)
		}
	case []interface{}:
		for _, u := range v {
			visitImages(u, visitor)
		}
	default:
	}
}

func stringMapValue(m map[string]interface{}, key string) (string, bool) {
	if value, ok := m[key]; ok {
		if valueStr, ok := value.(string); ok {
			return valueStr, true

		}
	}
	return "", false
}
