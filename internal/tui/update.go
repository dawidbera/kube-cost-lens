// Package tui contains state machine update handlers following The Elm Architecture.
package tui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/patcher"
	"github.com/dawidbera/kube-cost-lens/internal/tui/views"
)

// Update processes incoming tea.Msg events and advances the state machine.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		tableHeight := msg.Height - 8
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.Table.SetHeight(tableHeight)
		m.Table.SetWidth(msg.Width - 4)

		m.Viewport.Width = msg.Width - 10
		m.Viewport.Height = msg.Height - 10

		m.updateTableData()

	case TickMsg:
		cmds = append(cmds, TickCmd(m.RefreshInterval))
		if !m.IsFiltering {
			cmds = append(cmds, FetchClusterDataCmd(m.K8sProvider, m.SelectedNS, m.Calculator))
		}

	case ClusterDataUpdatedMsg:
		m.IsLoading = false
		if msg.Err != nil {
			m.ErrorMessage = msg.Err.Error()
		} else {
			m.ErrorMessage = ""
			m.Workloads = msg.Workloads
			m.Namespaces = msg.Namespaces
			m.Nodes = msg.Nodes
			m.updateTableData()
		}

	case ClipboardCopiedMsg:
		if msg.Success {
			m.StatusMessage = "✓ Patch copied to clipboard!"
		} else {
			m.StatusMessage = "✗ Failed to copy to clipboard"
		}

	case PatchSavedMsg:
		if msg.Err == nil {
			m.StatusMessage = fmt.Sprintf("✓ Saved patch to %s", msg.Path)
		} else {
			m.StatusMessage = fmt.Sprintf("✗ Failed to save patch: %v", msg.Err)
		}

	case tea.KeyMsg:
		// 1. Live Filter Input Handling
		if m.IsFiltering {
			switch msg.Type {
			case tea.KeyEnter:
				m.IsFiltering = false
				m.FilterInput.Blur()
				return m, nil
			case tea.KeyEsc:
				m.IsFiltering = false
				m.FilterInput.Blur()
				m.FilterInput.SetValue("")
				m.FilterQuery = ""
				m.updateTableData()
				return m, nil
			default:
				var tiCmd tea.Cmd
				m.FilterInput, tiCmd = m.FilterInput.Update(msg)
				m.FilterQuery = m.FilterInput.Value()
				m.updateTableData()
				return m, tiCmd
			}
		}

		// 2. Active Modal Dialog Handling
		if m.ActiveModal != ModalNone {
			switch {
			case key.Matches(msg, m.KeyMap.Escape), msg.String() == "q":
				m.ActiveModal = ModalNone
				m.StatusMessage = ""
				return m, nil

			case key.Matches(msg, m.KeyMap.Patch) && m.ActiveModal == ModalContainerDetails:
				if m.SelectedWorkload != nil {
					return m.openPatchModal()
				}

			case key.Matches(msg, m.KeyMap.Copy) && m.ActiveModal == ModalPatchPreview:
				if m.GeneratedPatch != "" {
					err := clipboard.WriteAll(m.GeneratedPatch)
					m.StatusMessage = "✓ Patch copied to clipboard!"
					if err != nil {
						m.StatusMessage = fmt.Sprintf("Clipboard note: %v", err)
					}
				}
				return m, nil

			case key.Matches(msg, m.KeyMap.Write) && m.ActiveModal == ModalPatchPreview:
				if m.SelectedWorkload != nil && m.GeneratedPatch != "" {
					fileName := fmt.Sprintf("patch-%s.yaml", m.SelectedWorkload.Name)
					err := patcher.WritePatchToFile(fileName, m.GeneratedPatch)
					if err == nil {
						m.StatusMessage = fmt.Sprintf("✓ Written %s", fileName)
					} else {
						m.StatusMessage = fmt.Sprintf("✗ Error: %v", err)
					}
				}
				return m, nil

			case key.Matches(msg, m.KeyMap.Help) && m.ActiveModal == ModalHelp:
				m.ActiveModal = ModalNone
				return m, nil

			default:
				var vpCmd tea.Cmd
				m.Viewport, vpCmd = m.Viewport.Update(msg)
				return m, vpCmd
			}
		}

		// 3. Normal Table View Handling
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.KeyMap.Help):
			m.ActiveModal = ModalHelp
			return m, nil

		case key.Matches(msg, m.KeyMap.Refresh):
			m.IsLoading = true
			m.StatusMessage = "Refreshing cluster data..."
			cmds = append(cmds, FetchClusterDataCmd(m.K8sProvider, m.SelectedNS, m.Calculator))

		case key.Matches(msg, m.KeyMap.Filter):
			m.IsFiltering = true
			m.FilterInput.Focus()
			return m, textinput.Blink

		case key.Matches(msg, m.KeyMap.NextTab):
			m.ActiveView = (m.ActiveView + 1) % 3
			m.updateTableData()

		case key.Matches(msg, m.KeyMap.PrevTab):
			m.ActiveView = (m.ActiveView + 2) % 3
			m.updateTableData()

		case msg.String() == "1":
			m.ActiveView = ViewWorkloads
			m.updateTableData()

		case msg.String() == "2":
			m.ActiveView = ViewNamespaces
			m.updateTableData()

		case msg.String() == "3":
			m.ActiveView = ViewNodes
			m.updateTableData()

		case key.Matches(msg, m.KeyMap.Enter):
			if m.ActiveView == ViewWorkloads {
				m.resolveSelectedWorkload()
				if m.SelectedWorkload != nil {
					content := views.RenderContainerDetailsModal(m.SelectedWorkload, m.Rightsizer, m.Width)
					m.Viewport.SetContent(content)
					m.Viewport.GotoTop()
					m.ActiveModal = ModalContainerDetails
					return m, nil
				}
			}

		case key.Matches(msg, m.KeyMap.Patch):
			if m.ActiveView == ViewWorkloads {
				m.resolveSelectedWorkload()
				if m.SelectedWorkload != nil {
					return m.openPatchModal()
				}
			}

		default:
			var tableCmd tea.Cmd
			m.Table, tableCmd = m.Table.Update(msg)
			m.resolveSelectedWorkload()
			return m, tableCmd
		}
	}

	return m, tea.Batch(cmds...)
}

