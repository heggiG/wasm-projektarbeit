# GO WebAssembly

Schaue die neueste Version auf [github pages](https://heggig.github.io/wasm-projektarbeit/) an.

# Wie funktioniert WebAssembly

WebAssembly ist kompilierter Byte Code der direkt vom Browser ausgeführt werden kann.
Dadurch können C, C++ und ähnliche low-level Sprachen im Web verwendet werden.

Moderne Browser bieten eine API mit der WebAssembly und JavaScript Code interagieren können, dadurch
kann man (hier zum Beispiel in Go) in der ausgewählten Sprache direkt mit dem Dom interargieren,
aufrufbare JavaScript methoden erzeugen und alles andere machen was auch in JavaScript möglich ist.

## Go Code im Browser verwenden

Zunächst muss der Code zu WASM kompiliert werden. Hierfür muss dem Go Compiler das korrekte
Target mitgegeben werden. Dazu setzt man `GOOS=js` und `GOARCH=wasm`.

Um das Go Modul zu wasm zu kompilieren lautet der Befehl also:
```shell
GOOS=js GOARCH=wasm go build -o app.wasm .
```


