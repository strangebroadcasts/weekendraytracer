# Development notes

* Use doubles/float64 when possible
* RTiaW views rays as parametric lines $p(t) = A + tB$ where A is the ray origin and B is the ray direction.
* Core idea: "calculate which ray goes from the eye to a pixel, compute what that ray intersects, and compute a color for that intersection point".
* Conventions:
  * Camera center at $[0, 0, 0]$
  * Right-handed coordinate system: X-axis goes right, Y-axis goes up, Z-axis goes towards the viewer (and conversely, negative Z-axis is "into the screen")
  * Traverse the screen from the lower-left corner
* No need to normalize ray direction vector
