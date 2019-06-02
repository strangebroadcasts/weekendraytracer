# Development notes

* Use doubles/float64 when possible
* RTiaW views rays as parametric lines $p(t) = A + tB$ where A is the ray origin and B is the ray direction.
* Core idea: "calculate which ray goes from the eye to a pixel, compute what that ray intersects, and compute a color for that intersection point".
* Conventions:
  * Camera center at $[0, 0, 0]$
  * Right-handed coordinate system: X-axis goes right, Y-axis goes up, Z-axis goes towards the viewer (and conversely, negative Z-axis is "into the screen")
  * Traverse the screen from the lower-left corner
* No need to normalize ray direction vector
* Recall the equation for a sphere of radius R centered at the origin is $x^2 + y^2 + z^2 = R^2$
* The ray from the center $C=[cx,cy,cz]$ to the point $p=[x,y,z]$ is $p-C$.
* We can then rewrite the sphere equation for rays as $(p-C) \cdot (p-C) = R^2$
* First, we want to determine if there's an intersection at all. We can solve this analytically:
  * $(p(t) - c) \cdot (p(t) - c) = R^2$
  * $(A + tB - C) \cdot (A + tB - C) - R^2 = 0$
  * Let $D = A - C$, so $(D + tB) \cdot (D + tB) - R^2 = 0$
  * TODO: Finish this derivation
* RTiaW's solution is $t^2 (B \cdot B) + 2t (B \cdot (A-C)) + (A-C) \cdot (A-C) - R^2 = 0$
* The surface normal is perpendicular to the surface. For spheres, the normal is in the direction of the intersection point minus the sphere center.
* We normalize normals to unit length (convenient for shading)
* To make the raytracer generalizable, we create a general interface "Hittable", which determines whether a ray intersects a particular primitive.
* This ray tracer doesn't use "stratification" - look up at a later date
* "Diffuse objects that don't emit light merely take on the color of their surroundings, but modulate it with their own intrinsic color."
* Light is reflected off diffuse surfaces in random directions.
* Light can also be absorbed by diffuse surfaces - the darker the surface, the more light it absorbs
* We use an approximation that should be "exactly correct for ideal diffuse surfaces":
  * We pick a random point $s$ from the unit radius sphere tangent to the intersection point $p$, and cast a ray from $p$ to $s$.
  * This sphere has center $(p+N)$
  * This approach requires generating random points inside a unit sphere. RTiaW uses a rejection method here.
* Look up color spaces, gamma transform: RTiaW uses a very simple transform here.
