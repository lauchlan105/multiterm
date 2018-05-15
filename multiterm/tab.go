package multiterm

//Tab asdf
type Tab struct {
	manager *Terminal
	id      string
	name    string
	active  bool
	buffer  []string
}

//Terminate kills the current tab
//via the Terminal objects removeTab(id) func
func (t *Tab) Terminate() {
	t.manager.removeTab(t.id)
}

//Open sets the tab's active state to true
// and adds it to the managers activeTabs slice
func (t *Tab) Open() {
	if t.active {
		return
	}
	t.active = true

	t.manager.activeTabs = append(t.manager.activeTabs, t)
}

//Close sets the tab's active state to false
// and removes it from the managers activeTabs slice
func (t *Tab) Close() {
	if !t.active {
		return
	}
	t.active = false

	//Iterate tabs and break when tab is found
	for i, tab := range t.manager.activeTabs {
		if tab.id == t.id {
			t.manager.activeTabs = append(t.manager.activeTabs[:i],
				t.manager.activeTabs[i+1:]...)
			break
		}
	}

}
