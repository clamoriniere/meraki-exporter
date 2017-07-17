package main

import (
	"log"

	"github.com/cedriclam/meraki-exporter/pkg/meraki"
)

func main() {
	config := meraki.NewConfig()
	config.Init()

	exporter := meraki.NewExporter(config)
	log.Fatal(exporter.ListenAndServe())
}
