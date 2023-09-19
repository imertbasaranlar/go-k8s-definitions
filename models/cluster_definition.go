package models

type ClusterDefinition struct {
	Name       string
	Pods       []*PodDefinition
	Nodes      []*NodeDefinition
	Namespaces []*NamespaceDefinition
	Endpoints  map[string][]*SubsetAddressDefinition
}

type PodDefinition struct {
	Name string
	IP   string
}

type NodeDefinition struct {
	Name string
}

type NamespaceDefinition struct {
	Name string
}

type SubsetAddressDefinition struct {
	IP string
}
