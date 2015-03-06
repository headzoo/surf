Surf makes it easy to download page assets, such as images, stylesheets, and scripts. They can even be downloaded asynchronously.

```go
bow := surf.NewBrowser()
err := bow.Open("http://www.reddit.com")
if err != nil { panic(err) }

// Download the images on the page and write them to files.
for _, image := range bow.Images() {
    filename := "/home/joe/Pictures" + image.URL.Path()
    fout, err := os.Create(filename)
    if err != nil {
    	log.Printf(
    	    "Error creating file '%s'.", filename)
    	continue
    }
    defer fout.Close()
    
    _, err = image.Download(fout)
    if err != nil {
    	log.Printf(
    	    "Error downloading file '%s'.", filename)
    }
}

// Downloading assets asynchronously takes a little more work, but isn't difficult.
// The DownloadAsync() method takes an io.Writer just like the Download() method,
// plus an instance of AsyncDownloadChannel. The DownloadAsync() method will send
// an instance of browser.AsyncDownloadResult to the channel when the download is
// complete.
ch := make(AsyncDownloadChannel, 1)
queue := 0
for _, image := range bow.Images() {
    filename := "/home/joe/Pictures" + image.URL.Path()
	fout, err := os.Create(filename)
	if err != nil {
		log.Printf(
			"Error creating file '%s'.", filename)
		continue
	}
	
	image.DownloadAsync(fout, ch)
	queue++
}

// Now we wait for each download to complete.
for {
	select {
	case result := <- ch:
	    // result is the instance of browser.AsyncDownloadResult sent by the
	    // DownloadAsync() method. It contains the writer which you need to
	    // close. It also contains the asset itself, and an error instance if
	    // there was an error.
		result.Writer.Close()
		if result.Error != nil {
		    log.Printf("Error download '%s'. %s\n", result.Asset.Url(), result.Error)
		} else {
		    log.Printf("Downloaded '%s'.\n", result.Asset.Url())
		}
		
		queue--
		if queue == 0 {
			goto FINISHED
		}
	}
}
	
FINISHED:
close(ch)
log.Println("Downloads complete!")
```
When downloading assets asynchronously, you should keep in mind the potentially large number of assets embedded into a typical web page. For that reason you should setup a queue that downloads only a few at a time.
