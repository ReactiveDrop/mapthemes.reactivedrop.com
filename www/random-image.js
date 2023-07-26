"use strict";

async function fetchLines(name) {
	const response = await fetch(name);
	const body = await response.text();

	const lines = body.split(/\r?\n/g);
	if (!lines[lines.length]) {
		lines.pop();
	}

	return lines;
}

(async function() {
	const variants = 10;
	const [adjectives, nouns] = await Promise.all([
		fetchLines("adjective.txt"),
		fetchLines("noun.txt"),
	]);

	const img = document.createElement("img");
	const h1 = document.createElement("h1");
	const button = document.createElement("button");

	function showRandomImage() {
		const adjective = adjectives[Math.floor(Math.random() * adjectives.length)];
		const noun = nouns[Math.floor(Math.random() * nouns.length)];
		const variant = Math.floor(Math.random() * variants);

		h1.textContent = adjective + " " + noun + " " + (variant + 1) + "000";
		img.src = "images/" + adjective + "-" + noun + "-" + variant.toString().padStart(4, "0") + ".avif";
	}

	document.body.insertBefore(button, document.body.firstChild);
	document.body.insertBefore(h1, document.body.firstChild);
	document.body.insertBefore(img, document.body.firstChild);
	img.alt = "computer-generated concept art image";
	button.textContent = "Generate Random Theme";
	button.onclick = function(e) {
		e.preventDefault();

		showRandomImage();
	};

	showRandomImage();
})();
