package multiterm

//Tab asdf
type Tab struct {
	manager *Terminal
	id      string
	name    string
	active  bool
	buffer  [][]string
}

//Terminate kills the current tab
//via the Terminal objects removeTab(id) func
func (t *Tab) Terminate() {
	t.manager.removeTab(t.id)
}

//Open asdf
func (t *Tab) Open() {
	if t.active {
		return
	}
	t.active = true

	t.manager.activeTabs = append(t.manager.activeTabs, t)
}

//Close asdf
func (t *Tab) Close() {
	if !t.active {
		return
	}
	t.active = false

}
