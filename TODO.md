DONE - use magnitude for colour scaling

DONE - add double-buffering for image update

DONE - justify the labels on the entry boxes

- move UI values into an opts structure?
    - add tick to the opts structure
    - make mandel reference it

- use xlo + xwidth instead of xlo and xhi (better update behaviour)

- try running on mobile and/or wasm

- provide some control (or presets) for colour animation

- show magnitude at mouse co-ords

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
      
