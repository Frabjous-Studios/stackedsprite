# stackedsprite

A simple [Ebitengine](https://ebitengine.org) prototype demonstrating rendering using stacked sprites.

## What is this?

Stacked Sprites are 3D rendered sprites from horizontal 2D slices like the ones in the sprite strip below.

![Stacked Sprite of a Blue Car](./cmd/test/img/BlueCar.png)

It's simple to render stacked sprites to the screen to create a fake 3D effect using rotations at any angle. 
Each slice of the sprite is rotated separately and rendered bottom to top. An illusion of orthographic perspective is
created with the y-coordinate shifting up by one for every layer. 
While this is inefficient when many sprites need to be rotated every frame, the technique can work for some games where
a large number of rotations isn't required per frame.

The implementation here caches frames after render and tracks the origin of the sprite at the center. It's a very 
early prototype, no guarantee is given about the stability of the API.

But it works!

https://github.com/Frabjous-Studios/stackedsprite/assets/22620342/41f01c3c-dfaa-4cac-b33e-2e758f5922f2
