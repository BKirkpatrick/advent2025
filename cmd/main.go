package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	s "strings"
)

type Vec3 struct {
	X, Y, Z float64
}

// clusterID == 0 means that those vectors are not linked
type DistanceLog struct {
	id1       int
	id2       int
	v1        Vec3
	v2        Vec3
	d         float64
	clusterId int
}

func main() {
	filepath := "../testdata/day8.txt"

	topN := 3

	coords, _ := loadVec3File(filepath)
	n := len(coords)

	nDistances := n * (n - 1) / 2
	nJoins := nDistances

	distanceInfo := fillDistanceLog(coords)

	clusterMap, nodeMap := getClusterMap(distanceInfo, nJoins, n)

	answer := scoreTopN(clusterMap, topN)

	_ = nodeMap

	fmt.Printf("\nFINAL ANSWER = %v\n", answer)

}

func mapComplete(clusterMap map[int][]int, nMembers int) bool {
	var ans bool
	keys := make([]int, 0, len(clusterMap))
	for k := range clusterMap {
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		if len(clusterMap[keys[0]]) == nMembers {
			ans = true
		} else {
			ans = false
		}
	}
	return ans
}

func lastDistance(distances []DistanceLog, nDistances int) float64 {
	last := distances[nDistances-1]
	return last.v1.X * last.v2.X
}

func scoreTopN(clusterMap map[int][]int, n int) int {
	runningProduct := 1
	// Extract keys from clusterMap
	keys := make([]int, 0, len(clusterMap))
	for k := range clusterMap {
		keys = append(keys, k)
	}

	//fmt.Printf("KEYS: %v\n", keys)

	// Sort keys by the length of their corresponding slices (descending order)
	sort.Slice(keys, func(i, j int) bool {
		return len(clusterMap[keys[i]]) > len(clusterMap[keys[j]])
	})

	// Now iterate through sorted keys
	for i, key := range keys {
		if i <= n-1 {
			//fmt.Printf("%v: %v (length: %d)\n", key, clusterMap[key], len(clusterMap[key]))
			runningProduct *= len(clusterMap[key])
		}
	}
	return runningProduct
}

func getClusterMap(distanceInfo []DistanceLog, nJoins int, nJunctions int) (map[int][]int, map[int]int) {
	clusterMap := make(map[int][]int)
	nodeMap := make(map[int]int)
	clusterID := 1
	for i := range nJoins {
		// join the two vectors in distanceInfo
		// what are there ids?
		id1 := distanceInfo[i].id1
		id2 := distanceInfo[i].id2
		val1, ok1 := nodeMap[id1]
		val2, ok2 := nodeMap[id2]

		fmt.Printf("CLUSTER MAP (%v): %v\n", i, clusterMap)

		if mapComplete(clusterMap, nJunctions) {
			fmt.Printf("DONE\n")
			fmt.Printf("The final connection was:\n")
			fmt.Printf("INFO %v - %v\n", i-1, distanceInfo[i-1])
			x1 := distanceInfo[i-1].v1.X
			x2 := distanceInfo[i-1].v2.X
			prod := x1 * x2
			fmt.Printf("X1: %v, X2: %v -- PROD: %v\n", x1, x2, prod)
			break
		}

		if ok1 && !ok2 {
			//fmt.Printf("%v is new. Joining to %v\n", id2, id1)
			// We have seen id1 already but not id2
			// look up nodeMap to see which cluster id1 is in
			clusterIDTemp := val1 // we already did this
			// now add id2 to this cluster
			clusterMap[clusterIDTemp] = append(clusterMap[clusterIDTemp], id2)
			// and make sure we add id2 to nodeMap now
			nodeMap[id2] = clusterIDTemp
		} else if !ok1 && ok2 {
			//fmt.Printf("%v is new. Joining to %v\n", id1, id2)
			// We have not seen id1, but we have seen id2
			clusterIDTemp := val2 // he was in this cluster
			// add id1 to this cluster
			clusterMap[clusterIDTemp] = append(clusterMap[clusterIDTemp], id1)
			//update nodeMap
			nodeMap[id1] = clusterIDTemp
		} else if !ok1 && !ok2 {
			//fmt.Printf("%v and %v are both new. Making new ID = %v\n", id1, id2, clusterID)
			// We haven't seen either of these ids before
			clusterIDTemp := clusterID // make new cluster ID
			// add both to cluster map
			clusterMap[clusterIDTemp] = append(clusterMap[clusterIDTemp], id1)
			clusterMap[clusterIDTemp] = append(clusterMap[clusterIDTemp], id2)
			// record that we have seen these
			nodeMap[id1] = clusterIDTemp
			nodeMap[id2] = clusterIDTemp
			// Increment globale clusterID ready for next time this happens
			clusterID++
		} else if ok1 && ok2 {
			//fmt.Printf("We have seen both %v and %v before.\n", id1, id2)
			// We have seen both ids already, so they are both already in clusters
			if val1 == val2 {
				//fmt.Printf("They are already in the same cluster = %v.\n", val1)
				// they already belong to the same cluster - do nothing
				// check to see if all nodes are now in the same (single) cluster
				continue
			} else {
				//fmt.Printf("Forcing %v to assimilate with %v.\n", id2, id1)
				// they are in different clusters. Force cluster 2 to assimilate to cluster 1.
				// get all the nodes that belong to cluster 2
				nodesFrom2 := clusterMap[val2]
				//fmt.Printf("Nodes from %v: %v\n", id2, nodesFrom2)
				// add them to cluster 1
				//fmt.Printf("Want to add them to %v\n", clusterMap[val1])
				clusterMap[val1] = append(clusterMap[val1], nodesFrom2...)
				// update where node is "pointing"
				//fmt.Printf("Updating Node references for %v\n", nodesFrom2)
				for _, node := range nodesFrom2 {
					nodeMap[node] = val1
				}
				//delete k,v pair for cluster 2
				delete(clusterMap, val2)
			}
		}

		//clusterMap[clusterID] = append(clusterMap[clusterID], distanceInfo[i].id1, distanceInfo[i].id2)
	}
	return clusterMap, nodeMap
}

func fillDistanceLog(coords []Vec3) []DistanceLog {
	var distances []DistanceLog
	n := len(coords)
	for i := range n - 1 {
		// take vector coords[i] and measure distances to other n-1 vectors
		for j := range n - 1 - i {
			d := calcDistEuclid(coords[i], coords[j+i+1])
			distances = append(distances, DistanceLog{i, j + i + 1, coords[i], coords[j+i+1], d, 0})
		}
	}
	// now we sort by ascending distance
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].d < distances[j].d
	})
	return distances
}

func calcDistEuclid(v1 Vec3, v2 Vec3) float64 {
	ans := math.Sqrt((v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y) + (v1.Z-v2.Z)*(v1.Z-v2.Z))
	return ans
}

func loadVec3File(path string) ([]Vec3, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var points []Vec3

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := s.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		coords := s.Split(line, ",")

		if len(coords) != 3 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}

		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		z, err3 := strconv.Atoi(coords[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("invalid integers on line: %q", line)
		}

		points = append(points, Vec3{float64(x), float64(y), float64(z)})
	}

	return points, scanner.Err()
}