// openPatchModal calculates recommendations and generates the YAML patch preview.
func (m Model) openPatchModal() (tea.Model, tea.Cmd) {
	if m.SelectedWorkload == nil {
		return m, nil
	}

	var cRecs []finops.ContainerRecommendation
	for _, c := range m.SelectedWorkload.Containers {
		rec := m.Rightsizer.RecommendContainer(c.Name, c.CPURequest, c.CPULimit, c.CPUUsage, c.MemoryRequest, c.MemoryLimit, c.MemoryUsage)
		cRecs = append(cRecs, rec)
	}

	wRec := m.Rightsizer.RecommendWorkload(m.SelectedWorkload.Namespace, m.SelectedWorkload.Name, m.SelectedWorkload.Kind, m.SelectedWorkload.Replicas, cRecs)
	patchText, err := patcher.GenerateStrategicMergePatch(wRec)
	if err != nil {
		m.StatusMessage = fmt.Sprintf("Error generating patch: %v", err)
		return m, nil
	}

	m.GeneratedPatch = patchText
	content := views.RenderPatchModal(m.SelectedWorkload.Name, m.SelectedWorkload.Namespace, patchText, m.Width)
	m.Viewport.SetContent(content)
	m.Viewport.GotoTop()
	m.ActiveModal = ModalPatchPreview
	return m, nil
}

// resolveSelectedWorkload identifies the active workload pointer based on table selection.
func (m *Model) resolveSelectedWorkload() {
	if m.ActiveView != ViewWorkloads || len(m.FilteredWorkloads) == 0 {
		return
	}
	idx := m.Table.Cursor()
	if idx >= 0 && idx < len(m.FilteredWorkloads) {
		m.SelectedWorkload = &m.FilteredWorkloads[idx]
	}
}

// updateTableData rebuilds columns and rows for the active tab view safely.
func (m *Model) updateTableData() {
	// Clear rows first to prevent column count mismatch during table re-render
	m.Table.SetRows([]table.Row{})

	switch m.ActiveView {
	case ViewWorkloads:
		m.Table.SetColumns(views.BuildWorkloadsColumns(m.Width))
		rows, filtered := views.BuildWorkloadsRows(m.Workloads, m.FilterQuery)
		m.FilteredWorkloads = filtered
		m.Table.SetRows(rows)
		m.resolveSelectedWorkload()

	case ViewNamespaces:
		m.Table.SetColumns(views.BuildNamespacesColumns(m.Width))
		rows := views.BuildNamespacesRows(m.Namespaces, m.FilterQuery)
		m.Table.SetRows(rows)

	case ViewNodes:
		m.Table.SetColumns(views.BuildNodesColumns(m.Width))
		rows := views.BuildNodesRows(m.Nodes, m.FilterQuery)
		m.Table.SetRows(rows)
	}
}
