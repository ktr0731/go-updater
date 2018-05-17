package brew

import updater "github.com/ktr0731/go-updater"

var _ updater.Means = (*HomebrewClient)(nil)
