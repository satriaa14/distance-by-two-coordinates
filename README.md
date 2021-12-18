# Distance By Two Coordinates

### Problems
- I’m trying to find the closest point from some points that I have, for example I have about 1000 set of geographhical coordinates (lat,long).
- Given one coordinates, I want to find the closest one from that set.
- Note that the list of point changes all the time, and the closes distance depend on when and where the user’s point.
- What is the best optimized solution for this ?
- Please implement this in a language you are comfortable with and push to github.

### Solving
- Will implemen this using Golang.
- Generate 1000 random coordinate points using _Geographic Information Systems (GIS) Algorithm_ and store it in hash map. 
- All points will update it's coordinate (lat, long) every 3 seconds
- Create dummy request, and the dummy request will generate random coordinate by _GIS Algorithm_.
- Get distance by two coordinates (calc every points coordinates, with dummy coordinate request) using Haversine formula. And store the distance in a slice (array) variable.
- Sort the slice that contains distance ascending (by nearbies).
- Print, request coordinate, and 5 nearbies points coordinate.

### Best Optimized ?
- It can be improved using [Redis](https://redis.io/commands/geodist) Geo-Redis to handle many requests and calculation. Command : _GEOADD, GEODIST, GEO*, etc..._

### License
MIT
**Free Software, Hell Yeah!**