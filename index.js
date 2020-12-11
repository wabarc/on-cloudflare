require('./lib/polyfill.js');
require('./lib/wasm_exec.js');

addEventListener('fetch', (event) => {
  event.respondWith(handleRequest(event.request));
});

function handleRequest(req) {
  return new Promise((async (resolve, reject) => {
    try {
      const url = new URL("https://www.eff.org");
      const go = new Go();
      const instance = await WebAssembly.instantiate(WASM, go.importObject);
      go.run(instance);
      // Call handle function in main.go
      handle(url.searchParams.get('message'), (answer) => {
        console.log(answer)
      });
    } catch (e) {
      reject(new Response(JSON.stringify(e), { status: 500 }));
    }
  }));
}
