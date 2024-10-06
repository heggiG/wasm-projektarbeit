# GO WebAssembly

[Die neueste Version auf github pages anaschauen.](https://heggig.github.io/wasm-projektarbeit/)

# Wie funktioniert WebAssembly

WebAssembly ist kompilierter Byte Code der direkt vom Browser ausgeführt werden kann.
Dadurch können C, C++ und ähnliche low-level Sprachen* im Web verwendet werden.

Moderne Browser bieten eine API mit der WebAssembly und JavaScript Code interagieren können, dadurch
kann man (hier zum Beispiel in Go) in der ausgewählten Sprache direkt mit dem Dom interargieren,
aufrufbare JavaScript methoden erzeugen und alles andere machen was auch in JavaScript möglich ist.

\* Programmiersprachen mit einem garbage-collected Speichermodell sind ein [Langzeitziel](https://webassembly.org/docs/high-level-goals/) des WASM Projekts

## Go Code im Browser verwenden

Zunächst muss der Code zu WASM kompiliert werden. Hierfür muss dem Go Compiler das korrekte
Target mitgegeben werden. Dazu setzt man `GOOS=js` und `GOARCH=wasm`.

Um das Go Modul zu wasm zu kompilieren lautet der Befehl also:
```shell
GOOS=js GOARCH=wasm go build -o app.wasm .
```

Anschliessend muss die Datei `wasm_exec.js` in per `<script>` tag eingebunden werden.
Im Fall von Go kann die Datei entweder aus dem `GOHOME` Verzeichnis kopiert werden 
(`cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm_exec.js`) oder als Quelle verlinkt 
werden (https://github.com/golang/go/blob/master/lib/wasm/wasm_exec.js). Dieses 
Skript erzeugt eine globale Go-Klasse, die die Interaktion mit der WebAssembly API wrapped.

