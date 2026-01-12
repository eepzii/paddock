package browser

func New(config Config) (Browser, error) {

	var b = Browser{}

	var browserPath = config.CustomBrowserPath
	if config.CustomBrowserPath == "" {
		browser, err := findBrowser()
		if err != nil {
			return b, err
		}
		browserPath = browser
	}

	b.commandPath = browserPath
	b.profilePath = config.FileManager.BrowserProfileDirPath()
	b.headless = config.Headless
	b.brand = findBrowserBrand(browserPath)

	return b, nil
}
