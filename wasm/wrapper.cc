#include <cstddef>
#include <cstdint>

#include "marisa/trie.h"

// note: even on 64-bit platforms with 64-bit size_t, key IDs and lengths are
// limited to uint32s in MARISA

// note: most throws only happen due to programmer error or out-of-bound values,
// other than runtime_error, which gets thrown when a trie is corrupt, and a few
// invalid_argument ones when it's parsing data (and luckily, the throws that we
// can't prevent are just during reads of the pre-allocated trie data, so we
// won't really have a problem with memory leaks even though our fake_throw
// won't call destructors on stack variables)

static marisa::Trie trie;
static marisa::Keyset build;

extern "C" void marisa_new(void *ptr, size_t size) {
    trie.map(ptr, size);
}

extern "C" void marisa_load() {
    trie.read(0);
}

extern "C" void marisa_save() {
    trie.write(0);
}

extern "C" void marisa_build_push(const char *ptr, size_t length, float weight) {
    build.push_back(ptr, length, weight); // copies ptr[:length] into a data block in keyset
    //assert(build[build.size()-1].ptr() != ptr);
}

extern "C" void marisa_build(int flags) {
    // warning: this will continue to reference pointers to data blocks in keyset
    trie.build(build, flags);
}

struct marisa_stat {
    uint32_t size;
    uint32_t io_size;
    uint32_t total_size;
    uint32_t num_tries;
    uint32_t num_nodes;
    uint32_t tail_mode;
    uint32_t node_order;
};

extern "C" struct marisa_stat marisa_stat() {
    return (struct marisa_stat){
        .size = trie.size(),
        .io_size = trie.io_size(),
        .total_size = trie.total_size(),
        .num_tries = trie.num_tries(),
        .num_nodes = trie.num_nodes(),
        .tail_mode = trie.tail_mode(),
        .node_order = trie.node_order(),
    };
}

extern "C" marisa::Agent *marisa_query_new() {
    auto agent = new marisa::Agent;
    agent->init_state(); // heap allocates
    return agent;
}

extern "C" void marisa_query_set_str(marisa::Agent *agent, const char *ptr, size_t len) {
    agent->set_query(ptr, len); // does not copy
}

extern "C" void marisa_query_set_id(marisa::Agent *agent, uint32_t id) {
    agent->set_query(static_cast<size_t>(id));
}

extern "C" void marisa_query_clear(marisa::Agent *agent) {
    agent->clear();
}

extern "C" void marisa_query_free(marisa::Agent *agent) {
    delete agent;
}

// note: we have it as a agent function even for the lookups which only
// return a single node since the agent contains heap-allocated memory which
// wouldn't get cleaned up if it threw and we had it on the stack

extern "C" bool marisa_query_lookup(marisa::Agent *agent) {
    return trie.lookup(*agent);
}

extern "C" bool marisa_query_reverse_lookup(marisa::Agent *agent) {
    if (agent->query().id() >= trie.num_keys()) return false;
    // note: this will always throw if id >= trie.num_keys()
    trie.reverse_lookup(*agent);
    return true;
}

extern "C" bool marisa_query_common_prefix_search(marisa::Agent *agent) {
    return trie.common_prefix_search(*agent);
}

extern "C" bool marisa_query_predictive_search(marisa::Agent *agent) {
    return trie.predictive_search(*agent);
}

extern "C" struct marisa_query_result {
    uint32_t id;
    const char *ptr;
    uint32_t len;
};

extern "C" struct marisa_query_result marisa_query_result(marisa::Agent *agent) {
    auto key = agent->key();
    return (struct marisa_query_result){
        .id = static_cast<uint32_t>(key.id()),
        .ptr = key.ptr(),
        .len = static_cast<uint32_t>(key.length()),
    };
}
