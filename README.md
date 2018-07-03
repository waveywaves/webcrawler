# webcrawler
Command Line utility for crawling a Website. 

This is a simple utility for a user to be able to generate a pseudo-sitemap to the internal and not the external links to the website. 
The main purpose of this tool was to make web scraping a site for the internal links possible on the command line.

This project is new and we are still building on added functionality.
Please look at the [CONTRIBUTING](https://github.com/waveywaves/webcrawler/blob/master/docs/CONTRIBUTING.md) section.

### Installing from source

Follow the commands below, step-by-step to install `scrape` on your system.

```sh
$ git clone https://github.com/waveywaves/webcrawler.git
$ cd webcrawler
$ # set your GOPATH correctly
$ make build

$ ./webcrawler webscraper.io #example
```
### Building for docker 

```sh
$ make docker-build
```