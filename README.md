# Reactive Drop Map Theme Generator

[This website](https://mapthemes.reactivedrop.com/) is intended as a source of inspiration for map makers who can't decide on a theme for their next map. It chooses a random adjective and/or noun from a list and shows you one of several pre-generated images associated with those words.

These [images](https://mapthemes.reactivedrop.com/images/) are generated using [Stable Diffusion](https://en.wikipedia.org/wiki/Stable_Diffusion), an algorithm that allows a computer to hallucinate an image based on its knowledge of an extremely large set of labelled images.

The generation process was run on an NVIDIA T1000 with 8GB of VRAM and took about a minute for the first pass and about five minutes for the upscaling pass for each image.

This process took 3 months for the initial set of 21,810 images (10 variants of each combination of 51 adjectives and 41 nouns, plus 10 variants of each word alone with a two-word overlap between the lists) and started on July 26th 2023.

This page simply displays a random image from the set of thousands, along with the adjective and noun used in the prompt for the image.

Questions? Suggestions? Feedback? [Create an issue on GitHub.](https://github.com/ReactiveDrop/mapthemes.reactivedrop.com/issues)

```
Model used: https://huggingface.co/stabilityai/stable-diffusion-2-1/blob/main/v2-1_768-nonema-pruned.ckpt
Sampling method: Euler a; Sampling steps: 50 (+50 upscaling)
Size: 640x360; Batch count: 10; CFG scale: 7
Hires. fix (Latent): Resize (2x) to 1280x720; Denoising strength 0.7

Seeds: 563560-563569
Prompts: adjective noun, level design render, wide view, dim volumetric lighting, retrofuturism
Negative Prompt: text, 2d, screenshot, watermark
```

For single-word prompts, the adjective slot is used for the word and the noun slot is left blank. [Here is what the prompt generates with no adjective and no noun before upscaling.](www/images/nowords.avif)
