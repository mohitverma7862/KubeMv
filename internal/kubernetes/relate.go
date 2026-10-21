package kubernetes

import "strings"

// TreeNode is one selectable line in the ownership tree.
type TreeNode struct {
	Kind      string
	Name      string
	Namespace string
	Depth     int
}

// PublicRow drops secret values. Secret rows keep type and key count only.
func PublicRow(row ResourceRow) ResourceRow {
	if row.Kind != "Secret" {
		return row
	}
	extra := map[string]string{}
	if row.Extra != nil {
		if value := row.Extra["type"]; value != "" {
			extra["type"] = value
		}
		if value := row.Extra["keys"]; value != "" {
			extra["keys"] = value
		}
	}
	row.Extra = extra
	row.Containers = nil
	if len(row.Labels) > 0 {
		labels := map[string]string{}
		for key, value := range row.Labels {
			if safeMetadata(key) && safeMetadata(value) {
				labels[key] = value
			}
		}
		row.Labels = labels
	}
	return row
}

func safeMetadata(value string) bool {
	text := strings.ToLower(value)
	if strings.Contains(text, "begin ") || strings.Contains(text, "private key") || strings.Contains(text, "bearer ") {
		return false
	}
	return true
}

// BuildTree walks owner references from the leaf toward the root.
// lookup resolves the next owner. A missing owner stops the walk and still appears in the tree.
func BuildTree(leaf ResourceRow, lookup func(OwnerRef) (ResourceRow, bool)) []TreeNode {
	leaf = PublicRow(leaf)
	if leaf.Kind == "" && leaf.Name == "" {
		return nil
	}
	upward := []OwnerRef{{Kind: leaf.Kind, Name: leaf.Name, Namespace: leaf.Namespace}}
	current := leaf
	seen := map[string]bool{ownerKey(upward[0]): true}
	for len(upward) < 8 && len(current.Owners) > 0 {
		owner := current.Owners[0]
		if owner.Name == "" || seen[ownerKey(owner)] {
			break
		}
		seen[ownerKey(owner)] = true
		upward = append(upward, owner)
		parent, ok := lookup(owner)
		if !ok {
			break
		}
		current = PublicRow(parent)
	}
	nodes := make([]TreeNode, 0, len(upward)+1+len(leaf.Containers))
	for index := len(upward) - 1; index >= 0; index-- {
		ref := upward[index]
		nodes = append(nodes, TreeNode{Kind: ref.Kind, Name: ref.Name, Namespace: ref.Namespace, Depth: len(upward) - 1 - index})
	}
	if len(leaf.Containers) > 0 && leaf.Kind != "Secret" {
		depth := len(nodes)
		nodes = append(nodes, TreeNode{Kind: "Containers", Depth: depth})
		for _, container := range leaf.Containers {
			nodes = append(nodes, TreeNode{Kind: "Container", Name: container.Name, Depth: depth + 1})
		}
	}
	return nodes
}

// FormatTree renders the compact ownership tree. cursor highlights one line when it is in range.
func FormatTree(nodes []TreeNode, cursor int) string {
	if len(nodes) == 0 {
		return "No owner chain."
	}
	lines := make([]string, len(nodes))
	for i, node := range nodes {
		label := node.Kind + " (" + node.Name + ")"
		if node.Kind == "Containers" {
			label = "Containers"
		}
		if node.Kind == "Container" {
			label = node.Name
		}
		if node.Depth > 0 {
			label = strings.Repeat("   ", node.Depth-1) + "└─ " + label
		}
		if i == cursor {
			label = "> " + label
		}
		lines[i] = label
	}
	return strings.Join(lines, "\n")
}

func ownerKey(ref OwnerRef) string {
	return ref.Kind + "/" + ref.Namespace + "/" + ref.Name
}
