package scan

import (
	"strings"

	. "github.com/dprotaso/go-yit"
	"gopkg.in/yaml.v3"
)

func SearchImageNodes(doc *yaml.Node) []*yaml.Node {
	return FromIterators(
		MatchImageKey(doc),
		MatchArgsMap(doc),
		MatchTemplateDefaults(doc),
	).ToArray()
}

func MatchImageKey(node *yaml.Node) Iterator {
	return FromNode(node).
		RecurseNodes().
		Filter(WithKind(yaml.MappingNode)).
		ValuesForMap(
			// Key Predicate
			Union(
				WithValue("image"),
				WithSuffix("Image"),
				WithSuffix("-image"),
			),
			// Value Predicate
			StringValue,
		).
		Filter(Negate(WithPrefix("$")))
}

func MatchTemplateDefaults(node *yaml.Node) Iterator {
	return FromNode(node).
		RecurseNodes().
		Filter(WithKind(yaml.MappingNode)).
		ValuesForMap(
			// Key Predicate
			WithValue("parameters"),
			// Value Prediate
			WithKind(yaml.SequenceNode),
		).
		Values(). // Unpack sequences
		Filter(
			Intersect(
				WithKind(yaml.MappingNode),
				WithMapKeyValue(WithValue("name"), WithSuffix("IMAGE")),
				WithMapKey("default"),
			),
		).
		ValuesForMap(
			// Key Predicate
			WithValue("default"),
			// Value Predicate
			StringValue,
		)
}

func MatchArgsMap(node *yaml.Node) Iterator {
	return FromNode(node).
		RecurseNodes().
		Filter(WithKind(yaml.MappingNode)).
		ValuesForMap(
			// Key Predicate
			WithValue("args"),
			// Value Predicate
			WithKind(yaml.SequenceNode),
		).
		Iterate(imageFlagValues).
		Filter(StringValue)
}

var imageFlagValues = func(next Iterator) Iterator {
	var content []*yaml.Node

	return func() (node *yaml.Node, ok bool) {
		for {
			for len(content) >= 2 {
				arg := content[0]
				node = content[1]

				// not all arg sequences are (flag, value) pairs
				content = content[1:]

				if ok = strings.HasPrefix(arg.Value, "-") &&
					strings.HasSuffix(arg.Value, "-image"); ok {
					return
				}
			}

			var parent *yaml.Node
			for parent, ok = next(); ok; parent, ok = next() {
				if parent.Kind == yaml.SequenceNode && len(parent.Content) > 0 {
					break
				}
			}

			if !ok {
				return
			}

			content = parent.Content
		}
	}
}
