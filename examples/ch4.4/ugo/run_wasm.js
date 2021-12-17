const fs = require('fs');

var wasmInstance = null

WebAssembly.instantiate(
	new Uint8Array(fs.readFileSync('./a.out.wasm')),
	{
		env: {
			ugo_builtin_println: function(n) {
				console.log(n);
				return 0;
			},
			ugo_builtin_exit: function(n) {
				console.log("exit:", n);
				return 0;
			}
		}
	}
).then(result => {
	wasmInstance = result.instance;
	wasmInstance.exports.main();
}).catch(e => {
	console.log(e);
});
