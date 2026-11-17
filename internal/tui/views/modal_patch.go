// Package views provides patch preview modal formatting.
package views

import (
	"fmt"
	"strings"
)

// RenderPatchModal formats the strategic merge patch preview modal.
func RenderPatchModal(workloadName, namespace, patchContent string, width int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Rightsizing Patch Preview: %s/%s\n", namespace, workloadName))
	sb.WriteString("--------------------------------------------------------------------------------\n")
	sb.WriteString(patchContent)
	sb.WriteString("\n--------------------------------------------------------------------------------\n")
	sb.WriteString("Keyboard shortcuts: [c] Copy to Clipboard | [w] Save to File | [Esc] Close\n")

	return sb.String()
}
