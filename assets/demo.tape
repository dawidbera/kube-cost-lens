# VHS Tape script for recording kube-cost-lens demo GIF
# Output: assets/demo.gif

Output assets/demo.gif

Set FontSize 14
Set Width 1200
Set Height 700
Set Padding 20
Set Theme "Dracula"

Type "./bin/kube-cost-lens --demo"
Enter
Sleep 2s

# Navigate down workloads
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms

# Inspect container drilldown modal
Enter
Sleep 2.5s

# Open patch preview
Type "p"
Sleep 2.5s

# Close modal
Escape
Sleep 1s

# Switch to Namespaces view
Tab
Sleep 2s

# Switch to Nodes view
Tab
Sleep 2s

# Filter demonstration
Type "/"
Sleep 500ms
Type "spot"
Sleep 1.5s
Escape
Sleep 1s

# Clean exit
Type "q"
Sleep 500ms
