package main

import "flag"

func main() {
	replicationFactor := flag.Int("r", 2, "replication factor")
	flag.Parse()
}
