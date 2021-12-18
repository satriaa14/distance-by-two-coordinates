package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// By : Agung Satria
// ------- Please Solve the Problem -------
// I’m trying to find the closest point from some points that I have, for example I have about 1000 set of geographhical coordinates (lat,long).
// Given one coordinates, I want to find the closest one from that set.
// Note that the list of point changes all the time, and the closes distance depend on when and where the user’s point.
// What is the best optimized solution for this ?
// Please implement this in a language you are comfortable with and push to github.

const (
	// Initial Point
	latitude  float64 = 0
	longitude float64 = 0

	// Earth Radius
	centerRadius float64 = 6371000 // In meter
)

type point struct {
	namePoint string
	lat       float64
	long      float64
	dist      float64 // In meter
}

var tempMem sync.Map

func initPoints() {
	for i := 1; i <= 1000; i++ {
		la, lo := randLatLngFromCenter(latitude, longitude, centerRadius)
		np := fmt.Sprintf("P%04s", fmt.Sprint(i))
		tempMem.Store(np, point{
			namePoint: np,
			lat:       la,
			long:      lo,
		})
	}
}

// Change all list every 2 seconds
func pointsMover(wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		for {
			tempMem.Range(func(k, v interface{}) bool {
				la, lo := randLatLngFromCenter(latitude, longitude, centerRadius)
				tempMem.Store(k, point{
					namePoint: k.(string),
					lat:       la,
					long:      lo,
				})
				return true
			})

			time.Sleep(time.Second * 3)
		}
	}()

	// Runner will run until 30 seconds
	time.Sleep(time.Second * 30)
}

// Get list 5 closest point from client lat-long
func getClosestPoint() {
	nearby := []point{}

	// Random lat-long
	la, lo := randLatLngFromCenter(latitude, longitude, centerRadius)
	tempMem.Range(func(k, v interface{}) bool {
		val := v.(point)
		nearby = append(nearby, point{
			namePoint: val.namePoint,
			lat:       val.lat,
			long:      val.long,
			dist:      distance(val.lat, la, val.long, lo),
		})
		// Un-Comment if you want to see current all lists
		// fmt.Println("range (): ", v)
		return true
	})

	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].dist < nearby[j].dist
	})

	// Your point
	fmt.Println(fmt.Sprintf("\n====================================\nYour Latitude  : %v\nYour Longitude : %v", la, lo))

	// Print candidata 5 nearbies
	for i, v := range nearby[:5] {
		fmt.Println(fmt.Sprintf("-----------------%d------------------\nName Point     : %v\nDistance       : %.4f Kilometers\nLatitude       : %v\nLogitude       : %v", i+1, v.namePoint, v.dist/1000, v.lat, v.long))
	}
	fmt.Println("====================================")
}

// Geographic information systems (GIS) Algorithm
// randLatLngFromCenter (center (for Lat, Long center location), radius (in meter)) returning location Lat and Long
func randLatLngFromCenter(centerLatitude, centerLongitude, radius float64) (float64, float64) {
	y0 := centerLatitude
	x0 := centerLongitude
	rd := radius / 111300

	rand.Seed(time.Now().UnixNano())
	u := rand.Float64()
	v := rand.Float64()

	w := rd * math.Sqrt(u)
	t := 2 * math.Pi * v
	x := w * math.Cos(t)
	y := w * math.Sin(t)

	x1 := x + x0
	y1 := y + y0

	return y1, x1
}

// ----- ----- ----- ----- CORE ----- ----- ----- -----
// Get distance from 2 coordinates by Haversine formula
func distance(lat1, lat2, lon1, lon2 float64) float64 {

	// The math module contains a function
	// named toRadians which converts from
	// degrees to radians.
	lon1 = lon1 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180
	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	var dlon = lon2 - lon1
	var dlat = lat2 - lat1
	var a = math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)

	var c = 2 * math.Asin(math.Sqrt(a))

	var r float64 = 6371000 // Radius of earth in meter.

	// calculate the result
	return c * r
}

// Flush points
func flush() {
	tempMem.Range(func(k, v interface{}) bool {
		tempMem.Delete(k)
		return true
	})
}

func main() {

	// This state will generate 1000 random coordinate (lat, long) points
	// based from center lat-long (0, 0) and on earth radius (in meters)
	initPoints()

	var wg sync.WaitGroup

	wg.Add(1)

	// All list Point Coordinates will update every 3 seconds
	go pointsMover(&wg)

	// Dummy clients will generate random coordinate (lat, long) points
	// Then system get 5 nearbies list points, and based on clients coordinate request
	// after get distance by Haversine formula
	// Client will get nearby every 5 seconds
	go func() {
		for {
			getClosestPoint()
			time.Sleep(time.Second * 5)
		}
	}()

	// Will wait ontil points mover process (30 seconds)
	wg.Wait()

	// flush memory
	flush()

	fmt.Println("Server Shut Down")

}
