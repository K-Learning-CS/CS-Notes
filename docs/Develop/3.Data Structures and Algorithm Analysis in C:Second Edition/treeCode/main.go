package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	conf, err := parse()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	if conf.help {
		flag.Usage()

		return
	}

	go func() { log.Println(http.ListenAndServe(":6060", nil)) }()

	log.Println("==========>>>>> Input Arguments")
	log.Printf("kind: %s, size: %v, rand: %v, seed: %v", conf.kind, conf.size, conf.rand, conf.seed)

	nodes, tree := newNodes(conf.size), newTree(conf.kind)

	rand.Seed(conf.seed)

	runInsert(conf, nodes, tree)
	runDelete(conf, nodes, tree)
}

func runInsert(conf *config, nodes []*Node[int], tree Tree[int]) {
	shuffle(conf, nodes)

	pause()

	log.Println("==========>>>>> Insert Begin")

	start := time.Now()

	Reset()

	for _, node := range nodes {
		tree.Insert(node)
	}

	log.Println("==========>>>>> Insert Results")
	log.Printf("insert elapse: %vms", time.Since(start).Milliseconds())
	log.Printf("insert search: %.2f", float64(GetSearchCounter())/float64(conf.size))
	log.Printf("insert  fixup: %.2f", float64(GetFixupCounter())/float64(conf.size))
	log.Printf("insert  extra: %.2f", float64(GetExtraCounter())/float64(conf.size))
	log.Printf("insert rotate: %.2f", float64(GetRotateCounter())/float64(conf.size))

	pause()
}

func runDelete(conf *config, nodes []*Node[int], tree Tree[int]) {
	shuffle(conf, nodes)

	pause()

	log.Println("==========>>>>> Delete Begin")

	start := time.Now()

	Reset()

	for _, node := range nodes {
		_ = tree.Find(node.Data())
		tree.Delete(node)
	}

	log.Println("==========>>>>> Delete Results")
	log.Printf("delete elapse: %vms", time.Since(start).Milliseconds())
	log.Printf("delete search: %.2f", float64(GetSearchCounter())/float64(conf.size))
	log.Printf("delete  fixup: %.2f", float64(GetFixupCounter())/float64(conf.size))
	log.Printf("delete  extra: %.2f", float64(GetExtraCounter())/float64(conf.size))
	log.Printf("delete rotate: %.2f", float64(GetRotateCounter())/float64(conf.size))

	pause()
}

func shuffle(conf *config, nodes []*Node[int]) {
	rand.Shuffle(len(nodes), func(i, j int) {
		k := math.Round(float64(j) * conf.rand)
		j = i - int(k)
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
}

func newTree(kind string) Tree[int] {
	switch kind {
	case "rb":
		return NewRBTree[int]()
	default:
		return NewAVLTree[int]()
	}
}

func newNodes(size int) []*Node[int] {
	nodes := make([]*Node[int], size)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}

	return nodes
}

func pause() {
	const duration = 1

	log.Printf("pause for %v seconds...", duration)
	time.Sleep(time.Duration(duration) * time.Second)
}
