DONE - move to Raster before adding tap - may sort out coord transformations

DONE - use magnitude for colour scaling

DONE - add double-buffering for image update

DONE - justify the labels on the entry boxes

- DONE respond to mouse clicks

  - DONE left click to recentre

- DONE finish tap-to-move

  - DONE remove PosX/PosY/PosMag
  - DONE add mandel x, y, width, height (float/complex)
  - DONE add 'onTap' normalised to 0.0 -> 1.0
    - in image space, convert tap event to 0.0->1.0 space
    - in mandel, convert tap event to calculated (x,y) space
    - recentre to that

- BUG: ? some clicks don't register?

- DONE break types out into files/pkgs

- DONE add pt/rect types

- DONE hook mouse wheel scroll for zoom

- fix flicker

- DONE fix aspect ratio (e.g. 1920x1080)

- update fyne controls to be more ergonmic (sliders etc)

- move to using canvas.Raster
  "If you wish to render a pixel-specific image please use canvas.Raster -
  everything else is scaled according to device and user preference."

- move UI values into an opts structure?
  - add tick to the opts structure
  - make mandel reference it

DONE - use xlo + xwidth instead of xlo and xhi (better update behaviour)

- try running on mobile and/or wasm

- provide some control (or presets) for colour animation

- show magnitude at mouse co-ords

DONE - add zoom button + and -
NO - also left/right/up/down?

- set pixel w+h based on window size

- add palette-rotating animation option

- fix entry scaling

DONE - interaction between refresh and high number of steps causes
flickering/blanking. Can we do better?

DONE - use a (fixed) number of goros to parallelise the mandel calc

- optimisation?
  - the work done by the goros in UpdateMagMap isn't equal. The early escape
    for threshold will be hit more on the early and late columns so they do
    less work.
    Perhaps a per-pixel approach would be faster?
    goro-per-pixel likely too much overhead, so channel+workers?
