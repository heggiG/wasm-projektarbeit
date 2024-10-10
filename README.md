# GO WebAssembly

[Die neueste Version auf github pages anschauen.](https://heggig.github.io/wasm-projektarbeit/)

# Wie funktioniert WebAssembly

WebAssembly ist kompilierter Byte Code der direkt vom Browser ausgeführt werden kann.
Dadurch können C, C++ und ähnliche low-level Sprachen* im Web verwendet werden.

Moderne Browser bieten eine API mit der WebAssembly und JavaScript Code interagieren können, dadurch
kann man (hier zum Beispiel in Go) in der ausgewählten Sprache direkt mit dem Dom interagieren,
aufrufbare JavaScript methoden erzeugen und alles andere machen, was auch in JavaScript möglich ist.

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

# Was wurde in diesem Repo implementiert?

In diesem Repo sind zwei verschiedene Themen in WebAssembly implementiert, einmal kann man
ein Bild (png) laden und dann einen von 4 verschiedenen "Filtern" anwenden (Sobel Filter, 
Gausssches Weichzeichnen, Farbshift und Vignetten) und als Zweites befindet sich unten auf der
Seite ein in Go implementiertes Pong Spiel.

## Filter

### Sobel

Der Sobel-Filter ist ein grundlegendes Verfahren zur Kantendetektion in der Bildverarbeitung, 
das auf der Berechnung von Gradienten basiert. Er verwendet zwei diskrete Differenzierungsoperatoren,
die auf die Intensitätswerte eines Bildes angewendet werden, um Änderungen in der Helligkeit 
entlang der horizontalen (x) und vertikalen (y) Achse zu erfassen. Diese Operatoren, auch als 
Sobel-Kernels bekannt, berechnen die erste Ableitung des Bildes in den jeweiligen Richtungen.

Das Ergebnis sind zwei Gradientenbilder, die die Kanteninformationen entlang der x- und y-Achse
darstellen. Durch die Kombination dieser Gradienten kann die Gesamtstärke und Ausrichtung der 
Kanten im Bild ermittelt werden. 

### Gausssches Weichzeichnen (Gaussian Blur)

Das gaußsche Weichzeichnen ist ein Filter um Bilder zu glätten und Rauschen zu reduzieren. 
Es basiert auf der Faltung eines Bildes mit einer gaußschen Funktion, die eine gewichtete 
Mittelung der Pixelwerte in einer bestimmten Umgebung durchführt. Diese Funktion erzeugt 
eine glockenförmige Kurve, bei der nahegelegene Pixel stärker gewichtet werden als 
weiter entfernte.

Die Gewichte des gaussschen Kernels sollten zu 1 aufsummieren da ansonsten die Helligkeit
des Bildes verringert oder erhöht wird, je nachdem ob die Summe der Gewichte größer oder kleiner
1 ist.

Durch diesen Filter werden feine Details im Bild abgeschwächt, was zu einer gleichmäßigeren 
Verteilung der Pixelwerte führt.

### Farbshift

Der Farbshift "shiftet" alle Pixel um einen angegebenen Faktor (hier mithilfe eines Sliders)
in Richtung einer gegebenen Farbe (hier per HTML Color Picker). Wird der Faktor auf 1 gesetzt
(100 %) besteht das Zielbild ausschließlich aus der ausgewählten Farbe.

### Vignette

Die Vignette ist ein Filter, bei dem die Ränder eines Bildes dunkler erscheinen als der 
zentrale Bereich. Mathematisch wird die Vignettierung durch die Anwendung einer Funktion realisiert, 
die die Helligkeit schrittweise von der Mitte des Bildes zu den Rändern hin reduziert.

Die Felder für die X und Y können hier auch durch Klicken auf das Bild ausgefüllt werden, der 
Radius bestimmt den Radius des Bereichs im Bild, der frei bleibt, die Distanz ist ein Faktor der 
die Distanz von Start und Ende der Verdunkelung erhöht bzw. verringert.

## Pong

Pong ist ein klassisches Videospiel, das als eines der ersten kommerziellen Spiele gilt. Es
simuliert eine vereinfachte Version von Tischtennis (Ping-Pong), bei dem zwei Spieler
gegeneinander antreten. Jeder Spieler steuert einen vertikal beweglichen Schläger, der am
linken oder rechten Rand des Bildschirms positioniert ist. Ziel des Spiels ist es, den Ball so
zu schlagen, dass der Gegner ihn nicht zurückspielen kann.

Der Ball bewegt sich in einer geraden Linie über den Bildschirm und ändert seine Richtung, sobald 
er einen Schläger oder die Ober- und Unterkante des Spielfelds berührt. Wenn ein Spieler den
Ball verfehlt, erzielt der andere einen Punkt.

Dieses Spiel ist eigentlich eine Vollbildanwendung und wird dann als iframe in html eingebettet, dies
vereinfacht einiges was einiges in bezug auf Größenberechnung des Spiels vereinfacht.

## WASM Interaktion

Diese zwei Teile zeigen zwei unterschiedliche Wege, um WASM zu benutzen. Einmal zeigt das Bild Filtern
wie man Go Methoden in JavaScript aufruft und Rückgabewerte behandelt und verarbeitet, zum anderen wie man
Vollbild Anwendungen per WASM im Web verwenden kann.
