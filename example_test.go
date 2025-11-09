package marisa_test

import (
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"

	"github.com/pgaskin/go-marisa"
)

func EnglishWords() iter.Seq[string] {
	return slices.Values(words)
}

func ExampleTrie_Build() {
	keys := []string{
		"a",
		"a/b",
		"a/b/c",
		"b",
		"b/c",
		"c",
	}

	var trie marisa.Trie
	if err := trie.Build(slices.Values(keys), marisa.Config{}); err != nil {
		panic(err)
	}

	fmt.Println(&trie)

	// Output:
	// marisa.Trie(size=6 io_size=4088 total_size=3403 num_tries=3 num_nodes=7 tail_mode=text node_order=weight)
}

func ExampleTrie_BuildWeights() {
	keys := map[string]float32{
		"a":     1.0,
		"a/b":   1.0,
		"a/b/c": 1.0,
		"b":     1.0,
		"b/c":   5.0,
		"c":     3.5,
	}
	// note: weights are cumulative and include the weights of all children

	cfg := marisa.Config{
		NodeOrder: marisa.WeightOrder,
	}

	var trie marisa.Trie
	if err := trie.BuildWeights(maps.All(keys), cfg); err != nil {
		panic(err)
	}

	fmt.Println(&trie)

	var err error
	for _, key := range trie.DumpSeq()(&err) {
		fmt.Println(key)
	}
	if err != nil {
		panic(err)
	}

	// Output:
	// marisa.Trie(size=6 io_size=4088 total_size=3403 num_tries=3 num_nodes=7 tail_mode=text node_order=weight)
	// b
	// b/c
	// c
	// a
	// a/b
	// a/b/c
}

func Example_save() {
	var trie marisa.Trie
	if err := trie.Build(EnglishWords(), marisa.Config{}); err != nil {
		panic(err)
	}

	f, err := os.Create("words.dat")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := trie.WriteTo(f); err != nil {
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}
}

func Example_load() {
	trie, err := marisa.Open("words.dat")
	if err != nil {
		panic(err)
	}
	fmt.Println(trie)
}

func Example_marshal() {
	var trie marisa.Trie
	if err := trie.Build(EnglishWords(), marisa.Config{}); err != nil {
		panic(err)
	}

	buf, err := trie.MarshalBinary()
	if err != nil {
		panic(err)
	}

	fmt.Println(trie.Size(), trie.DiskSize(), len(buf))

	if err := trie.UnmarshalBinary(buf); err != nil {
		panic(err)
	}

	fmt.Println(trie.Size())

	// Output:
	// 466550 1413352 1413352
	// 466550
}

func ExampleTrie_query() {
	var trie marisa.Trie
	if err := trie.Build(EnglishWords(), marisa.Config{}); err != nil {
		panic(err)
	}

	id, ok, err := trie.Lookup("iterate")
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("not found")
	}

	key, ok, err := trie.ReverseLookup(id)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("not found")
	}

	fmt.Println("l", id, key)

	for id, key := range trie.PredictiveSearchSeq("iterat")(&err) {
		fmt.Println("p", id, key)
	}
	if err != nil {
		panic(err)
	}

	for id, key := range trie.CommonPrefixSearchSeq("iterated")(&err) {
		fmt.Println("c", id, key)
	}
	if err != nil {
		panic(err)
	}

	// Output:
	// l 262491 iterate
	// p 352923 iterative
	// p 413344 iteratively
	// p 413345 iterativeness
	// p 352924 iteration
	// p 413346 iterations
	// p 352925 iterating
	// p 262491 iterate
	// p 352926 iterated
	// p 352927 iterately
	// p 352928 iterates
	// p 262492 iterator
	// p 352929 iterator's
	// p 352930 iterators
	// c 4 i
	// c 2235 ite
	// c 17192 iter
	// c 262491 iterate
	// c 352926 iterated
}
