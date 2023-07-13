// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

{
    // hist is the counters of opcode bigrams
    hist: {},
    // lastOp is last operation
    lastOp: '',
    // execution depth of last op
    lastDepth: 0,
    // step is invoked for every opcode that the VM executes.
    step: function(log, db) {
        var op = log.op.toString();
        var depth = log.getDepth();
        if (depth == this.lastDepth){
            var key = this.lastOp+'-'+op;
            if (this.hist[key]){
                this.hist[key]++;
            }
            else {
                this.hist[key] = 1;
            }
        }
        this.lastOp = op;
        this.lastDepth = depth;
    },
    // fault is invoked when the actual execution of an opcode fails.
    fault: function(log, db) {},
    // result is invoked when all the opcodes have been iterated over and returns
    // the final result of the tracing.
    result: function(ctx) {
        return this.hist;
    },
}
