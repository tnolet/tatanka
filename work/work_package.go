package work

func (w *WorkPackage) Shift() WorkItem {
	var item WorkItem
	if len(w.WorkItems) > 0 {
		item, w.WorkItems = w.WorkItems[0], w.WorkItems[1:]
	}
	return item
}
